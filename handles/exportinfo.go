package lemin

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
)

// readFile parses the lem-in input from text.txt in the current working directory
// * and returns a populated Antsfarm structure.
func Exportinfo(f_name string) *Antsfarm {
	if !strings.HasSuffix(strings.TrimSpace(f_name), ".txt") {
		fmt.Fprintln(os.Stderr, "File name must end with .txt")
		os.Exit(1)
	}
	//? f_name := "input file name"
	f, err := os.Open(f_name)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error opening file:", err)
		return nil
	}
	defer f.Close()

	farm, err := ParseMap(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil
	}
	return farm
}

// ParseMap is the shared, in-memory parser used by the CLI and WebAssembly.
// It has the original parsing rules and returns diagnostics without exiting.
func ParseMap(reader io.Reader) (*Antsfarm, error) {
	farm := &Antsfarm{Rooms: make(map[string]*Room)}
	coordSet := make(map[string]string) // key: "x:y" -> room name

	scanner := bufio.NewScanner(reader)
	lineNo := 0
	fail := func(parts ...any) (*Antsfarm, error) {
		return nil, fmt.Errorf("line %d: %s", lineNo, strings.TrimSpace(fmt.Sprintln(parts...)))
	}
	parsedAnts := false
	linkSection := false
	expectStart := false
	expectEnd := false

	for scanner.Scan() {
		raw := scanner.Text()
		lineNo++
		line := strings.TrimSpace(raw)

		if line == "" { // skip empty lines (not specified but harmless)
			continue
		}

		// Check for valid ants line
		if !parsedAnts {
			if strings.HasPrefix(line, "#") { // comments before ants line
				continue
			}
			ants, err := strconv.Atoi(line)
			if err != nil || ants <= 0 {
				return fail("Invalid ants:", line)
			}

			// Check for too many ants
			if ants > 1000000 {
				return fail("Too many ants (>1000000)")
			}

			farm.TotalAnts = ants
			parsedAnts = true
			continue
		}

		// Commands or comments
		if strings.HasPrefix(line, "##") {

			if linkSection && (line == "##start" || line == "##end") { // commands not allowed after link section begins
				return fail("Command after links started:", line)
			}

			// Check for valid start/end commands
			switch line {
			case "##start":
				if farm.StartRoom != nil {
					return fail("Duplicate start command:", line)
				}
				expectStart = true
				expectEnd = false
			case "##end":
				if farm.EndRoom != nil {
					return fail("Duplicate end command:", line)
				}
				expectEnd = true
				expectStart = false
			default:
				// other double-hash lines treated as comments per spec (only start/end meaningful)
			}
			continue
		}

		if strings.HasPrefix(line, "#") { // normal comment
			continue
		}

		// * Determine if link line
		if isLinkLine(line) {

			linkSection = true

			// Reset start/end expectations
			name1, name2 := parseLinkLine(line)
			// check that both room names exist in your room map before proceeding
			r1, ok1 := farm.Rooms[name1]
			r2, ok2 := farm.Rooms[name2]
			if !ok1 || !ok2 {
				return fail("Link references unknown room:", line)
			}

			// * Add link, checking for self-loops and duplicates
			if err := addLink(r1, r2); err != nil {
				return fail("Error adding link:", err)
			}

			// After consuming a link, pending start/end expectation is invalid
			if expectStart || expectEnd {
				return fail("Expected room after command, found link:", line)
			}
			continue
		}

		// * If link section already started, no more room definitions allowed
		if linkSection {
			return fail("Room definition after links started:", line)
		}

		// Room line
		// Format: name x y
		parts := strings.Fields(line)
		if len(parts) != 3 {
			return fail("Invalid room format:", line)
		}

		// Validate room name
		name := parts[0]
		if err := validateRoomName(name); err != nil {
			return fail("Invalid room name:", line)
		}
		//exists is a boolean that is true if the key exists, otherwise false.
		if _, exists := farm.Rooms[name]; exists {
			return fail("Duplicate room name:", line)
		}
		x, errX := strconv.Atoi(parts[1])
		y, errY := strconv.Atoi(parts[2])
		if errX != nil || errY != nil {
			return fail("Invalid room coordinates:", line)
		}
		// key is "x:y"
		// value is room name
		coordKey := fmt.Sprintf("%d:%d", x, y)
		if other := coordSet[coordKey]; other != "" {
			return fail("Duplicate room coordinates:", line)
		}
		// This creates a new Room struct instance and takes its pointer.
		r := &Room{Name: name, X: x, Y: y}
		// This line adds the new room r to that map using the room’s name as the key.
		farm.Rooms[name] = r
		coordSet[coordKey] = name
		if expectStart {
			farm.StartRoom = r
			expectStart = false
		} else if expectEnd {
			farm.EndRoom = r
			expectEnd = false
		}
	}
	if err := scanner.Err(); err != nil {
		return fail("Scan error:", err)
	}

	if !parsedAnts {
		return fail("Missing ants line")
	}
	if farm.StartRoom == nil {
		return fail("Missing start room")
	}
	if farm.EndRoom == nil {
		return fail("Missing end room")
	}
	if expectStart || expectEnd {
		return fail("File ended before room after start/end command")
	}
	return farm, nil
}

// --------------------- Helpers ----------------------------------------------------------------------

// *String returns a concise summary of the parsed farm.
func (f *Antsfarm) String() string {
	if f == nil {
		return "<nil farm>"
	}
	start := "<nil>"
	end := "<nil>"
	if f.StartRoom != nil {
		start = f.StartRoom.Name
	}
	if f.EndRoom != nil {
		end = f.EndRoom.Name
	}
	linkCount := 0
	for _, r := range f.Rooms {
		linkCount += len(r.Links)
	}
	// each undirected link counted twice
	linkCount /= 2
	return fmt.Sprintf("ants=%d rooms=%d links=%d start=%s end=%s", f.TotalAnts, len(f.Rooms), linkCount, start, end)
}

// *validateRoomName enforces lem-in constraints on room names.
func validateRoomName(name string) error {
	if name == "" {
		return errors.New("empty room name")
	}
	if name[0] == 'L' {
		return fmt.Errorf("invalid room name %q (starts with 'L')", name)
	}
	if strings.Contains(name, "-") {
		return fmt.Errorf("invalid room name %q (contains '-')", name)
	}
	return nil
}

// *isLinkLine decides if a trimmed line represents a link (format name1-name2, no spaces).
func isLinkLine(line string) bool {
	if strings.Contains(line, " ") || strings.Contains(line, "\t") {
		return false
	}
	if strings.Count(line, "-") != 1 { // exactly one hyphen
		return false
	}
	parts := strings.Split(line, "-")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return false
	}
	return true
}

// *parseLinkLine splits a validated link line into room names.
func parseLinkLine(line string) (string, string) {
	parts := strings.SplitN(line, "-", 2)
	return parts[0], parts[1]
}

// * addLink adds an undirected edge between two rooms, enforcing constraints.
func addLink(a, b *Room) error {
	if a == nil || b == nil { // prevents linking to a room that does not exist (nil pointer)
		return errors.New("nil room in link")
	}

	if a == b { // self loop could cause infinite loops or incorrect graph logic.
		return errors.New("self-loop is forbidden")
	}

	// check duplicates in a's adjacency
	if slices.Contains(a.Links, b) {
		return errors.New("duplicate link")
	}

	a.Links = append(a.Links, b)
	b.Links = append(b.Links, a)
	return nil
}
