import {
  Component,
  Suspense,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import type { ReactNode } from "react";
import { Canvas, useFrame, useThree } from "@react-three/fiber";
import { Grid, Line, OrbitControls } from "@react-three/drei";
import { Color, InstancedMesh, Object3D, Vector3 } from "three";
import type { OrbitControls as OrbitImpl } from "three-stdlib";
import { palette } from "../types";
import type { Solution, Step } from "../types";
import { antState } from "../solver/playback";
import { layout, roomPoint } from "./layout";
import type { Point } from "./layout";

type Props = {
  solution: Solution;
  tracks: Step[][];
  time: number;
  reset: number;
  radial: boolean;
  labels: boolean;
  selected: string | null;
  onSelect: (name: string) => void;
  activePath: number | null;
};

function CameraRig({ reset }: { reset: number }) {
  const controls = useRef<OrbitImpl>(null);
  const { camera, size } = useThree();
  const home = useMemo(
    () =>
      new Vector3(13, 22, 25).multiplyScalar(
        Math.max(1, 1.15 / (size.width / size.height)),
      ),
    [size.width, size.height],
  );
  const settling = useRef(true);
  const reducedMotion = useMemo(
    () => window.matchMedia("(prefers-reduced-motion: reduce)").matches,
    [],
  );
  useEffect(() => {
    settling.current = true;
  }, [reset, home]);
  useFrame((_, delta) => {
    if (settling.current) {
      camera.position.lerp(home, reducedMotion ? 1 : 1 - Math.exp(-delta * 7));
      controls.current?.target.lerp(
        new Vector3(0, 0, 0),
        reducedMotion ? 1 : 1 - Math.exp(-delta * 7),
      );
      controls.current?.update();
      if (camera.position.distanceTo(home) < 0.03) settling.current = false;
    }
  });
  return (
    <OrbitControls
      ref={controls}
      makeDefault
      enableDamping={!reducedMotion}
      dampingFactor={0.1}
      minDistance={4}
      maxDistance={95}
      maxPolarAngle={Math.PI / 2 - 0.05}
      onStart={() => {
        settling.current = false;
      }}
    />
  );
}

function Ants({
  solution,
  tracks,
  time,
  points,
}: {
  solution: Solution;
  tracks: Step[][];
  time: number;
  points: Map<string, Point>;
}) {
  const abdomen = useRef<InstancedMesh>(null),
    thorax = useRef<InstancedMesh>(null),
    head = useRef<InstancedMesh>(null),
    legs = useRef<InstancedMesh>(null);
  const dummy = useMemo(() => new Object3D(), []);
  const colors = useMemo(
    () =>
      solution.assignments.map((p) => new Color(palette[p % palette.length])),
    [solution],
  );
  useEffect(() => {
    for (let i = 0; i < solution.antCount; i++) {
      abdomen.current?.setColorAt(i, colors[i]);
      thorax.current?.setColorAt(i, colors[i]);
      head.current?.setColorAt(i, colors[i]);
    }
    for (const ref of [abdomen, thorax, head])
      if (ref.current?.instanceColor)
        ref.current.instanceColor.needsUpdate = true;
  }, [solution, colors]);
  useFrame(() => {
    if (!abdomen.current || !thorax.current || !head.current || !legs.current)
      return;
    tracks.forEach((track, i) => {
      const s = antState(track, solution.start, time),
        from = roomPoint(points, s.from, i + 1, solution),
        to = roomPoint(points, s.to, i + 1, solution);
      const t = s.fraction * s.fraction * (3 - 2 * s.fraction),
        x = from[0] + (to[0] - from[0]) * t,
        z = from[2] + (to[2] - from[2]) * t;
      const yaw =
        s.from === s.to
          ? i * 2.399963
          : Math.atan2(to[0] - from[0], to[2] - from[2]);
      const sin = Math.sin(yaw),
        cos = Math.cos(yaw);
      for (const [mesh, offset, sx, sy, sz] of [
        [abdomen.current, -0.13, 0.12, 0.11, 0.19],
        [thorax.current, 0.06, 0.075, 0.075, 0.1],
        [head.current, 0.21, 0.095, 0.085, 0.09],
      ] as const) {
        dummy.position.set(x + sin * offset, 0.73, z + cos * offset);
        dummy.rotation.set(0, yaw, 0);
        dummy.scale.set(sx * 1.35, sy * 1.35, sz * 1.35);
        dummy.updateMatrix();
        mesh!.setMatrixAt(i, dummy.matrix);
      }
      for (let k = 0; k < 6; k++) {
        const side = k % 2 === 0 ? -1 : 1,
          forward = (Math.floor(k / 2) - 1) * 0.12;
        const stride =
          s.from === s.to ? 0 : Math.sin(time * 22 + i + k) * 0.025;
        dummy.position.set(
          x + cos * 0.13 * side + sin * forward,
          0.67,
          z - sin * 0.13 * side + cos * forward + stride,
        );
        dummy.rotation.set(0.2, yaw + side * 0.45, side * 0.4);
        dummy.scale.set(0.19, 0.018, 0.018);
        dummy.updateMatrix();
        legs.current!.setMatrixAt(i * 6 + k, dummy.matrix);
      }
    });
    for (const ref of [abdomen, thorax, head, legs])
      ref.current!.instanceMatrix.needsUpdate = true;
  });
  return (
    <group>
      {[abdomen, thorax, head].map((ref, i) => (
        <instancedMesh
          key={i}
          ref={ref}
          args={[undefined, undefined, solution.antCount]}
          frustumCulled={false}
        >
          <sphereGeometry args={[1, 10, 8]} />
          <meshStandardMaterial roughness={0.35} metalness={0.45} />
        </instancedMesh>
      ))}
      <instancedMesh
        ref={legs}
        args={[undefined, undefined, solution.antCount * 6]}
        frustumCulled={false}
      >
        <boxGeometry />
        <meshStandardMaterial color="#d8c99e" roughness={0.5} />
      </instancedMesh>
    </group>
  );
}

function Scene(
  props: Props & {
    points: Map<string, Point>;
    labelElements: Map<string, HTMLButtonElement>;
  },
) {
  const { solution, points, selected, onSelect, activePath } = props;
  const [hover, setHover] = useState<string | null>(null);
  const { camera, size } = useThree();
  const projected = useMemo(() => new Vector3(), []);
  useFrame(() => {
    for (const room of solution.rooms) {
      const element = props.labelElements.get(room.name);
      if (!element) continue;
      const p = points.get(room.name)!;
      projected
        .set(
          p[0],
          room.name === solution.start || room.name === solution.end
            ? 1.65
            : 0.95,
          p[2],
        )
        .project(camera);
      element.style.visibility =
        projected.z > 1 || projected.z < -1 ? "hidden" : "visible";
      element.style.transform = `translate(${(projected.x * 0.5 + 0.5) * size.width}px,${(-projected.y * 0.5 + 0.5) * size.height}px) translate(-50%,-50%)`;
    }
  });
  const edges = useMemo(() => {
    const map = new Map<string, number>();
    solution.paths.forEach((p, i) =>
      p
        .slice(1)
        .forEach((name, j) => map.set([p[j], name].sort().join("\0"), i)),
    );
    return map;
  }, [solution]);
  return (
    <>
      <color attach="background" args={["#0a131e"]} />
      <fog attach="fog" args={["#0a131e", 42, 95]} />
      <ambientLight intensity={1.3} />
      <hemisphereLight args={["#c5e6f5", "#172c30", 1.5]} />
      <directionalLight
        position={[4, 15, 10]}
        intensity={2.5}
        color="#b5e7e3"
      />
      <pointLight position={[-12, 8, -10]} intensity={30} color="#6cbda9" />
      <Grid
        position={[0, -0.16, 0]}
        args={[70, 70]}
        cellSize={1.5}
        sectionSize={7.5}
        cellThickness={0.45}
        sectionThickness={0.6}
        cellColor="#213345"
        sectionColor="#36505f"
        fadeDistance={52}
        fadeStrength={1.5}
      />
      <mesh position={[0, -0.24, 0]} rotation={[-Math.PI / 2, 0, 0]}>
        <circleGeometry args={[29, 80]} />
        <meshStandardMaterial color="#0e1c29" roughness={1} />
      </mesh>
      <mesh rotation={[-Math.PI / 2, 0, 0]} position={[0, -0.12, 0]}>
        <ringGeometry args={[16.8, 16.84, 100]} />
        <meshBasicMaterial color="#29404b" transparent opacity={0.65} />
      </mesh>
      {solution.tunnels.map(([a, b]) => {
        const path = edges.get([a, b].sort().join("\0")),
          active =
            path !== undefined && (activePath === null || path === activePath);
        const p = points.get(a)!,
          q = points.get(b)!;
        return (
          <group key={`${a}-${b}`}>
            <Line
              points={[
                [p[0], 0.3, p[2]],
                [q[0], 0.3, q[2]],
              ]}
              color={active ? palette[path! % palette.length] : "#334453"}
              lineWidth={active ? 2.3 : 1}
              transparent
              opacity={active ? 0.65 : 0.38}
            />
            {active && (
              <Line
                points={[
                  [p[0], 0.38, p[2]],
                  [q[0], 0.38, q[2]],
                ]}
                color={palette[path! % palette.length]}
                lineWidth={0.6}
                transparent
                opacity={0.35}
              />
            )}
          </group>
        );
      })}
      {solution.rooms.map((room) => {
        const p = points.get(room.name)!,
          start = room.name === solution.start,
          end = room.name === solution.end,
          special = start || end;
        const color = start
          ? "#62dec2"
          : end
            ? "#e9b666"
            : selected === room.name
              ? "#ffffff"
              : "#7193aa";
        return (
          <group key={room.name} position={p}>
            <mesh position={[0, 0.07, 0]}>
              <cylinderGeometry
                args={[special ? 0.72 : 0.42, special ? 0.85 : 0.52, 0.32, 6]}
              />
              <meshStandardMaterial
                color="#263b48"
                metalness={0.65}
                roughness={0.4}
              />
            </mesh>
            <mesh
              position={[0, 0.27, 0]}
              onPointerOver={(e) => {
                e.stopPropagation();
                setHover(room.name);
              }}
              onPointerOut={() => setHover(null)}
              onClick={(e) => {
                e.stopPropagation();
                onSelect(room.name);
              }}
            >
              <cylinderGeometry
                args={[special ? 0.57 : 0.3, special ? 0.57 : 0.3, 0.1, 6]}
              />
              <meshStandardMaterial
                color={color}
                emissive={color}
                emissiveIntensity={hover === room.name ? 0.65 : 0.22}
                roughness={0.25}
                metalness={0.5}
              />
            </mesh>
            {special && (
              <mesh rotation={[-Math.PI / 2, 0, 0]} position={[0, 0.15, 0]}>
                <ringGeometry args={[1.02, 1.05, 50]} />
                <meshBasicMaterial color={color} transparent opacity={0.5} />
              </mesh>
            )}
          </group>
        );
      })}
      <Ants
        solution={solution}
        tracks={props.tracks}
        time={props.time}
        points={points}
      />
      <CameraRig reset={props.reset} />
    </>
  );
}

function Fallback(props: Props & { points: Map<string, Point> }) {
  const { solution, points, time, tracks } = props;
  return (
    <div className="fallback">
      <p>3D graphics unavailable · accessible graph view</p>
      <svg
        viewBox="-15 -15 30 30"
        role="img"
        aria-label="Colony graph with real ant positions"
      >
        {solution.tunnels.map(([a, b]) => {
          const p = points.get(a)!,
            q = points.get(b)!;
          return (
            <line
              key={`${a}-${b}`}
              x1={p[0]}
              y1={p[2]}
              x2={q[0]}
              y2={q[2]}
              stroke="#456474"
              strokeWidth=".09"
            />
          );
        })}
        {solution.rooms.map((r) => {
          const p = points.get(r.name)!;
          return (
            <g key={r.name}>
              <circle
                cx={p[0]}
                cy={p[2]}
                r=".28"
                fill={r.name === solution.end ? "#e9b666" : "#62dec2"}
              />
              <text x={p[0] + 0.4} y={p[2] - 0.3} fill="#d1e0e6" fontSize=".5">
                {r.name}
              </text>
            </g>
          );
        })}
        {tracks.map((track, i) => {
          const s = antState(track, solution.start, time),
            p = roomPoint(points, s.from, i + 1, solution),
            q = roomPoint(points, s.to, i + 1, solution);
          return (
            <circle
              key={i}
              cx={p[0] + (q[0] - p[0]) * s.fraction}
              cy={p[2] + (q[2] - p[2]) * s.fraction}
              r=".11"
              fill={palette[solution.assignments[i] % palette.length]}
            />
          );
        })}
      </svg>
    </div>
  );
}

class GraphicsBoundary extends Component<
  { children: ReactNode; fallback: ReactNode },
  { failed: boolean }
> {
  state = { failed: false };
  static getDerivedStateFromError() {
    return { failed: true };
  }
  render() {
    return this.state.failed ? this.props.fallback : this.props.children;
  }
}
export default function Colony(props: Props) {
  const points = useMemo(
    () => layout(props.solution, props.radial),
    [props.solution, props.radial],
  );
  const labelElements = useRef(new Map<string, HTMLButtonElement>()).current;
  const fallback = <Fallback {...props} points={points} />;
  return (
    <GraphicsBoundary fallback={fallback}>
      <div className="colony-canvas">
        <Canvas
          dpr={[1, 1.5]}
          camera={{ position: [15, 26, 29], fov: 45, near: 0.1, far: 150 }}
          gl={{ antialias: true, powerPreference: "low-power" }}
          aria-label="Interactive 3D colony. Drag to orbit, scroll to zoom, right-drag to pan."
        >
          <Suspense fallback={null}>
            <Scene {...props} points={points} labelElements={labelElements} />
          </Suspense>
        </Canvas>
        <div className="projected-labels">
          {props.solution.rooms.map((room) => {
            const start = room.name === props.solution.start,
              end = room.name === props.solution.end;
            return (
              (props.labels ||
                start ||
                end ||
                props.selected === room.name) && (
                <button
                  key={room.name}
                  ref={(el) => {
                    if (el) labelElements.set(room.name, el);
                    else labelElements.delete(room.name);
                  }}
                  className={`room-label ${start ? "start" : end ? "end" : ""} ${props.selected === room.name ? "chosen" : ""}`}
                  onClick={() => props.onSelect(room.name)}
                  aria-label={`Inspect room ${room.name}`}
                  title={`${room.name} · (${room.x}, ${room.y}) · ${props.solution.tunnels.filter((t) => t.includes(room.name)).length} tunnels`}
                >
                  {(start || end) && <span>{start ? "START" : "END"}</span>}
                  {room.name}
                </button>
              )
            );
          })}
        </div>
      </div>
    </GraphicsBoundary>
  );
}
