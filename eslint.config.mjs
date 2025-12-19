import { defineConfig, globalIgnores } from "eslint/config";
import nextVitals from "eslint-config-next/core-web-vitals";
import nextTs from "eslint-config-next/typescript";

const eslintConfig = defineConfig([
  ...nextVitals,
  ...nextTs,
  // Override default ignores of eslint-config-next.
  globalIgnores([
    // Default ignores of eslint-config-next:
    ".next/**",
    "out/**",
    "build/**",
    "next-env.d.ts",
    "node_modules/**",
    "!.prettierrc.cjs",
    "!.eslintrc.cjs",
    // Custom ignores:
    "public/**",
    "scripts/**",
    "commitlint.config.mjs",
    "eslint.config.mjs",
    "tailwind.config.cjs",
    "postcss.config.cjs",
  ]),
]);

export default eslintConfig;
