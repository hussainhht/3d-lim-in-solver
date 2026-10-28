package lemin

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestBrowserParity(t *testing.T) {
	for i := 0; i < 8; i++ {
		t.Run(fmt.Sprintf("example%02d", i), func(t *testing.T) {
			text, err := os.ReadFile(fmt.Sprintf("../test/example%02d.txt", i))
			if err != nil {
				t.Fatal(err)
			}
			web, err := SolveText(string(text))
			if err != nil {
				t.Fatal(err)
			}
			farm, err := ParseMap(strings.NewReader(string(text)))
			if err != nil {
				t.Fatal(err)
			}
			native := Simulate(farm)
			if !reflect.DeepEqual(web.Turns, native.Turns) {
				t.Fatal("browser turns differ from native Simulate")
			}
			for p, path := range native.Paths {
				for j, r := range path.Rooms {
					if web.Paths[p][j+1] != r.Name {
						t.Fatal("path mismatch")
					}
				}
			}
			for id, track := range native.Tracks(farm) {
				path := web.Paths[web.Assignments[id]]
				if len(track) != len(path)-1 {
					t.Fatalf("L%d assignment length mismatch", id+1)
				}
				for j, step := range track {
					if step.Room.Name != path[j+1] {
						t.Fatalf("L%d assignment mismatch", id+1)
					}
				}
			}
			for j, turn := range native.Turns {
				var expected strings.Builder
				for _, m := range turn {
					fmt.Fprintf(&expected, "L%d-%s ", m.AntID, m.Room)
				}
				if web.Output[j] != expected.String() {
					t.Fatal("CLI output line mismatch")
				}
			}
		})
	}
}

func TestBrowserErrorsAndRecovery(t *testing.T) {
	valid := "3\n##start\ns 0 0\n##end\ne 1 0\ns-e\n"
	cases := map[string]string{
		"disconnected":     "3\n##start\ns 0 0\n##end\ne 1 0\n",
		"ants limit":       strings.Replace(valid, "3\n", "1001\n", 1),
		"parser ant limit": strings.Replace(valid, "3\n", "1000001\n", 1),
		"invalid":          strings.Replace(valid, "s-e", "s-missing", 1),
		"size":             strings.Repeat("# comment\n", 15000),
		"empty":            "",
	}
	for name, text := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := SolveText(text); err == nil {
				t.Fatal("expected error")
			}
			if _, err := SolveText(valid); err != nil {
				t.Fatalf("recovery failed: %v", err)
			}
		})
	}
}

func TestBrowserDenseMapBudget(t *testing.T) {
	// Complete 12-room graph overwhelms simple-path enumeration; must reject.
	var text strings.Builder
	text.WriteString("4\n")
	for i := 0; i < 12; i++ {
		if i == 0 {
			text.WriteString("##start\n")
		}
		if i == 11 {
			text.WriteString("##end\n")
		}
		fmt.Fprintf(&text, "r%d %d 0\n", i, i)
	}
	for i := 0; i < 12; i++ {
		for j := i + 1; j < 12; j++ {
			fmt.Fprintf(&text, "r%d-r%d\n", i, j)
		}
	}
	if _, err := SolveText(text.String()); err == nil || !strings.Contains(err.Error(), "budget") {
		t.Fatalf("expected budget rejection, got %v", err)
	}
}

func BenchmarkBrowserExamples(b *testing.B) {
	for _, name := range []string{"../test/example00.txt", "../test/example05.txt", "../test/example07.txt", "../3D/input.txt"} {
		text, _ := os.ReadFile(name)
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if _, err := SolveText(string(text)); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
