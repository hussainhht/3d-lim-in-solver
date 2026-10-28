// Command 3D is a raylib-go visualizer for the lem-in solver.
//
// It parses the colony, runs the real solver (lemin.Simulate) and replays the
// resulting turns in 3D: every ant moves exactly as in the CLI output.
package main

import (
	"flag"
	"fmt"
	lemin "lemin/handles"
	"log"
	"math"
	"os"
	"sort"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// ---- short aliases ----
type v3 = rl.Vector3

func V3(x, y, z float32) v3 { return rl.NewVector3(x, y, z) }
func lerp(a, b v3, t float32) v3 {
	return V3(a.X+(b.X-a.X)*t, a.Y+(b.Y-a.Y)*t, a.Z+(b.Z-a.Z)*t)
}

const (
	roomRadius    = 0.5
	endRoomRadius = 0.7
	antRadius     = 0.24
	antHeight     = 0.55 // ants hover above the room spheres so they stay visible
	maxIDLabels   = 40   // ant ids are only drawn on small colonies
	maxRoomLabels = 80   // room names are drawn by default on small maps
)

var pathPalette = []rl.Color{
	{R: 255, G: 196, B: 61, A: 255},
	{R: 255, G: 112, B: 166, A: 255},
	{R: 94, G: 226, B: 160, A: 255},
	{R: 178, G: 140, B: 255, A: 255},
	{R: 255, G: 146, B: 76, A: 255},
	{R: 90, G: 214, B: 255, A: 255},
}

// scene holds everything that is derived from the solver result.
type scene struct {
	farm    *lemin.Antsfarm
	result  *lemin.Result
	tracks  [][]lemin.Step // per ant (id-1): the rooms it enters, by turn
	pos     map[string]v3  // room name -> world position
	arrived []int          // arrived[k] = ants in the end room after k turns
	turns   int
}

func newScene(farm *lemin.Antsfarm, res *lemin.Result, scale float32) *scene {
	s := &scene{
		farm:   farm,
		result: res,
		tracks: res.Tracks(farm),
		pos:    make(map[string]v3, len(farm.Rooms)),
		turns:  len(res.Turns),
	}

	// Rooms: centered on the map's centroid, then scaled.
	var cx, cz float32
	for _, r := range farm.Rooms {
		cx += float32(r.X)
		cz += float32(r.Y)
	}
	if n := float32(len(farm.Rooms)); n > 0 {
		cx /= n
		cz /= n
	}
	for name, r := range farm.Rooms {
		s.pos[name] = V3((float32(r.X)-cx)*scale, 0, (float32(r.Y)-cz)*scale)
	}

	s.arrived = make([]int, s.turns+1)
	for _, track := range s.tracks {
		if n := len(track); n > 0 && track[n-1].Room == farm.EndRoom {
			s.arrived[track[n-1].Turn]++
		}
	}
	for k := 1; k <= s.turns; k++ {
		s.arrived[k] += s.arrived[k-1]
	}
	return s
}

// extent is the largest horizontal distance of any room from the origin.
func (s *scene) extent() float32 {
	var e float32
	for _, p := range s.pos {
		if d := float32(math.Hypot(float64(p.X), float64(p.Z))); d > e {
			e = d
		}
	}
	return e
}

// clusterOffset spreads ants that wait in the start room or sit in the end room
// around it (golden-angle spiral) so they do not overlap.
func clusterOffset(id int) v3 {
	angle := float64(id) * 2.399963
	r := 1.1 + 0.09*math.Sqrt(float64(id))
	return V3(float32(math.Cos(angle)*r), 0, float32(math.Sin(angle)*r))
}

func (s *scene) roomPoint(id int, room *lemin.Room) v3 {
	p := s.pos[room.Name]
	if room == s.farm.StartRoom || room == s.farm.EndRoom {
		o := clusterOffset(id)
		p = V3(p.X+o.X, p.Y, p.Z+o.Z)
	}
	p.Y = antHeight
	return p
}

// antPosition returns where ant id is at animation time tau (in turns) and
// whether it is travelling between two rooms or resting in the start/end room.
func (s *scene) antPosition(id int, tau float32) (p v3, moving, parked bool) {
	track := s.tracks[id-1]
	k := int(math.Floor(float64(tau)))

	// number of moves the ant has completed after k turns
	n := sort.Search(len(track), func(i int) bool { return track[i].Turn > k })
	cur := s.farm.StartRoom
	if n > 0 {
		cur = track[n-1].Room
	}
	from := s.roomPoint(id, cur)

	if n < len(track) && track[n].Turn == k+1 { // moving during turn k+1
		f := tau - float32(k)
		f = f * f * (3 - 2*f) // smoothstep
		return lerp(from, s.roomPoint(id, track[n].Room), f), true, false
	}
	return from, false, cur == s.farm.StartRoom || cur == s.farm.EndRoom
}

// turnLine renders the solver output line of a turn ("L1-2 L2-3 ...").
func (s *scene) turnLine(turn int) string {
	if turn < 1 || turn > s.turns {
		return ""
	}
	var b strings.Builder
	for i, m := range s.result.Turns[turn-1] {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "L%d-%s", m.AntID, m.Room)
	}
	return b.String()
}

