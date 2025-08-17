package lemin

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// readFile parses the lem-in input from text.txt in the current working directory
// and returns a populated Antsfarm structure.
func ReadFile() *Antsfarm {
	f, err := os.Open("text.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer f.Close()

	farm := &Antsfarm{rooms: make(map[string]*Room)}
	coordSet := make(map[string]string) // key: "x:y" -> room name

	scanner := bufio.NewScanner(f)
	lineNo := 0
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

		if !parsedAnts {
			if strings.HasPrefix(line, "#") { // comments before ants line
				continue
			}
			ants, err := strconv.Atoi(line)
			if err != nil || ants <= 0 {
				fmt.Println("Invalid ants:", line)
				return nil
			}
			farm.antsnum = ants
			parsedAnts = true
			continue
		}

		// Commands or comments
		if strings.HasPrefix(line, "##") {
			if linkSection { // commands not allowed after link section begins
				fmt.Println("Command after links started:", line)
				return nil
			}
			switch line {
			case "##start":
				if farm.startroom != nil {
					fmt.Println("Duplicate start command:", line)
					return nil
				}
				expectStart = true
				expectEnd = false
			case "##end":
				if farm.endroom != nil {
					fmt.Println("Duplicate end command:", line)
					return nil
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

		// Determine if link line
		if isLinkLine(line) {
			linkSection = true
			name1, name2 := parseLinkLine(line)
			r1, ok1 := farm.rooms[name1]
			r2, ok2 := farm.rooms[name2]
			if !ok1 || !ok2 {
				fmt.Println("Link references unknown room:", line)
				return nil
			}
			if err := addLink(r1, r2); err != nil {
				fmt.Println("Error adding link:", err)
				return nil
			}
			// After consuming a link, pending start/end expectation is invalid
			if expectStart || expectEnd {
				fmt.Println("Expected room after command, found link:", line)
				return nil
			}
			continue
		}

		// If link section already started, no more room definitions allowed
		if linkSection {
			fmt.Println("Room definition after links started:", line)
			return nil
		}

		// Room line
		parts := strings.Fields(line)
		if len(parts) != 3 {
			fmt.Println("Invalid room format:", line)
			return nil
		}
		name := parts[0]
		if err := validateRoomName(name); err != nil {
			fmt.Println("Invalid room name:", line)
			return nil
		}
		if _, exists := farm.rooms[name]; exists {
			fmt.Println("Duplicate room name:", line)
			return nil
		}
		x, errX := strconv.Atoi(parts[1])
		y, errY := strconv.Atoi(parts[2])
		if errX != nil || errY != nil {
			fmt.Println("Invalid room coordinates:", line)
			return nil
		}
		coordKey := fmt.Sprintf("%d:%d", x, y)
		if other := coordSet[coordKey]; other != "" {
			fmt.Println("Duplicate room coordinates:", line)
			return nil
		}
		r := &Room{name: name, x: x, y: y}
		farm.rooms[name] = r
		coordSet[coordKey] = name
		if expectStart {
			farm.startroom = r
			expectStart = false
		} else if expectEnd {
			farm.endroom = r
			expectEnd = false
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Scan error:", err)
		return nil
	}

	if !parsedAnts {
		fmt.Println("Missing ants line")
		return nil
	}
	if farm.startroom == nil {
		fmt.Println("Missing start room")
		return nil
	}
	if farm.endroom == nil {
		fmt.Println("Missing end room")
		return nil
	}
	if expectStart || expectEnd {
		fmt.Println("File ended before room after start/end command")
		return nil
	}
	return farm
}

// String returns a concise summary of the parsed farm.
func (f *Antsfarm) String() string {
	if f == nil {
		return "<nil farm>"
	}
	start := "<nil>"
	end := "<nil>"
	if f.startroom != nil {
		start = f.startroom.name
	}
	if f.endroom != nil {
		end = f.endroom.name
	}
	linkCount := 0
	for _, r := range f.rooms {
		linkCount += len(r.links)
	}
	// each undirected link counted twice
	linkCount /= 2
	return fmt.Sprintf("ants=%d rooms=%d links=%d start=%s end=%s", f.antsnum, len(f.rooms), linkCount, start, end)
}

// validateRoomName enforces lem-in constraints on room names.
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

// isLinkLine decides if a trimmed line represents a link (format name1-name2, no spaces).
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

// parseLinkLine splits a validated link line into room names.
func parseLinkLine(line string) (string, string) {
	parts := strings.SplitN(line, "-", 2)
	return parts[0], parts[1]
}

// addLink adds an undirected edge between two rooms, enforcing constraints.
func addLink(a, b *Room) error {
	if a == nil || b == nil {
		return errors.New("nil room in link")
	}
	if a == b {
		return errors.New("self-loop is forbidden")
	}
	// check duplicates in a's adjacency
	for _, l := range a.links {
		if l == b {
			return errors.New("duplicate link")
		}
	}
	a.links = append(a.links, b)
	b.links = append(b.links, a)
	return nil
}
