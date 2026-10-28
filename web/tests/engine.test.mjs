import { test, before, after } from "node:test";
import assert from "node:assert/strict";
import { Worker } from "node:worker_threads";
import { execFileSync } from "node:child_process";
import { readFileSync, mkdtempSync } from "node:fs";
import { tmpdir } from "node:os";
import { join, resolve } from "node:path";
import { once } from "node:events";
import { antState, buildTracks, arrivedAt } from "../src/solver/playback.ts";
const cwd = mkdtempSync(join(tmpdir(), "lem-in-web-tests-"));
const root = resolve("..");
const suffix = process.platform === "win32" ? ".exe" : "";
const native = join(cwd, `native${suffix}`),
  cli = join(cwd, `cli${suffix}`);
let worker;
before(async () => {
  const options = { cwd: root, env: { ...process.env, CGO_ENABLED: "0" } };
  execFileSync("go", ["build", "-o", native, "./web/tests/native"], options);
  execFileSync("go", ["build", "-o", cli, "."], options);
  worker = new Worker(new URL("./wasm-runner.mjs", import.meta.url));
  assert.equal((await once(worker, "message"))[0].ready, true);
});
after(async () => {
  await worker?.terminate();
});
async function solve(text) {
  worker.postMessage(text);
  return (await once(worker, "message"))[0];
}

for (const file of [
  "input.txt",
  "3D/input.txt",
  ...Array.from({ length: 8 }, (_, i) => `test/example0${i}.txt`),
]) {
  test(`WASM matches native paths, assignments and exact CLI turns: ${file}`, async () => {
    const text = readFileSync(resolve(root, file), "utf8");
    const { result: web, error } = await solve(text);
    assert.equal(error, undefined);
    const oracle = JSON.parse(
      execFileSync(native, [], { input: text, encoding: "utf8" }),
    );
    const { computeMs: _a, ...expected } = oracle,
      { computeMs: _b, ...actual } = web;
    assert.deepEqual(actual, expected);
    const out = execFileSync(cli, [resolve(root, file)], {
      cwd,
      encoding: "utf8",
    });
    const lines = out
      .replaceAll("\r", "")
      .split("\n")
      .filter((l) => /^L\d+-/.test(l));
    assert.deepEqual(web.output, lines);
    const tracks = buildTracks(web);
    assert.equal(arrivedAt(tracks, web.end, 0), 0);
    assert.equal(arrivedAt(tracks, web.end, web.totalTurns), web.antCount);
    for (let turn = 0; turn <= web.totalTurns; turn++) {
      const pos = Array(web.antCount).fill(web.start);
      web.turns
        .slice(0, turn)
        .flat()
        .forEach((m) => (pos[m.AntID - 1] = m.Room));
      for (let i = 0; i < web.antCount; i++) {
        assert.equal(antState(tracks[i], web.start, turn).from, pos[i]);
        const mid = antState(tracks[i], web.start, turn + 0.5);
        const next =
          web.turns[turn]?.find((m) => m.AntID === i + 1)?.Room ?? pos[i];
        assert.equal(mid.to, next);
      }
    }
  });
}
test("WASM rejects malformed and disconnected maps and recovers", async () => {
  for (const text of [
    "",
    "0\n",
    "3\n##start\ns 0 0\n##end\ne 1 0\n",
    "1000001\n",
  ]) {
    assert.equal(typeof (await solve(text)).error, "string");
  }
  assert.equal(
    (await solve(readFileSync(resolve(root, "input.txt"), "utf8"))).result
      .totalTurns,
    6,
  );
});