// ---- orbit camera ----

type orbit struct {
	target     v3
	yaw, pitch float32
	dist       float32
	home       float32
}

func newOrbit(extent float32) *orbit {
	d := float32(math.Min(300, math.Max(16, float64(extent)*1.7+6)))
	return &orbit{pitch: 0.75, dist: d, home: d}
}

func (o *orbit) reset() {
	o.target, o.yaw, o.pitch, o.dist = v3{}, 0, 0.75, o.home
}

func (o *orbit) camera() rl.Camera3D {
	cp, sp := math.Cos(float64(o.pitch)), math.Sin(float64(o.pitch))
	sy, cy := math.Sin(float64(o.yaw)), math.Cos(float64(o.yaw))
	eye := V3(
		o.target.X+o.dist*float32(cp*sy),
		o.target.Y+o.dist*float32(sp),
		o.target.Z+o.dist*float32(cp*cy),
	)
	return rl.Camera3D{
		Position:   eye,
		Target:     o.target,
		Up:         V3(0, 1, 0),
		Fovy:       55,
		Projection: rl.CameraPerspective,
	}
}

func (o *orbit) update(dt float32) {
	if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		d := rl.GetMouseDelta()
		o.yaw -= d.X * 0.005
		o.pitch += d.Y * 0.005
		o.pitch = float32(math.Min(1.5, math.Max(0.05, float64(o.pitch))))
	}
	if w := rl.GetMouseWheelMove(); w != 0 {
		o.dist *= 1 - w*0.1
		o.dist = float32(math.Min(400, math.Max(3, float64(o.dist))))
	}

	sy, cy := float32(math.Sin(float64(o.yaw))), float32(math.Cos(float64(o.yaw)))
	forward := V3(-sy, 0, -cy)
	right := V3(cy, 0, -sy)

	pan := func(f, r float32) {
		o.target.X += forward.X*f + right.X*r
		o.target.Z += forward.Z*f + right.Z*r
	}
	if rl.IsMouseButtonDown(rl.MouseButtonRight) {
		d := rl.GetMouseDelta()
		k := o.dist * 0.0015
		pan(d.Y*k, -d.X*k)
	}
	step := o.dist * 0.8 * dt
	if rl.IsKeyDown(rl.KeyW) {
		pan(step, 0)
	}
	if rl.IsKeyDown(rl.KeyS) {
		pan(-step, 0)
	}
	if rl.IsKeyDown(rl.KeyD) {
		pan(0, step)
	}
	if rl.IsKeyDown(rl.KeyA) {
		pan(0, -step)
	}
}

// ---- drawing helpers ----

func antColor(id int) rl.Color {
	return rl.ColorFromHSV(float32((id*47)%360), 0.7, 1.0)
}

// onScreen reports whether a world point is in front of the camera.
func onScreen(cam rl.Camera3D, p v3) bool {
	f := V3(cam.Target.X-cam.Position.X, cam.Target.Y-cam.Position.Y, cam.Target.Z-cam.Position.Z)
	d := V3(p.X-cam.Position.X, p.Y-cam.Position.Y, p.Z-cam.Position.Z)
	return f.X*d.X+f.Y*d.Y+f.Z*d.Z > 0
}

// drawGround draws a faint square grid (2 units per cell) just below the rooms.
func drawGround(half int) {
	col := rl.Fade(rl.White, 0.07)
	h := float32(half)
	for i := -half; i <= half; i += 2 {
		f := float32(i)
		rl.DrawLine3D(V3(f, -0.6, -h), V3(f, -0.6, h), col)
		rl.DrawLine3D(V3(-h, -0.6, f), V3(h, -0.6, f), col)
	}
}

func drawCentered(text string, cx, y, size int32, col rl.Color) {
	rl.DrawText(text, cx-rl.MeasureText(text, size)/2, y, size, col)
}

