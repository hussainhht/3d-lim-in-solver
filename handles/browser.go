package lemin

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// BrowserResult is a pointer-free transport, not a second solver.
type BrowserResult struct {
	Rooms       []BrowserRoom `json:"rooms"`
	Tunnels     [][2]string   `json:"tunnels"`
	Start       string        `json:"start"`
	End         string        `json:"end"`
	AntCount    int           `json:"antCount"`
	Paths       [][]string    `json:"paths"`       // includes start for presentation
	Assignments []int         `json:"assignments"` // path index, indexed by ant ID - 1
	Turns       [][]Movement  `json:"turns"`
	Output      []string      `json:"output"` // exact CLI lines, including trailing spaces
	TotalTurns  int           `json:"totalTurns"`
	ComputeMS   float64       `json:"computeMs"`
	Engine      string        `json:"engine"`
}

type BrowserRoom struct {
	Name string `json:"name"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
}

type searchLimit string

func (e searchLimit) Error() string { return string(e) }

type searchBudget struct{ branches int }

func (b *searchBudget) enumeration(queued, found int) {
	if queued > 50000 || found > 256 {
		panic(searchLimit("This map exceeds the browser path-enumeration budget (50,000 queued paths / 256 candidates). Try fewer interconnected tunnels."))
	}
}
func (b *searchBudget) selection() {
	b.branches++
	if b.branches > 2000000 {
		panic(searchLimit("This map exceeds the browser search budget (2,000,000 backtracking steps). Try a simpler map."))
	}
}

// SolveText invokes the existing parser, Simulate and independent Verify in memory.
// Browser limits reject expensive inputs; they never prune or approximate a result.
func SolveText(text string) (out *BrowserResult, err error) {
	started := time.Now()
	defer func() {
		if recovered := recover(); recovered != nil {
			if limit, ok := recovered.(searchLimit); ok {
				out, err = nil, limit
				return
			}
			panic(recovered)
		}
	}()
	if len(text) > 128*1024 {
		return nil, fmt.Errorf("browser map limit: 128 KiB")
	}
	farm, err := ParseMap(strings.NewReader(text))
	if err != nil {
		return nil, err
	}
	tunnelCount := 0
	for _, r := range farm.Rooms {
		tunnelCount += len(r.Links)
	}
	if len(farm.Rooms) > 80 || tunnelCount/2 > 160 || farm.TotalAnts > 1000 {
		return nil, fmt.Errorf("browser limits: 80 rooms, 160 tunnels and 1,000 ants")
	}
	// Avoid native FindPaths' process-exit behavior for disconnected maps.
	seen := map[*Room]bool{farm.StartRoom: true}
	queue := []*Room{farm.StartRoom}
	for len(queue) > 0 {
		r := queue[0]
		queue = queue[1:]
		for _, n := range r.Links {
			if !seen[n] {
				seen[n] = true
				queue = append(queue, n)
			}
		}
	}
	if !seen[farm.EndRoom] {
		return nil, fmt.Errorf("no path connects start %q to end %q", farm.StartRoom.Name, farm.EndRoom.Name)
	}
	farm.budget = &searchBudget{}
	res := Simulate(farm)
	if err := Verify(farm, res); err != nil {
		return nil, fmt.Errorf("solver verification failed: %w", err)
	}
	out = &BrowserResult{Start: farm.StartRoom.Name, End: farm.EndRoom.Name, AntCount: farm.TotalAnts,
		Turns: res.Turns, TotalTurns: len(res.Turns), Assignments: make([]int, farm.TotalAnts), Engine: "Go · WebAssembly"}
	names := make([]string, 0, len(farm.Rooms))
	for name := range farm.Rooms {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		r := farm.Rooms[name]
		out.Rooms = append(out.Rooms, BrowserRoom{Name: r.Name, X: r.X, Y: r.Y})
		for _, n := range r.Links {
			if name < n.Name {
				out.Tunnels = append(out.Tunnels, [2]string{name, n.Name})
			}
		}
	}
	firstRoomPath := map[string]int{}
	for i, p := range res.Paths {
		names := []string{farm.StartRoom.Name}
		for _, r := range p.Rooms {
			names = append(names, r.Name)
		}
		out.Paths = append(out.Paths, names)
		firstRoomPath[p.Rooms[0].Name] = i
	}
	assigned := make([]bool, farm.TotalAnts)
	for _, turn := range res.Turns {
		var line strings.Builder
		for _, m := range turn {
			fmt.Fprintf(&line, "L%d-%s ", m.AntID, m.Room)
			if !assigned[m.AntID-1] {
				out.Assignments[m.AntID-1] = firstRoomPath[m.Room]
				assigned[m.AntID-1] = true
			}
		}
		out.Output = append(out.Output, line.String())
	}
	out.ComputeMS = float64(time.Since(started).Microseconds()) / 1000
	return out, nil
}
