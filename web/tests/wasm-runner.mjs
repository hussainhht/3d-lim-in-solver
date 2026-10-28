import { parentPort } from "node:worker_threads";
import { readFile } from "node:fs/promises";
await import("../public/wasm_exec.js");
const go = new globalThis.Go();
const { instance } = await WebAssembly.instantiate(
  await readFile(new URL("../public/solver.wasm", import.meta.url)),
  go.importObject,
);
void go.run(instance);
parentPort.on("message", (text) =>
  parentPort.postMessage(JSON.parse(globalThis.lemInSolve(text))),
);
parentPort.postMessage({ ready: true });
