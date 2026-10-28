export interface Room {
  name: string;
  x: number;
  y: number;
}
export interface Move {
  AntID: number;
  Room: string;
}
export interface Solution {
  rooms: Room[];
  tunnels: [string, string][];
  start: string;
  end: string;
  antCount: number;
  paths: string[][];
  assignments: number[];
  turns: Move[][];
  output: string[];
  totalTurns: number;
  computeMs: number;
  engine: string;
}
export interface Step {
  turn: number;
  room: string;
}
export const palette = [
  "#62dec2",
  "#e9b666",
  "#aa98f0",
  "#74b9f4",
  "#ef8fab",
  "#b4d77d",
];
