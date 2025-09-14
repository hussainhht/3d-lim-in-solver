package main

import (
	"flag"
	"fmt"
	lemin "lemin/handles"
	"log"
	"math"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// ---- short aliases ----
type v3 = rl.Vector3
func V3(x, y, z float32) v3 { return rl.NewVector3(x, y, z) }
func lerp(a, b v3, t float32) v3 {
	return V3(a.X+(b.X-a.X)*t, a.Y+(b.Y-a.Y)*t, a.Z+(b.Z-a.Z)*t)
}

// Visual ant assigned to a path
type VisAnt struct {
	PathKey string
	Phase   float32 // segments traversed (accumulated)
}

func main() {
	// CLI flags
	input := flag.String("input", "input.txt", "lem-in input file")
	speed := flag.Float64("speed", 0.6, "ant speed (segments per second)")
	scale := flag.Float64("scale", 1.6, "layout scale (spread rooms)")
	flag.Parse()

	// Parse farm
	farm := lemin.Exportinfo(*input)
	if farm == nil {
		log.Fatal("failed to parse input file")
	}

	// Best set of paths (pruned + non-overlapping + optimized for numberOfRounds)
	paths := lemin.FindPaths(farm)
	if len(paths) == 0 {
		log.Fatal("no paths from start to end")
	}

	// Room placement: centered + scaled
	roomPos := make(map[string]v3, len(farm.Rooms))
	var cx, cz float32
	for _, r := range farm.Rooms {
		cx += float32(r.X)
		cz += float32(r.Y)
	}
	if n := float32(len(farm.Rooms)); n > 0 {
		cx /= n
		cz /= n
	}
	s := float32(*scale)
	for name, r := range farm.Rooms {
		X := (float32(r.X) - cx) * s
		Z := (float32(r.Y) - cz) * s
		roomPos[name] = V3(X, 0, Z)
	}

	// Precompute path points (once)
	pathPoints := make(map[string][]v3, len(paths))
	for _, p := range paths {
		pts := make([]v3, len(p.Rooms))
		for i, rm := range p.Rooms {
			pts[i] = roomPos[rm.Name]
		}
		pathPoints[p.String()] = pts
	}

	// Ant distribution across paths (turns equalization algorithm)
	dist := computeAntDistribution(paths, farm.TotalAnts)

	// Visual ant herd according to distribution
	var ants []VisAnt
	for _, p := range paths {
		key := p.String()
		n := dist[key]
		if len(pathPoints[key]) < 2 || n <= 0 {
			continue
		}
		for i := 0; i < n; i++ {
			ants = append(ants, VisAnt{
				PathKey: key,
				Phase:   float32(i) * 0.35, // simple staggering
			})
		}
	}

	// Window + camera
	rl.InitWindow(1200, 800, "lem-in 3D Visualizer (raylib-go)")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	cam := rl.Camera3D{
		Position:   V3(0, 30, 40),
		Target:     V3(0, 0, 0),
		Up:         V3(0, 1, 0),
		Fovy:       60,
		Projection: rl.CameraPerspective,
	}

	antSpeed := float32(*speed)
	last := time.Now()

	// Main drawing loop
	for !rl.WindowShouldClose() {
		rl.UpdateCamera(&cam, rl.CameraFree)

		now := time.Now()
		dt := float32(now.Sub(last).Seconds())
		last = now

		for i := range ants {
			ants[i].Phase += antSpeed * dt
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.NewColor(18, 18, 22, 255))
		rl.BeginMode3D(cam)

		// Ground grid
		rl.DrawGrid(120, 1)

		// Links between rooms (gray lines)
		for _, r := range farm.Rooms {
			a := roomPos[r.Name]
			for _, nb := range r.Links {
				b := roomPos[nb.Name]
				if r.Name < nb.Name {
					rl.DrawLine3D(a, b, rl.Gray)
				}
			}
		}

		// Rooms (Start=Green, End=Red, normal=Blueish)
		for name, pos := range roomPos {
			col := rl.Color{R: 120, G: 200, B: 255, A: 255}
			radius := float32(0.5)
			if farm.StartRoom != nil && name == farm.StartRoom.Name {
				col, radius = rl.Green, 0.7
			}
			if farm.EndRoom != nil && name == farm.EndRoom.Name {
				col, radius = rl.Red, 0.7
			}
			rl.DrawSphere(pos, radius, col)
		}

		// Selected paths
		for _, p := range paths {
			pts := pathPoints[p.String()]
			if len(pts) < 2 {
				continue
			}
			for i := 0; i < len(pts)-1; i++ {
				rl.DrawLine3D(pts[i], pts[i+1], rl.LightGray)
			}
		}

		// Ants moving along their paths
		for _, a := range ants {
			pts := pathPoints[a.PathKey]
			totalSeg := len(pts) - 1
			if totalSeg <= 0 {
				continue
			}
			p := a.Phase
			seg := int(math.Floor(float64(p))) % totalSeg
			if seg < 0 {
				seg += totalSeg
			}
			t := p - float32(int(p))
			A := pts[seg]
			B := pts[seg+1]
			pos := lerp(A, B, t)
			rl.DrawSphere(pos, 0.24, rl.Yellow)
		}

		rl.EndMode3D()

		// HUD
		rl.DrawText("WASD + mouse  |  Wheel: zoom  |  R: reset camera", 10, 10, 18, rl.RayWhite)
		rl.DrawText(
			fmt.Sprintf("Rooms: %d  Links: ~%d  Paths: %d  Ants: %d",
				len(farm.Rooms), approxLinks(farm), len(paths), len(ants)),
			10, 34, 18, rl.RayWhite,
		)

		// Reset camera
		if rl.IsKeyPressed(rl.KeyR) {
			cam.Position = V3(0, 30, 40)
			cam.Target = V3(0, 0, 0)
			cam.Up = V3(0, 1, 0)
			cam.Fovy = 60
		}

		rl.EndDrawing()
	}
}

// Approximate number of links (each edge counted twice)
func approxLinks(f *lemin.Antsfarm) int {
	sum := 0
	for _, r := range f.Rooms {
		sum += len(r.Links)
	}
	return sum / 2
}

// Turns equalization: in each step, pick the path with minimum (len(path)+assigned)
func computeAntDistribution(paths []*lemin.Path, totalAnts int) map[string]int {
	dist := make(map[string]int)
	if len(paths) == 0 || totalAnts <= 0 {
		return dist
	}
	for _, p := range paths {
		dist[p.String()] = 0
	}
	remaining := totalAnts
	for remaining > 0 {
		minTurns := int(^uint(0) >> 1)
		var best *lemin.Path
		for _, p := range paths {
			turns := len(p.Rooms) + dist[p.String()]
			if turns < minTurns {
				minTurns = turns
				best = p
			}
		}
		dist[best.String()]++
		remaining--
	}
	return dist
}
