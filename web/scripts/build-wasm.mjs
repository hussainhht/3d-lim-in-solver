import { execFileSync } from "node:child_process";
import { copyFileSync, mkdirSync, existsSync } from "node:fs";
import { fileURLToPath } from "node:url";
import { resolve } from "node:path";
const root = fileURLToPath(new URL("../../", import.meta.url));
const publicDir = resolve(root, "web/public");
mkdirSync(publicDir, { recursive: true });
execFileSync(
  "go",
  ["build", "-trimpath", "-o", resolve(publicDir, "solver.wasm"), "./wasm"],
  {
    cwd: root,
    env: { ...process.env, GOOS: "js", GOARCH: "wasm", CGO_ENABLED: "0" },
    stdio: "inherit",
  },
);
const goroot = execFileSync("go", ["env", "GOROOT"], {
  encoding: "utf8",
}).trim();
const runtime = [
  resolve(goroot, "lib/wasm/wasm_exec.js"),
  resolve(goroot, "misc/wasm/wasm_exec.js"),
].find(existsSync);
if (!runtime) throw new Error("Matching Go WebAssembly runtime not found");
copyFileSync(runtime, resolve(publicDir, "wasm_exec.js"));
console.log("Built real Go solver and copied matching wasm_exec.js.");
