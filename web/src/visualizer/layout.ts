import type { Solution } from "../types";
export type Point = [number, number, number];
export function layout(solution: Solution, radial = false): Map<string, Point> {
  const rooms = solution.rooms;
  const xs = rooms.map((r) => r.x),
    ys = rooms.map((r) => r.y);
  const minX = Math.min(...xs),
    maxX = Math.max(...xs),
    minY = Math.min(...ys),
    maxY = Math.max(...ys);
  const span = Math.max(maxX - minX, maxY - minY, 1);
  return new Map(
    rooms.map((r, i) => [
      r.name,
      radial
        ? [
            Math.cos((i / rooms.length) * Math.PI * 2) * 11,
            0,
            Math.sin((i / rooms.length) * Math.PI * 2) * 11,
          ]
        : ([
            ((r.x - minX) / span - (maxX - minX) / span / 2) * 24,
            0,
            ((r.y - minY) / span - (maxY - minY) / span / 2) * 24,
          ] as Point),
    ]),
  );
}
export function roomPoint(
  points: Map<string, Point>,
  name: string,
  id: number,
  solution: Solution,
): Point {
  const p = points.get(name)!;
  if (name !== solution.start && name !== solution.end)
    return [p[0], 0.72, p[2]];
  // Endpoint parking offsets are visual only; every logical position stays in its room.
  const a = id * 2.399963,
    r = 0.68 + 0.045 * Math.sqrt(id);
  return [p[0] + Math.cos(a) * r, 0.72, p[2] + Math.sin(a) * r];
}
