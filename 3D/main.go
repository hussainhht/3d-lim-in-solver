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

type v3 = rl.Vector3

func V3(x, y, z float32) v3 { return rl.NewVector3(x, y, z) }

func lerp(a, b v3, t float32) v3 {
	return V3(a.X+(b.X-a.X)*t, a.Y+(b.Y-a.Y)*t, a.Z+(b.Z-a.Z)*t)
}

func main() {
	input := flag.String("input", "input.txt", "lem-in input file")
	speed := flag.Float64("speed", 0.6, "ant speed (segments per second)")
	scale := flag.Float64("scale", 1.6, "layout scale (spread rooms)")
	flag.Parse()

	farm := lemin.Exportinfo(*input)
	if farm == nil {
		log.Fatal("failed to parse input file")
	}

	paths := lemin.FindPaths(farm)
	if len(paths) == 0 {
		log.Fatal("no paths from start to end")
	}
	shortestIdx := 0
	for i := 1; i < len(paths); i++ {
		if len(paths[i].Rooms) < len(paths[shortestIdx].Rooms) {
			shortestIdx = i
		}
	}
	shortest := paths[shortestIdx]

	roomPos := make(map[string]v3, len(farm.Rooms))
	var cx, cz float32
	for _, r := range farm.Rooms {
		cx += float32(r.X)
		cz += float32(r.Y)
	}
	if len(farm.Rooms) > 0 {
		cx /= float32(len(farm.Rooms))
		cz /= float32(len(farm.Rooms))
	}
	s := float32(*scale)
	for name, r := range farm.Rooms {
		X := (float32(r.X) - cx) * s
		Z := (float32(r.Y) - cz) * s
		roomPos[name] = V3(X, 0, Z)
	}

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

	//! Ants
	ants := farm.TotalAnts
	phase := make([]float32, ants)
	for i := 0; i < ants; i++ {
		phase[i] = float32(i) * 0.35
	}
	antSpeed := float32(*speed)

	last := time.Now()

	for !rl.WindowShouldClose() {
		rl.UpdateCamera(&cam, rl.CameraFree)

		now := time.Now()
		dt := float32(now.Sub(last).Seconds())
		last = now

		for i := 0; i < ants; i++ {
			phase[i] += antSpeed * dt
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.NewColor(18, 18, 22, 255))
		rl.BeginMode3D(cam)

		rl.DrawGrid(120, 1)

		for _, r := range farm.Rooms {
			a := roomPos[r.Name]
			for _, nb := range r.Links {
				b := roomPos[nb.Name]
				if r.Name < nb.Name {
					rl.DrawLine3D(a, b, rl.Gray)
				}
			}
		}

		for name, pos := range roomPos {
			col := rl.Color{R: 120, G: 200, B: 255, A: 255}
			radius := float32(0.5)
			if farm.StartRoom != nil && name == farm.StartRoom.Name {
				col = rl.Green
				radius = 0.7
			}
			if farm.EndRoom != nil && name == farm.EndRoom.Name {
				col = rl.Red
				radius = 0.7
			}
			rl.DrawSphere(pos, radius, col)
		}

		if len(shortest.Rooms) >= 2 {
			points := make([]v3, len(shortest.Rooms))
			for i, rm := range shortest.Rooms {
				points[i] = roomPos[rm.Name]
			}

			for i := 0; i < len(points)-1; i++ {
				rl.DrawLine3D(points[i], points[i+1], rl.LightGray)
			}

			totalSeg := len(points) - 1
			for i := 0; i < ants; i++ {
				p := phase[i]
				seg := int(math.Floor(float64(p))) % totalSeg
				t := p - float32(int(p))
				a := points[seg]
				b := points[seg+1]
				pos := lerp(a, b, t)
				rl.DrawSphere(pos, 0.24, rl.Yellow)
			}
		}

		rl.EndMode3D()

		rl.DrawText("WASD + mouse: move   |  Wheel: zoom   |  R: reset camera", 10, 10, 18, rl.RayWhite)
		rl.DrawText(fmt.Sprintf("Rooms: %d  Links: ~%d  Ants demo: %d  Path len: %d",
			len(farm.Rooms), approxLinks(farm), ants, len(shortest.Rooms)), 10, 34, 18, rl.RayWhite)

		if rl.IsKeyPressed(rl.KeyR) {
			cam.Position = V3(0, 30, 40)
			cam.Target = V3(0, 0, 0)
			cam.Up = V3(0, 1, 0)
			cam.Fovy = 60
		}

		rl.EndDrawing()
	}
}

func approxLinks(f *lemin.Antsfarm) int {
	sum := 0
	for _, r := range f.Rooms {
		sum += len(r.Links)
	}
	return sum / 2
}
