import type { Solution, Step } from "../types";

export function buildTracks(solution: Solution): Step[][] {
  const tracks = Array.from({ length: solution.antCount }, () => [] as Step[]);
  solution.turns.forEach((moves, i) =>
    moves.forEach((m) =>
      tracks[m.AntID - 1].push({ turn: i + 1, room: m.Room }),
    ),
  );
  return tracks;
}

// The completed state is floor(time). Only the next genuine movement interpolates.
export function antState(track: Step[], start: string, time: number) {
  const turn = Math.floor(time);
  let lo = 0,
    hi = track.length;
  while (lo < hi) {
    const mid = (lo + hi) >>> 1;
    if (track[mid].turn <= turn) lo = mid + 1;
    else hi = mid;
  }
  const from = lo ? track[lo - 1].room : start;
  const next = track[lo];
  return {
    from,
    to: next?.turn === turn + 1 ? next.room : from,
    fraction: time - turn,
  };
}

export function arrivedAt(tracks: Step[][], end: string, time: number) {
  return tracks.filter(
    (t) =>
      t.length &&
      t[t.length - 1].room === end &&
      t[t.length - 1].turn <= Math.floor(time),
  ).length;
}
