import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  base: "/3d-lim-in-solver/",
  plugins: [react()],
  worker: { format: "es" },
  build: { chunkSizeWarningLimit: 1200 },
});