func main() {
	input := flag.String("input", "input.txt", "lem-in input file")
	speed := flag.Float64("speed", 0.8, "playback speed (turns per second)")
	scale := flag.Float64("scale", 1.6, "layout scale (spread rooms)")
	check := flag.Bool("check", false, "solve, verify the result and exit without opening a window")
	shot := flag.String("screenshot", "", "save a screenshot to this PNG file and exit (development aid)")
	shotTurn := flag.Float64("shot-turn", -1, "turn to capture with -screenshot (-1 = after the last turn)")
	flag.Parse()

	// Parse the colony and run the real solver. Diagnostics go to stderr.
	farm := lemin.Exportinfo(*input)
	if farm == nil {
		log.Fatal("failed to parse input file")
	}
	res := lemin.Simulate(farm)
	if err := lemin.Verify(farm, res); err != nil {
		log.Fatalf("solver produced an illegal result: %v", err)
	}
	sc := newScene(farm, res, float32(*scale))

	if *check {
		fmt.Fprintf(os.Stderr, "ok: %d ants, %d rooms, %d paths, %d turns\n",
			farm.TotalAnts, len(farm.Rooms), len(res.Paths), sc.turns)
		return
	}

	// Path polylines (start room included) for highlighting.
	pathPoints := make([][]v3, len(res.Paths))
	for i, p := range res.Paths {
		pts := []v3{sc.pos[farm.StartRoom.Name]}
		for _, r := range p.Rooms {
			pts = append(pts, sc.pos[r.Name])
		}
		pathPoints[i] = pts
	}

	// Window + camera
	rl.SetTraceLogLevel(rl.LogWarning) // keep raylib's INFO chatter out of the console
	rl.SetConfigFlags(rl.FlagWindowResizable | rl.FlagMsaa4xHint)
	rl.InitWindow(1200, 800, "lem-in 3D Visualizer (raylib-go)")
	defer rl.CloseWindow()
	rl.SetWindowPosition(60, 40)
	rl.SetTargetFPS(60)

	cam3 := newOrbit(sc.extent())
	gridHalf := int(math.Ceil(float64(sc.extent())/2))*2 + 10

	tau := float32(0) // animation time, in turns
	playing := true
	playSpeed := float32(*speed)
	showLabels := len(farm.Rooms) <= maxRoomLabels
	showIDs := farm.TotalAnts <= 12
	T := float32(sc.turns)

	frame := 0
	if *shot != "" {
		playing = false
		tau = T
		if *shotTurn >= 0 {
			tau = float32(math.Min(float64(T), *shotTurn))
		}
	}

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()

		// ---- input ----
		cam3.update(dt)
		if rl.IsKeyPressed(rl.KeySpace) {
			if tau >= T {
				tau = 0
			}
			playing = !playing
		}
		if rl.IsKeyPressed(rl.KeyEnter) {
			tau, playing = 0, true
		}
		if rl.IsKeyPressed(rl.KeyRight) {
			playing = false
			tau = float32(math.Min(float64(T), math.Floor(float64(tau))+1))
		}
		if rl.IsKeyPressed(rl.KeyLeft) {
			playing = false
			tau = float32(math.Max(0, math.Ceil(float64(tau))-1))
		}
		if rl.IsKeyPressed(rl.KeyEqual) || rl.IsKeyPressed(rl.KeyKpAdd) {
			playSpeed = float32(math.Min(30, float64(playSpeed)*1.25))
		}
		if rl.IsKeyPressed(rl.KeyMinus) || rl.IsKeyPressed(rl.KeyKpSubtract) {
			playSpeed = float32(math.Max(0.1, float64(playSpeed)/1.25))
		}
		if rl.IsKeyPressed(rl.KeyL) {
			showLabels = !showLabels
		}
		if rl.IsKeyPressed(rl.KeyI) {
			showIDs = !showIDs
		}
		if rl.IsKeyPressed(rl.KeyR) {
			cam3.reset()
		}

		// ---- update ----
		if playing {
			tau += playSpeed * dt
			if tau >= T {
				tau, playing = T, false
			}
		}
		cam := cam3.camera()
		k := int(math.Floor(float64(tau)))
		turnNow := int(math.Ceil(float64(tau)))

		// ---- draw ----
		rl.BeginDrawing()
		rl.ClearBackground(rl.NewColor(18, 18, 22, 255))
		rl.BeginMode3D(cam)

		drawGround(gridHalf)

		// All tunnels (gray)
		for _, r := range farm.Rooms {
			a := sc.pos[r.Name]
			for _, nb := range r.Links {
				if r.Name < nb.Name {
					rl.DrawLine3D(a, sc.pos[nb.Name], rl.Gray)
				}
			}
		}

		// Tunnels of the paths chosen by the solver (one color per path)
		for i, pts := range pathPoints {
			col := pathPalette[i%len(pathPalette)]
			for j := 0; j < len(pts)-1; j++ {
				rl.DrawCylinderEx(pts[j], pts[j+1], 0.07, 0.07, 6, col)
			}
		}

		// Rooms (Start=Green, End=Red, normal=Blueish)
		for name, pos := range sc.pos {
			col := rl.Color{R: 120, G: 200, B: 255, A: 255}
			radius := float32(roomRadius)
			if name == farm.StartRoom.Name {
				col, radius = rl.Green, endRoomRadius
			}
			if name == farm.EndRoom.Name {
				col, radius = rl.Red, endRoomRadius
			}
			rl.DrawSphere(pos, radius, col)
		}

		// Ants
		type label struct {
			id int
			p  v3
		}
		var labels []label
		for id := 1; id <= farm.TotalAnts; id++ {
			p, moving, parked := sc.antPosition(id, tau)
			rl.DrawSphereEx(p, antRadius, 8, 8, antColor(id))
			if showIDs && farm.TotalAnts <= maxIDLabels && (moving || !parked) {
				labels = append(labels, label{id, p})
			}
		}

		rl.EndMode3D()

		// ---- 2D overlay ----
		w, h := int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight())

		if showLabels {
			for name, p := range sc.pos {
				if !onScreen(cam, p) {
					continue
				}
				below := rl.GetWorldToScreen(V3(p.X, -0.9, p.Z), cam)
				drawCentered(name, int32(below.X), int32(below.Y), 16, rl.RayWhite)
				tag, col := "", rl.Green
				switch name {
				case farm.StartRoom.Name:
					tag = "START"
				case farm.EndRoom.Name:
					tag, col = "END", rl.Red
				}
				if tag != "" {
					above := rl.GetWorldToScreen(V3(p.X, 1.3, p.Z), cam)
					drawCentered(tag, int32(above.X), int32(above.Y)-8, 18, col)
				}
			}
		}
		for _, l := range labels {
			if onScreen(cam, l.p) {
				sp := rl.GetWorldToScreen(V3(l.p.X, l.p.Y+0.5, l.p.Z), cam)
				drawCentered(fmt.Sprintf("L%d", l.id), int32(sp.X), int32(sp.Y)-8, 14, antColor(l.id))
			}
		}

		// HUD panel
		status := "playing"
		switch {
		case tau >= T:
			status = "finished"
		case !playing:
			status = "paused"
		}
		rl.DrawRectangle(0, 0, 440, 118, rl.Fade(rl.Black, 0.55))
		rl.DrawText("LEM-IN  3D", 12, 8, 22, rl.RayWhite)
		rl.DrawText(fmt.Sprintf("Turn %d / %d   (%s, %.2g turns/s)", turnNow, sc.turns, status, playSpeed), 12, 36, 18, rl.RayWhite)
		rl.DrawText(fmt.Sprintf("Ants arrived: %d / %d", sc.arrived[k], farm.TotalAnts), 12, 58, 18, rl.RayWhite)
		rl.DrawText(fmt.Sprintf("Rooms: %d   Paths used: %d   Ants: %d", len(farm.Rooms), len(res.Paths), farm.TotalAnts), 12, 80, 18, rl.LightGray)

		// Progress bar
		rl.DrawRectangle(0, h-6, w, 6, rl.Fade(rl.White, 0.15))
		if T > 0 {
			rl.DrawRectangle(0, h-6, int32(float32(w)*tau/T), 6, rl.Green)
		}

		// Solver output for the turn being animated
		if line := sc.turnLine(turnNow); line != "" {
			if limit := int(w / 9); len(line) > limit && limit > 3 {
				line = line[:limit-3] + "..."
			}
			rl.DrawText(line, 12, h-80, 18, rl.Yellow)
		}
		rl.DrawText("Space: play/pause  Enter: restart  Left/Right: step turn  +/-: speed  L: room names  I: ant ids",
			12, h-52, 16, rl.LightGray)
		rl.DrawText("LMB drag: orbit  RMB drag / WASD: pan  Wheel: zoom  R: reset camera  Esc: quit", 12, h-30, 16, rl.LightGray)

		if tau >= T {
			drawCentered(fmt.Sprintf("All %d ants reached %s in %d turns", farm.TotalAnts, farm.EndRoom.Name, sc.turns),
				w/2, 14, 26, rl.Green)
		}

		rl.EndDrawing()

		// Development aid: capture the framebuffer after a few frames and quit.
		if *shot != "" {
			frame++
			if frame == 5 {
				img := rl.LoadImageFromScreen()
				// On scaled (HiDPI) displays the read-back is larger than the rendered area.
				if sw, sh := int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()); img.Width > sw || img.Height > sh {
					rl.ImageCrop(img, rl.NewRectangle(0, float32(img.Height-sh), float32(sw), float32(sh)))
				}
				if !rl.ExportImage(*img, *shot) {
					log.Fatalf("could not write screenshot %s", *shot)
				}
				rl.UnloadImage(img)
				fmt.Fprintln(os.Stderr, "screenshot saved:", *shot)
				return
			}
		}
	}
}
