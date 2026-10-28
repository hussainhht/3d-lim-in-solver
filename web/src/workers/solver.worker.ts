/// <reference lib="webworker" />
// Each run owns a fresh runtime. Terminating the worker safely cancels Go as well.
export {};
declare const self: DedicatedWorkerGlobalScope & {
  Go: new () => {
    importObject: WebAssembly.Imports;
    run(instance: WebAssembly.Instance): Promise<void>;
  };
  lemInSolve: (text: string) => string;
};
self.onmessage = async ({
  data,
}: MessageEvent<{ text: string; base: string }>) => {
  try {
    self.postMessage({ status: "Loading Go WebAssembly…" });
    const runtime = await fetch(`${data.base}wasm_exec.js`);
    if (!runtime.ok)
      throw new Error("Unable to load the Go runtime. Refresh and try again.");
    // The official Go runtime is a classic script. A blob import works in module workers
    // without eval, including browsers that do not support importScripts in modules.
    const url = URL.createObjectURL(
      new Blob([await runtime.text()], { type: "text/javascript" }),
    );
    try {
      await import(/* @vite-ignore */ url);
    } finally {
      URL.revokeObjectURL(url);
    }
    const go = new self.Go();
    const response = await fetch(`${data.base}solver.wasm`);
    if (!response.ok)
      throw new Error("Unable to load solver.wasm. Refresh and try again.");
    const { instance } = await WebAssembly.instantiate(
      await response.arrayBuffer(),
      go.importObject,
    );
    void go
      .run(instance)
      .catch(() =>
        self.postMessage({
          error: "The Go runtime stopped unexpectedly. Please try again.",
        }),
      );
    self.postMessage({
      status: "Enumerating paths · selecting routes · simulating…",
    });
    self.postMessage(JSON.parse(self.lemInSolve(data.text)));
  } catch (e) {
    self.postMessage({
      error:
        e instanceof Error ? e.message : "Solver could not start. Try again.",
    });
  }
};
