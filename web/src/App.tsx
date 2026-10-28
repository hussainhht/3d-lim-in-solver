import {
  lazy,
  Suspense,
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import {
  ArrowDownToLine,
  ArrowRight,
  Box,
  Check,
  CheckCircle2,
  ChevronDown,
  CircleHelp,
  Code2,
  Copy,
  Cpu,
  FileUp,
  Focus,
  GitBranch,
  Github,
  Layers3,
  LoaderCircle,
  Maximize2,
  Network,
  Pause,
  Play,
  RotateCcw,
  Route,
  ScanLine,
  Settings2,
  ShieldCheck,
  SkipBack,
  SkipForward,
  Terminal,
  Timer,
  Trash2,
  X,
} from "lucide-react";
import type { LucideIcon } from "lucide-react";
import { examples } from "./examples";
import { useSolver } from "./hooks/useSolver";
import { antState, arrivedAt, buildTracks } from "./solver/playback";
import { palette } from "./types";
const Colony = lazy(() => import("./visualizer/Colony"));

function Mark() {
  return (
    <svg viewBox="0 0 40 40" fill="none" aria-hidden="true">
      <path
        d="m20 4 14 8v16l-14 8-14-8V12L20 4Z"
        stroke="currentColor"
        strokeWidth="1.4"
      />
      <path
        d="m20 12 7 4v8l-7 4-7-4v-8l7-4Zm0-8v8m14 0-7 4m7 12-7-4M20 36v-8M6 28l7-4M6 12l7 4"
        stroke="currentColor"
        strokeWidth="1.4"
      />
      <circle cx="20" cy="20" r="3" fill="currentColor" />
    </svg>
  );
}
function IconButton({
  icon: Icon,
  label,
  onClick,
  disabled = false,
  active = false,
}: {
  icon: LucideIcon;
  label: string;
  onClick: () => void;
  disabled?: boolean;
  active?: boolean;
}) {
  return (
    <button
      className={`icon-button ${active ? "active" : ""}`}
      type="button"
      onClick={onClick}
      disabled={disabled}
      title={label}
      aria-label={label}
    >
      <Icon size={16} />
    </button>
  );
}
function Metric({
  label,
  value,
  icon: Icon,
}: {
  label: string;
  value: string | number;
  icon: LucideIcon;
}) {
  return (
    <div className="metric">
      <span>
        <Icon size={13} />
        {label}
      </span>
      <strong
        data-testid={`metric-${label.toLowerCase().replaceAll(" ", "-")}`}
      >
        {value}
      </strong>
    </div>
  );
}

export default function App() {
  const [exampleId, setExampleId] = useState(examples[0].id);
  const example = examples.find((e) => e.id === exampleId);
  const [text, setText] = useState(examples[0].text);
  const [time, setTime] = useState(0),
    [playing, setPlaying] = useState(false),
    [speed, setSpeed] = useState(1);
  const [budget, setBudget] = useState(15),
    [camera, setCamera] = useState(0),
    [radial, setRadial] = useState(false),
    [labels, setLabels] = useState(true);
  const [selected, setSelected] = useState<string | null>(null),
    [activePath, setActivePath] = useState<number | null>(null);
  const [notice, setNotice] = useState(""),
    [inspector, setInspector] = useState<"output" | "ants">("output");
  const [antId, setAntId] = useState(1);
  const fileInput = useRef<HTMLInputElement>(null);
  const editorLines = useRef<HTMLDivElement>(null);
  const { solution, status, error, busy, solve, clear, cancel } = useSolver();
  const tracks = useMemo(
    () => (solution ? buildTracks(solution) : []),
    [solution],
  );
  const total = solution?.totalTurns ?? 0;
  const turn = Math.floor(time);
  const arrived = solution ? arrivedAt(tracks, solution.end, time) : 0;
  const complete = !!solution && time >= total;
  const resetPlayback = useCallback(() => {
    setTime(0);
    setPlaying(false);
    setSelected(null);
    setActivePath(null);
    setAntId(1);
  }, []);
  const run = useCallback(() => {
    resetPlayback();
    solve(text, budget);
  }, [text, budget, solve, resetPlayback]);
  useEffect(() => {
    solve(examples[0].text, 15);
  }, [solve]);
  useEffect(() => {
    resetPlayback();
    setCamera((c) => c + 1);
  }, [solution, resetPlayback]);
  useEffect(() => {
    if (!playing || !solution) return;
    let frame = 0,
      last = performance.now();
    const reduced = window.matchMedia(
      "(prefers-reduced-motion: reduce)",
    ).matches;
    let accumulated = 0;
    const tick = (now: number) => {
      const delta = Math.min((now - last) / 1000, 0.1) * speed;
      last = now;
      accumulated += delta;
      if (!reduced || accumulated >= 1) {
        const step = reduced ? Math.floor(accumulated) : delta;
        if (reduced) accumulated -= step;
        setTime((t) => Math.min(total, t + step));
      }
      frame = requestAnimationFrame(tick);
    };
    frame = requestAnimationFrame(tick);
    return () => cancelAnimationFrame(frame);
  }, [playing, solution, speed, total]);
  useEffect(() => {
    if (complete) setPlaying(false);
  }, [complete]);
  const jump = useCallback(
    (value: number) => {
      setPlaying(false);
      setTime(Math.max(0, Math.min(total, value)));
    },
    [total],
  );
  const toggle = useCallback(() => {
    if (!solution) return;
    if (time >= total) setTime(0);
    setPlaying((p) => !p);
  }, [solution, time, total]);
  useEffect(() => {
    const key = (e: KeyboardEvent) => {
      if (
        (e.target as HTMLElement).closest(
          "input,textarea,select,button,a,summary",
        )
      )
        return;
      if (e.code === "Space") {
        e.preventDefault();
        toggle();
      }
      if (e.code === "ArrowRight") {
        e.preventDefault();
        jump(Math.floor(time) + 1);
      }
      if (e.code === "ArrowLeft") {
        e.preventDefault();
        jump(Math.ceil(time) - 1);
      }
    };
    window.addEventListener("keydown", key);
    return () => window.removeEventListener("keydown", key);
  }, [toggle, jump, time]);
  const edit = (value: string) => {
    clear();
    resetPlayback();
    setText(value);
    setNotice("");
  };
  const choose = (id: string) => {
    const e = examples.find((e) => e.id === id)!;
    setExampleId(id);
    edit(e.text);
  };
  const upload = async (file?: File) => {
    if (!file) return;
    if (!file.name.toLowerCase().endsWith(".txt")) {
      setNotice("Choose a Lem-in map with a .txt extension.");
      return;
    }
    if (file.size > 128 * 1024) {
      setNotice("The browser accepts maps up to 128 KiB.");
      return;
    }
    try {
      edit(await file.text());
      setExampleId("custom");
      setNotice(`Loaded ${file.name}`);
    } catch {
      setNotice("This file could not be read. Please try again.");
    }
  };
  const download = () => {
    const url = URL.createObjectURL(new Blob([text], { type: "text/plain" }));
    const a = document.createElement("a");
    a.href = url;
    a.download = "colony.txt";
    a.click();
    URL.revokeObjectURL(url);
  };
  const copy = async (value: string) => {
    try {
      await navigator.clipboard.writeText(value);
      setNotice("Copied to clipboard.");
    } catch {
      setNotice("Clipboard unavailable. Select the text to copy it.");
    }
  };
  const antLine = text
    .split(/\r?\n/)
    .findIndex((l) => l.trim() && !l.trim().startsWith("#"));
  const antCount =
    antLine >= 0 && /^\d+$/.test(text.split(/\r?\n/)[antLine].trim())
      ? Number(text.split(/\r?\n/)[antLine])
      : "";
  const setAntCount = (value: string) => {
    const lines = text.split(/\r?\n/);
    if (antLine >= 0) lines[antLine] = value;
    else lines.unshift(value);
    edit(lines.join("\n"));
  };
  const room = solution?.rooms.find((r) => r.name === selected);
  const currentAnt = solution
    ? antState(tracks[antId - 1] ?? [], solution.start, time)
    : null;
  return (
    <>
      <a className="skip-link" href="#playground">
        Skip to interactive demo
      </a>
      <header className="header">
        <a href="#" className="brand">
          <div className="brand-mark">
            <Mark />
          </div>
          <span>
            lem<span className="brand-hyphen">-</span>in
            <span className="brand-sub">COLONY INTELLIGENCE</span>
          </span>
        </a>
        <nav aria-label="Main navigation">
          <a href="#playground" className="nav-active">
            Playground
          </a>
          <a
            href="https://github.com/hussainhht/3d-lim-in-solver"
            target="_blank"
            rel="noreferrer"
          >
            <Github size={16} /> Source code <ArrowRight size={13} />
          </a>
        </nav>
        <span className="live-badge">
          <i /> Interactive demo
        </span>
      </header>
      <main>
        <div className="workspace-heading" id="playground">
          <div>
            <span className="section-index">01</span>
            <h2>Colony playground</h2>
          </div>
          <span className="privacy-note">
            <ShieldCheck size={13} /> Your maps stay on your device
          </span>
        </div>
        <div className="workspace">
          <aside className="panel config-panel" aria-label="Map configuration">
            <div className="panel-title">
              <span>
                <Settings2 size={15} /> Map configuration
              </span>
              <span className="tiny-dot" />
            </div>
            <div className="config-body">
              <label className="field-label" htmlFor="example">
                EXPLORE A COLONY
              </label>
              <div className="select-wrap">
                <Layers3 size={16} />
                <select
                  id="example"
                  value={exampleId}
                  onChange={(e) => choose(e.target.value)}
                >
                  {!example && <option value="custom">Custom colony</option>}
                  {examples.map((e) => (
                    <option key={e.id} value={e.id}>
                      {e.name}
                    </option>
                  ))}
                </select>
                <ChevronDown size={14} />
              </div>
              <p className="map-description">
                {example?.description ??
                  "Your own colony. Edit the map below, then run the Go solver."}
              </p>
              <div className="source-caption">
                <span>{example?.tag ?? "CUSTOM MAP"}</span>
                <code>{example?.file ?? "colony.txt"}</code>
              </div>
              <div className="editor-heading">
                <label htmlFor="map-editor">
                  <Code2 size={14} /> Map definition
                </label>
                <div>
                  <IconButton
                    icon={RotateCcw}
                    label="Reset to example"
                    onClick={() => choose(example?.id ?? examples[0].id)}
                  />
                  <IconButton
                    icon={Trash2}
                    label="Clear editor"
                    onClick={() => {
                      setExampleId("custom");
                      edit("");
                    }}
                  />
                </div>
              </div>
              <div className="editor">
                <div
                  ref={editorLines}
                  className="line-numbers"
                  aria-hidden="true"
                >
                  {text.split("\n").map((_, i) => (
                    <div key={i}>{i + 1}</div>
                  ))}
                </div>
                <textarea
                  id="map-editor"
                  maxLength={128 * 1024}
                  spellCheck={false}
                  value={text}
                  aria-label="Map definition"
                  onChange={(e) => edit(e.target.value)}
                  onScroll={(e) => {
                    if (editorLines.current)
                      editorLines.current.scrollTop = e.currentTarget.scrollTop;
                  }}
                />
              </div>
              <div className="file-actions">
                <button onClick={() => fileInput.current?.click()}>
                  <FileUp size={13} /> Upload .txt
                </button>
                <div>
                  <IconButton
                    icon={Copy}
                    label="Copy map"
                    onClick={() => void copy(text)}
                  />
                  <IconButton
                    icon={ArrowDownToLine}
                    label="Download map"
                    onClick={download}
                  />
                </div>
                <input
                  ref={fileInput}
                  type="file"
                  accept=".txt,text/plain"
                  aria-label="Upload map file"
                  hidden
                  onChange={(e) => {
                    void upload(e.target.files?.[0]);
                    e.target.value = "";
                  }}
                />
              </div>
              <div className="config-fields">
                <label>
                  Colony size
                  <input
                    aria-label="Ant count"
                    type="number"
                    min="1"
                    max="1000"
                    value={antCount}
                    onChange={(e) => setAntCount(e.target.value)}
                  />
                </label>
                <label>
                  Time budget
                  <select
                    aria-label="Solver time budget"
                    value={budget}
                    onChange={(e) => setBudget(Number(e.target.value))}
                  >
                    <option value={5}>5 seconds</option>
                    <option value={15}>15 seconds</option>
                    <option value={30}>30 seconds</option>
                  </select>
                </label>
              </div>
              <button
                className={`solve-button ${busy ? "cancel" : ""}`}
                onClick={busy ? cancel : run}
              >
                {busy ? (
                  <>
                    <X size={16} /> Cancel computation
                  </>
                ) : (
                  <>
                    <Play size={15} fill="currentColor" /> Run solver{" "}
                    <span>↗</span>
                  </>
                )}
              </button>
              <div
                className={`validation ${error ? "invalid" : solution ? "valid" : ""}`}
                role="status"
                aria-live="polite"
              >
                {busy ? (
                  <LoaderCircle className="spin" size={13} />
                ) : solution ? (
                  <CheckCircle2 size={13} />
                ) : (
                  <span className="status-dot" />
                )}
                <span>{status}</span>
              </div>
              {error && (
                <p className="error" role="alert">
                  {error}
                </p>
              )}
              {notice && (
                <p className="notice" role="status">
                  {notice}
                </p>
              )}
              <details className="format-help">
                <summary>
                  <CircleHelp size={13} /> Map format & browser limits
                </summary>
                <p>
                  Positive ant count, then <code>##start</code> /{" "}
                  <code>##end</code> before their rooms. Rooms:{" "}
                  <code>name x y</code>. Tunnels: <code>a-b</code>, after all
                  rooms. Names and integer coordinates must be unique. No
                  self-loops or duplicate tunnels.
                </p>
                <p>
                  Up to 1,000 ants, 80 rooms, 160 tunnels, 128 KiB. Dense graphs
                  can hit the search or time budget. The Go parser validates
                  when you run.
                </p>
              </details>
            </div>
          </aside>
          <section
            className="observatory panel"
            aria-label="Colony visualization"
          >
            <div className="panel-title">
              <span>
                <Box size={15} /> Colony observatory
              </span>
              <span className={`scene-status ${complete ? "finished" : ""}`}>
                <i />
                {busy
                  ? "COMPUTING"
                  : complete
                    ? "COLONY ARRIVED"
                    : playing
                      ? "IN MOTION"
                      : solution
                        ? "READY TO EXPLORE"
                        : "AWAITING MAP"}
              </span>
            </div>
            <div className="scene">
              {solution ? (
                <Suspense
                  fallback={
                    <div className="scene-empty">
                      <LoaderCircle className="spin" />
                      <p>Preparing the observatory…</p>
                    </div>
                  }
                >
                  <Colony
                    solution={solution}
                    tracks={tracks}
                    time={time}
                    reset={camera}
                    radial={radial}
                    labels={labels}
                    selected={selected}
                    onSelect={setSelected}
                    activePath={activePath}
                  />
                </Suspense>
              ) : (
                <div className="scene-empty">
                  <div className="empty-mark">
                    <Mark />
                  </div>
                  {busy ? (
                    <>
                      <LoaderCircle className="spin" size={24} />
                      <h3>Finding a way forward.</h3>
                      <p>The real Go solver is working on your colony.</p>
                    </>
                  ) : (
                    <>
                      <h3>Your next expedition starts here.</h3>
                      <p>
                        {error
                          ? "Review your map, then try again."
                          : "Choose a map or create your own, then run the solver."}
                      </p>
                      <button className="secondary-button" onClick={run}>
                        Run this colony <ArrowRight size={14} />
                      </button>
                    </>
                  )}
                </div>
              )}
              <div className="scene-top">
                <span className="scene-chip">
                  <span className="tiny-dot" />{" "}
                  {radial ? "RADIAL LAYOUT" : "MAP COORDINATES"}
                </span>
                <div className="scene-tools">
                  <IconButton
                    icon={ScanLine}
                    label={labels ? "Hide room labels" : "Show room labels"}
                    active={labels}
                    disabled={!solution}
                    onClick={() => setLabels((l) => !l)}
                  />
                  <IconButton
                    icon={Network}
                    label="Toggle coordinate or radial layout"
                    active={radial}
                    disabled={!solution}
                    onClick={() => {
                      setRadial((r) => !r);
                      setCamera((c) => c + 1);
                    }}
                  />
                  <IconButton
                    icon={RotateCcw}
                    label="Reset camera"
                    disabled={!solution}
                    onClick={() => setCamera((c) => c + 1)}
                  />
                  <IconButton
                    icon={Maximize2}
                    label="Fit map to view"
                    disabled={!solution}
                    onClick={() => setCamera((c) => c + 1)}
                  />
                </div>
              </div>
              <div className="scene-bottom">
                <span>
                  <span className="coordinate-axis">
                    Y<span>↟</span>
                  </span>{" "}
                  XZ / colony plane
                </span>
                <span>
                  Drag to orbit <b>·</b> Scroll to zoom
                </span>
              </div>
              {room && (
                <div className="room-popover">
                  <div>
                    <Focus size={13} />
                    <strong>{room.name}</strong>
                    <IconButton
                      icon={X}
                      label="Close room information"
                      onClick={() => setSelected(null)}
                    />
                  </div>
                  <p>
                    Coordinates {room.x}, {room.y} ·{" "}
                    {
                      solution!.tunnels.filter((t) => t.includes(room.name))
                        .length
                    }{" "}
                    tunnels
                  </p>
                  <p>
                    {
                      tracks.filter(
                        (t) =>
                          antState(t, solution!.start, time).from === room.name,
                      ).length
                    }{" "}
                    ants after turn {turn}
                  </p>
                </div>
              )}
            </div>
            <div className="playback">
              <div className="playback-top">
                <div className="transport">
                  <IconButton
                    icon={RotateCcw}
                    label="Restart simulation"
                    onClick={() => jump(0)}
                    disabled={!solution}
                  />
                  <IconButton
                    icon={SkipBack}
                    label="Previous turn"
                    onClick={() => jump(Math.ceil(time) - 1)}
                    disabled={!solution || time <= 0}
                  />
                  <button
                    className="play-button"
                    aria-label={
                      playing ? "Pause simulation" : "Play simulation"
                    }
                    onClick={toggle}
                    disabled={!solution}
                  >
                    {playing ? (
                      <Pause size={17} fill="currentColor" />
                    ) : (
                      <Play size={17} fill="currentColor" />
                    )}
                  </button>
                  <IconButton
                    icon={SkipForward}
                    label="Next turn"
                    onClick={() => jump(Math.floor(time) + 1)}
                    disabled={!solution || complete}
                  />
                </div>
                <span className="turn-count">
                  TURN{" "}
                  <strong data-testid="current-turn">
                    {turn.toString().padStart(2, "0")}
                  </strong>
                  <span>/ {total.toString().padStart(2, "0")}</span>
                </span>
                <label className="speed">
                  SPEED
                  <select
                    aria-label="Playback speed"
                    value={speed}
                    onChange={(e) => setSpeed(Number(e.target.value))}
                  >
                    <option value={0.25}>0.25×</option>
                    <option value={0.5}>0.5×</option>
                    <option value={1}>1×</option>
                    <option value={2}>2×</option>
                    <option value={4}>4×</option>
                    <option value={8}>8×</option>
                  </select>
                </label>
              </div>
              <input
                className="timeline"
                aria-label="Simulation timeline"
                type="range"
                min="0"
                max={total || 1}
                step="0.01"
                value={time}
                disabled={!solution}
                onChange={(e) => jump(Number(e.target.value))}
                style={
                  {
                    "--progress": `${total ? (time / total) * 100 : 0}%`,
                  } as React.CSSProperties
                }
              />
              <div className="timeline-caption">
                <span>
                  {complete ? (
                    <>
                      <Check size={12} /> Every ant has arrived.
                    </>
                  ) : playing ? (
                    "Following the solver’s next move."
                  ) : (
                    "A journey, one turn at a time."
                  )}
                </span>
                <span>SPACE play / pause</span>
              </div>
            </div>
          </section>
          <aside className="panel telemetry" aria-label="Colony statistics">
            <div className="panel-title">
              <span>
                <Cpu size={15} /> Colony telemetry
              </span>
            </div>
            <div className="telemetry-body">
              <div className="metrics">
                <Metric
                  label="Total ants"
                  value={solution?.antCount ?? "—"}
                  icon={GitBranch}
                />
                <Metric
                  label="Rooms"
                  value={solution?.rooms.length ?? "—"}
                  icon={Box}
                />
                <Metric
                  label="Tunnels"
                  value={solution?.tunnels.length ?? "—"}
                  icon={Network}
                />
                <Metric
                  label="Paths"
                  value={solution?.paths.length ?? "—"}
                  icon={Route}
                />
              </div>
              <div className="arrival-card">
                <div>
                  <span>ARRIVED AT DESTINATION</span>
                  <CheckCircle2 size={14} />
                </div>
                <strong data-testid="arrived">
                  {arrived}
                  <span> / {solution?.antCount ?? "—"}</span>
                </strong>
                <div className="arrival-bar">
                  <span
                    style={{
                      width: `${solution ? (arrived / solution.antCount) * 100 : 0}%`,
                    }}
                  />
                </div>
                <p>
                  {complete
                    ? "Expedition complete."
                    : solution
                      ? `${solution.antCount - arrived} ants on their journey`
                      : "Waiting for a solution"}
                </p>
              </div>
              <div className="timing">
                <span>
                  <Timer size={13} /> Total turns{" "}
                  <b data-testid="total-turns">{total || "—"}</b>
                </span>
                <span>
                  <Cpu size={13} /> Computation{" "}
                  <b>
                    {solution ? `${solution.computeMs.toFixed(1)} ms` : "—"}
                  </b>
                </span>
              </div>
              <div className="routes-title">
                <Route size={14} /> Selected routes{" "}
                <span>{solution?.paths.length ?? 0}</span>
              </div>
              <p className="route-hint">Select a route to trace its tunnels.</p>
              <div className="routes">
                {solution?.paths.map((path, i) => (
                  <button
                    className={`route-card ${activePath === i ? "selected" : ""}`}
                    key={i}
                    onClick={() => setActivePath((p) => (p === i ? null : i))}
                    aria-pressed={activePath === i}
                    style={
                      {
                        "--route": palette[i % palette.length],
                      } as React.CSSProperties
                    }
                  >
                    <div>
                      <span>
                        <i /> ROUTE {String(i + 1).padStart(2, "0")}
                      </span>
                      <b>
                        {solution.assignments.filter((p) => p === i).length}{" "}
                        ants
                      </b>
                    </div>
                    <p>{path.join(" → ")}</p>
                    <span>{path.length - 1} tunnels</span>
                  </button>
                ))}
                {!solution && (
                  <div className="no-routes">
                    Computed routes will appear here.
                  </div>
                )}
              </div>
            </div>
          </aside>
        </div>
        <section className="panel output-panel" aria-label="Output inspector">
          <div className="output-header">
            <div
              className="output-tabs"
              role="tablist"
              aria-label="Inspector view"
            >
              <button
                role="tab"
                aria-selected={inspector === "output"}
                onClick={() => setInspector("output")}
              >
                <Terminal size={15} /> Solver output
              </button>
              <button
                role="tab"
                aria-selected={inspector === "ants"}
                onClick={() => setInspector("ants")}
              >
                <GitBranch size={15} /> Ant inspector
              </button>
            </div>
            <span>
              <span className="tiny-dot" />{" "}
              {solution ? "REAL COMPUTED MOVEMENTS" : "NO SOLUTION LOADED"}
            </span>
          </div>
          {inspector === "output" ? (
            <div
              className="output-content"
              role="tabpanel"
              aria-label="Solver output"
            >
              <div className="output-turn">
                {turn ? "TURN" : "STATE"}{" "}
                <b>{turn.toString().padStart(2, "0")}</b>
              </div>
              <pre data-testid="turn-output">
                {solution
                  ? turn
                    ? solution.output[turn - 1]
                    : "All ants are waiting in the start room."
                  : busy
                    ? "Computing in an isolated Web Worker…"
                    : "Run the solver to inspect the exact CLI output for each completed turn."}
              </pre>
              <IconButton
                icon={Copy}
                label="Copy all solver output"
                disabled={!solution}
                onClick={() => void copy(solution!.output.join("\n") + "\n")}
              />
            </div>
          ) : (
            <div
              className="ant-inspector"
              role="tabpanel"
              aria-label="Ant inspector"
            >
              {solution ? (
                <>
                  <label>
                    Ant ID{" "}
                    <input
                      aria-label="Inspect ant ID"
                      type="number"
                      min="1"
                      max={solution.antCount}
                      value={antId}
                      onChange={(e) =>
                        setAntId(
                          Math.max(
                            1,
                            Math.min(
                              solution.antCount,
                              Number(e.target.value) || 1,
                            ),
                          ),
                        )
                      }
                    />
                  </label>
                  <span>
                    <small>ASSIGNED ROUTE</small>
                    <b
                      style={{
                        color:
                          palette[
                            solution.assignments[antId - 1] % palette.length
                          ],
                      }}
                    >
                      Route {solution.assignments[antId - 1] + 1}
                    </b>
                  </span>
                  <span>
                    <small>CURRENT ROOM</small>
                    <b>{currentAnt!.from}</b>
                  </span>
                  <span>
                    <small>NEXT MOVEMENT</small>
                    <b>
                      {currentAnt!.from === currentAnt!.to
                        ? "At rest"
                        : `${currentAnt!.from} → ${currentAnt!.to}`}
                    </b>
                  </span>
                  <span>
                    <small>STATUS</small>
                    <b>
                      {currentAnt!.from === solution.end
                        ? "Arrived"
                        : currentAnt!.from === solution.start
                          ? "Waiting"
                          : "In transit"}
                    </b>
                  </span>
                </>
              ) : (
                <p>
                  Run a colony to inspect each ant’s actual assignment and
                  position.
                </p>
              )}
            </div>
          )}
        </section>
      </main>
      <footer>
        <div className="footer-brand">
          <Mark />
          <span>
            lem-in <b>/</b> a study in collective movement
          </span>
        </div>
        <span>Go + WebAssembly + Three.js</span>
        <a
          href="https://github.com/hussainhht/3d-lim-in-solver"
          target="_blank"
          rel="noreferrer"
        >
          Explore the source <ArrowRight size={13} />
        </a>
      </footer>
    </>
  );
}
