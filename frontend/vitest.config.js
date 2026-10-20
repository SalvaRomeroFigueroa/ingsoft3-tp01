import { defineConfig } from "vitest/config";

export default defineConfig({
  test: {
    include: ["test/**/*.test.js"],
    coverage: {
      provider: "v8",
      reporter: ["text", "html", "lcov", "json-summary"],
      include: ["lib/**"],   // QUÉ entra en la cuenta: solo la lógica. app.js es pegamento de DOM.
      thresholds: { lines: 90, branches: 90 },
    },
  },
});