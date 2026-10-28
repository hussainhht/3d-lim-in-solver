import { useCallback, useEffect, useRef, useState } from "react";
import type { Solution } from "../types";

export function useSolver() {
  const [solution, setSolution] = useState<Solution | null>(null);
  const [status, setStatus] = useState("Ready to explore");
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);
  const worker = useRef<Worker | null>(null);
  const timeout = useRef<ReturnType<typeof setTimeout> | undefined>(undefined);
  const stop = useCallback(() => {
    worker.current?.terminate();
    worker.current = null;
    clearTimeout(timeout.current);
  }, []);
  useEffect(() => stop, [stop]);
  const clear = useCallback(() => {
    stop();
    setBusy(false);
    setSolution(null);
    setError("");
    setStatus("Map changed · run to validate");
  }, [stop]);
  const cancel = useCallback(() => {
    stop();
    setBusy(false);
    setStatus("Cancelled · ready to run again");
  }, [stop]);
  const solve = useCallback(
    (text: string, seconds: number) => {
      stop();
      setSolution(null);
      setError("");
      setBusy(true);
      setStatus("Starting isolated solver…");
      try {
        const current = new Worker(
          new URL("../workers/solver.worker.ts", import.meta.url),
          { type: "module" },
        );
        worker.current = current;
        const fail = (message: string) => {
          stop();
          setBusy(false);
          setError(message);
          setStatus("Map could not be solved");
        };
        current.onmessage = ({ data }) => {
          if (data.status) setStatus(data.status);
          else if (data.error) fail(data.error);
          else if (data.result) {
            stop();
            setBusy(false);
            setSolution(data.result);
            setStatus("Validated · solution verified");
          }
        };
        current.onerror = () =>
          fail(
            "The worker stopped unexpectedly. Reduce map complexity and try again.",
          );
        timeout.current = setTimeout(
          () =>
            fail(
              `Stopped after ${seconds}s to keep your browser responsive. Exhaustive search can be exponential; try fewer connected rooms or increase the time budget.`,
            ),
          seconds * 1000,
        );
        current.postMessage({
          text,
          base: new URL(import.meta.env.BASE_URL, window.location.origin).href,
        });
      } catch {
        stop();
        setBusy(false);
        setError(
          "This browser could not start a Web Worker. Use a current browser with WebAssembly support.",
        );
      }
    },
    [stop],
  );
  return { solution, status, error, busy, solve, clear, cancel };
}
