import { test, expect } from "@playwright/test";
import type { Page } from "@playwright/test";
import { readFileSync } from "node:fs";
import { resolve } from "node:path";
const intro = readFileSync(resolve("../test/example00.txt"), "utf8");
async function scrub(page: Page, value: number) {
  await page.getByLabel("Simulation timeline").evaluate((element, value) => {
    const setter = Object.getOwnPropertyDescriptor(
      HTMLInputElement.prototype,
      "value",
    )!.set!;
    setter.call(element, String(value));
    element.dispatchEvent(new Event("input", { bubbles: true }));
    element.dispatchEvent(new Event("change", { bubbles: true }));
  }, value);
}

test("production load, real WASM, exact turns, playback, camera and completion", async ({
  page,
}, testInfo) => {
  const failures: string[] = [];
  page.on("pageerror", (e) => failures.push(e.message));
  page.on("requestfailed", (r) =>
    failures.push(`${r.url()}: ${r.failure()?.errorText}`),
  );
  page.on("response", (r) => {
    if (r.status() >= 400) failures.push(`${r.status()} ${r.url()}`);
  });
  await page.goto("./");
  await expect(page.getByText("Validated · solution verified")).toBeVisible();
  await expect(page.getByTestId("metric-total-ants")).toHaveText("22");
  await expect(page.getByTestId("metric-rooms")).toHaveText("27");
  await expect(page.getByTestId("metric-paths")).toHaveText("4");
  await expect(page.getByTestId("total-turns")).toHaveText("11");
  await expect(page.locator("canvas, .fallback svg")).toBeVisible();
  if (await page.locator("canvas").count()) {
    await expect(
      page.getByRole("button", { name: "Inspect room start", exact: true }),
    ).toHaveAttribute("style", /translate/);
  }
  if (testInfo.project.name === "chromium")
    await expect(page.locator("canvas")).toBeVisible();
  const overflow = await page.evaluate(
    () => document.documentElement.scrollWidth > window.innerWidth,
  );
  expect(overflow).toBe(false);
  await page.screenshot({
    path: testInfo.outputPath("initial.png"),
    fullPage: true,
    animations: "disabled",
  });
  if (await page.locator("canvas").count()) {
    await page
      .getByRole("button", { name: "Inspect room start", exact: true })
      .click();
    await expect(page.locator(".room-popover")).toContainText(
      "Coordinates 0, 3",
    );
    await page.getByLabel("Close room information", { exact: true }).click();
  }
  await page.getByLabel("Next turn", { exact: true }).click();
  await expect(page.getByTestId("current-turn")).toHaveText("01");
  await page.getByLabel("Previous turn", { exact: true }).click();
  await expect(page.getByTestId("current-turn")).toHaveText("00");
  await page.getByLabel("Playback speed").selectOption("4");
  await page.getByLabel("Play simulation", { exact: true }).click();
  await expect(page.getByTestId("current-turn")).not.toHaveText("00");
  await page.getByLabel("Pause simulation", { exact: true }).click();
  const paused = await page.getByLabel("Simulation timeline").inputValue();
  await page.waitForTimeout(300);
  await expect(page.getByLabel("Simulation timeline")).toHaveValue(paused);
  await page.getByLabel("Restart simulation", { exact: true }).click();
  await expect(page.getByTestId("current-turn")).toHaveText("00");
  await page.locator("#example").selectOption("intro");
  await expect(page.getByTestId("metric-total-ants")).toHaveText("—");
  await page.getByRole("button", { name: "Run solver" }).click();
  await expect(page.getByText("Validated · solution verified")).toBeVisible();
  await expect(page.getByTestId("metric-total-ants")).toHaveText("4");
  await expect(page.getByTestId("total-turns")).toHaveText("6");
  await page.getByLabel("Next turn", { exact: true }).click();
  await expect(page.getByTestId("turn-output")).toHaveText("L1-2 ");
  await page.getByRole("tab", { name: "Ant inspector" }).click();
  await expect(
    page.getByRole("tabpanel", { name: "Ant inspector" }),
  ).toContainText("Route 1");
  await expect(
    page.getByRole("tabpanel", { name: "Ant inspector" }),
  ).toContainText("2");
  await page.getByRole("tab", { name: "Solver output" }).click();
  await scrub(page, 6);
  await expect(page.getByTestId("arrived")).toHaveText("4 / 4");
  await expect(page.getByTestId("turn-output")).toHaveText("L4-1 ");
  await expect(page.getByText("Every ant has arrived.")).toBeVisible();
  await scrub(page, 2);
  await expect(page.getByTestId("arrived")).toHaveText("0 / 4");
  await expect(page.getByTestId("turn-output")).toHaveText("L1-3 L2-2 ");
  await page.getByLabel("Reset camera", { exact: true }).click();
  await page.getByLabel("Fit map to view", { exact: true }).click();
  await page.getByLabel("Toggle coordinate or radial layout").click();
  await expect(page.getByText("RADIAL LAYOUT", { exact: true })).toBeVisible();
  await page.getByLabel("Hide room labels", { exact: true }).click();
  await page.getByLabel("Show room labels", { exact: true }).click();
  await page.getByRole("button", { name: /ROUTE 01/ }).click();
  await expect(page.getByRole("button", { name: /ROUTE 01/ })).toHaveAttribute(
    "aria-pressed",
    "true",
  );
  await page.getByText("How Lem-in Works", { exact: true }).click();
  await expect(
    page.getByText("Exhaustive search. Measured correctness."),
  ).toBeVisible();
  expect(failures).toEqual([]);
});

