package lemin

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testDir = "../test"

// expectedTurns counts the turns of a reference result file.
func expectedTurns(t *testing.T, name string) int {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "result", name))
	if err != nil {
		t.Fatal(err)
	}
	n := 0
	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) != "" {
			n++
		}
	}
	return n
}

// TestExamples simulates every example map, checks the result is a legal
// solution and, where the repository ships a reference result, that it uses
// the same number of turns.
func TestExamples(t *testing.T) {
	cases := []struct {
		file      string
		reference string // "" when there is no reference turn count
	}{
		{"example00.txt", "result00.txt"},
		{"example01.txt", "result01.txt"},
		{"example02.txt", "result02.txt"},
		{"example03.txt", "result03.txt"},
		{"example04.txt", "result04.txt"},
		{"example05.txt", "result05.txt"},
		{"example06.txt", ""},
		{"example07.txt", ""},
	}
	for _, c := range cases {
		t.Run(c.file, func(t *testing.T) {
			farm := Exportinfo(filepath.Join(testDir, c.file))
			if farm == nil {
				t.Fatal("valid map was rejected")
			}
			res := Simulate(farm)
			if err := Verify(farm, res); err != nil {
				t.Fatal(err)
			}
			if c.reference != "" {
				if got, want := len(res.Turns), expectedTurns(t, c.reference); got != want {
					t.Errorf("turns = %d, want %d", got, want)
				}
			}
		})
	}
}

// TestBadInput checks that malformed maps are rejected by the parser.
func TestBadInput(t *testing.T) {
	cases := map[string]string{
		"zero ants":        "0\n##start\ns 0 0\n##end\ne 1 0\ns-e\n",
		"not a number":     "abc\n##start\ns 0 0\n##end\ne 1 0\ns-e\n",
		"no start":         "3\ns 0 0\n##end\ne 1 0\ns-e\n",
		"no end":           "3\n##start\ns 0 0\ne 1 0\ns-e\n",
		"duplicate room":   "3\n##start\ns 0 0\n##end\ne 1 0\ne 2 2\ns-e\n",
		"duplicate coords": "3\n##start\ns 0 0\n##end\ne 0 0\ns-e\n",
		"room starts in L": "3\n##start\ns 0 0\n##end\ne 1 0\nLx 3 3\ns-e\n",
		"unknown room":     "3\n##start\ns 0 0\n##end\ne 1 0\ns-x\n",
		"self loop":        "3\n##start\ns 0 0\n##end\ne 1 0\ns-s\ns-e\n",
		"duplicate link":   "3\n##start\ns 0 0\n##end\ne 1 0\ns-e\ne-s\n",
		"room after links": "3\n##start\ns 0 0\n##end\ne 1 0\ns-e\nz 5 5\n",
		"duplicate start":  "3\n##start\ns 0 0\n##start\nt 2 2\n##end\ne 1 0\ns-e\n",
		"empty file":       "",
	}
	for name, content := range cases {
		t.Run(name, func(t *testing.T) {
			file := filepath.Join(t.TempDir(), "map.txt")
			if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
				t.Fatal(err)
			}
			if Exportinfo(file) != nil {
				t.Error("malformed map was accepted")
			}
		})
	}
}

// TestSmallMaps covers shapes that stress the scheduler: a direct start-end
// tunnel, a single ant, a shared bottleneck and two competing routes.
func TestSmallMaps(t *testing.T) {
	cases := []struct {
		name    string
		content string
		turns   int
	}{
		{"direct tunnel", "3\n##start\ns 0 0\n##end\ne 1 0\ns-e\n", 3},
		{"single ant", "1\n##start\ns 0 0\na 1 0\n##end\ne 2 0\ns-a\na-e\n", 2},
		{"bottleneck", "5\n##start\ns 0 0\nm 1 0\na 2 1\nb 2 -1\n##end\ne 3 0\ns-m\nm-a\nm-b\na-e\nb-e\n", 7},
		{"two routes", "4\n##start\ns 0 0\na 1 0\nb 1 1\nc 2 0\nd 2 1\n##end\ne 3 0\ns-a\ns-b\na-c\na-d\nb-d\nc-e\nd-e\n", 4},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			file := filepath.Join(t.TempDir(), "map.txt")
			if err := os.WriteFile(file, []byte(c.content), 0o644); err != nil {
				t.Fatal(err)
			}
			farm := Exportinfo(file)
			if farm == nil {
				t.Fatal("valid map was rejected")
			}
			res := Simulate(farm)
			if err := Verify(farm, res); err != nil {
				t.Fatal(err)
			}
			if len(res.Turns) != c.turns {
				t.Errorf("turns = %d, want %d", len(res.Turns), c.turns)
			}
		})
	}
}

// TestTracks checks that per-ant tracks reproduce the turn list.
func TestTracks(t *testing.T) {
	farm := Exportinfo(filepath.Join(testDir, "example00.txt"))
	if farm == nil {
		t.Fatal("valid map was rejected")
	}
	res := Simulate(farm)
	tracks := res.Tracks(farm)
	if len(tracks) != farm.TotalAnts {
		t.Fatalf("tracks = %d, want %d", len(tracks), farm.TotalAnts)
	}
	steps := 0
	for id, track := range tracks {
		if len(track) == 0 || track[len(track)-1].Room != farm.EndRoom {
			t.Errorf("ant %d does not finish in the end room", id+1)
		}
		steps += len(track)
	}
	moves := 0
	for _, turn := range res.Turns {
		moves += len(turn)
	}
	if steps != moves {
		t.Errorf("tracks hold %d steps, turns hold %d movements", steps, moves)
	}
}

// TestVerifyRejectsIllegalResults makes sure the verifier is not vacuous.
func TestVerifyRejectsIllegalResults(t *testing.T) {
	newFarm := func() *Antsfarm {
		farm := Exportinfo(filepath.Join(testDir, "example00.txt"))
		if farm == nil {
			t.Fatal("valid map was rejected")
		}
		return farm
	}
	cases := map[string]*Result{
		"teleport":     {Turns: [][]Movement{{{1, "1"}}}},
		"unknown room": {Turns: [][]Movement{{{1, "nope"}}}},
		"ants missing": {Turns: [][]Movement{{{1, "2"}}}},
		"double move":  {Turns: [][]Movement{{{1, "2"}, {1, "3"}}}},
		"shared room":  {Turns: [][]Movement{{{1, "2"}, {2, "2"}}}},
		"unknown ant":  {Turns: [][]Movement{{{99, "2"}}}},
	}
	for name, res := range cases {
		t.Run(name, func(t *testing.T) {
			if Verify(newFarm(), res) == nil {
				t.Error("illegal result was accepted")
			}
		})
	}
}
