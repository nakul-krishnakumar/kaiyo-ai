import js from "@eslint/js";
import globals from "globals";
import react from "eslint-plugin-react";
import tseslint from "typescript-eslint";
import reactHooks from "eslint-plugin-react-hooks";
import jsxA11y from "eslint-plugin-jsx-a11y";

export default tseslint.config(
  js.configs.recommended,
  ...tseslint.configs.recommended,
  {
    files: ["**/*.{ts,tsx}"],
    ignores: ["dist", "build", "node_modules"],
    languageOptions: {
      globals: {
        ...globals.browser,
        ...globals.node,
      },
      parserOptions: {
        ecmaVersion: "latest",
        sourceType: "module",
        ecmaFeatures: {
          jsx: true,
        },
      },
    },
    plugins: {
      react,
      "react-hooks": reactHooks,
      "jsx-a11y": jsxA11y,
    },
    rules: {
      // ✅ Disable legacy React import requirement
      "react/react-in-jsx-scope": "off",

      // ✅ Hooks linting
      "react-hooks/rules-of-hooks": "error",
      "react-hooks/exhaustive-deps": "warn",

      // ✅ Accessibility
      "jsx-a11y/alt-text": "warn",

      // ✅ TypeScript relaxed rules for dev
      "@typescript-eslint/no-unused-vars": "warn",
      "@typescript-eslint/no-explicit-any": "off",

      // ✅ Allow .tsx extensions
      "react/jsx-filename-extension": ["warn", { extensions: [".tsx"] }],
    },
    settings: {
      react: {
        version: "detect",
      },
    },
  }
);