test("uploads, edits, invalid maps, downloads and solver recovery", async ({
  page,
}) => {
  await page.goto("./");
  await expect(page.getByText("Validated · solution verified")).toBeVisible();
  await page.getByLabel("Upload map file").setInputFiles({
    name: "custom.txt",
    mimeType: "text/plain",
    buffer: Buffer.from(intro),
  });
  await expect(page.getByLabel("Map definition", { exact: true })).toHaveValue(
    intro.replaceAll("\r", ""),
  );
  await page.getByRole("button", { name: "Run solver" }).click();
  await expect(page.getByTestId("total-turns")).toHaveText("6");
  await page.getByLabel("Ant count", { exact: true }).fill("7");
  await expect(page.getByTestId("metric-total-ants")).toHaveText("—");
  await page.getByRole("button", { name: "Run solver" }).click();
  await expect(page.getByTestId("metric-total-ants")).toHaveText("7");
  await expect(page.getByTestId("total-turns")).toHaveText("9");
  const downloading = page.waitForEvent("download");
  await page.getByLabel("Download map", { exact: true }).click();
  expect((await downloading).suggestedFilename()).toBe("colony.txt");
  await page.getByLabel("Map definition", { exact: true }).fill("0\n");
  await page.getByRole("button", { name: "Run solver" }).click();
  await expect(page.getByRole("alert")).toContainText("line 1: Invalid ants");
  await expect(page.locator("canvas")).toHaveCount(0);
  await page
    .getByLabel("Map definition", { exact: true })
    .fill("3\n##start\ns 0 0\n##end\ne 1 0\n");
  await page.getByRole("button", { name: "Run solver" }).click();
  await expect(page.getByRole("alert")).toContainText("no path connects");
  await page.getByLabel("Reset to example", { exact: true }).click();
  await page.getByRole("button", { name: "Run solver" }).click();
  await expect(page.getByText("Validated · solution verified")).toBeVisible();
  await page.getByLabel("Clear editor", { exact: true }).click();
  await expect(page.getByLabel("Map definition", { exact: true })).toHaveValue(
    "",
  );
});

test("timeout, cancellation and worker recovery keep the UI responsive", async ({
  page,
}) => {
  await page.goto("./");
  await expect(page.getByText("Validated · solution verified")).toBeVisible();
  // Delay only the actual WASM request to deterministically test lifecycle timeouts.
  await page.route("**/solver.wasm", async (route) => {
    await new Promise((r) => setTimeout(r, 6500));
    await route.continue().catch(() => {});
  });
  await page.getByLabel("Solver time budget").selectOption("5");
  await page.getByRole("button", { name: "Run solver" }).click();
  await page.getByRole("button", { name: "Cancel computation" }).click();
  await expect(page.getByText("Cancelled · ready to run again")).toBeVisible();
  await page.getByRole("button", { name: "Run solver" }).click();
  await expect(page.getByRole("alert")).toContainText("Stopped after 5s");
  await page.unroute("**/solver.wasm");
  await page.getByRole("button", { name: "Run solver" }).click();
  await expect(page.getByText("Validated · solution verified")).toBeVisible();
});

test("reduced motion and accessible graph fallback retain working playback", async ({
  page,
}) => {
  await page.emulateMedia({ reducedMotion: "reduce" });
  await page.addInitScript(() => {
    const original = HTMLCanvasElement.prototype.getContext;
    HTMLCanvasElement.prototype.getContext = function (
      type: string,
      ...args: unknown[]
    ) {
      if (
        type === "webgl" ||
        type === "webgl2" ||
        type === "experimental-webgl"
      )
        return null;
      return original.call(this, type, ...(args as [])) as never;
    } as typeof original;
  });
  await page.goto("./");
  await expect(page.getByText("Validated · solution verified")).toBeVisible();
  await expect(page.locator(".scene > .fallback svg")).toBeVisible();
  await page.getByLabel("Next turn", { exact: true }).click();
  await expect(page.getByTestId("current-turn")).toHaveText("01");
  await scrub(page, 11);
  await expect(page.getByTestId("arrived")).toHaveText("22 / 22");
});
