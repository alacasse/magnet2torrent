# Repository Guidelines

## Project Structure & Module Organization
- Keep runnable code in `src/` (library and CLI entry points); group helpers by feature (e.g., `src/parser/`, `src/http/`).
- Place automated tests in `tests/` mirroring source paths (`tests/parser/…`); add fixtures under `tests/fixtures/`.
- Use `scripts/` for repeatable dev automation (formatting, lint, release prep) and `docs/` for user-facing notes or ADRs.

## Build, Test, and Development Commands
- Install dependencies before anything else: `npm install` (or `pnpm install` if preferred and documented).
- Local dev loop: `npm run dev` for a watcher-friendly build/run; ensure it exercises the main CLI entry.
- Quality gates: `npm run lint` for static checks and `npm test` for the full test suite; add coverage flags (`npm test -- --coverage`) when relevant.
- Production build: `npm run build` should output distributable artifacts (e.g., `dist/` bundle or packaged CLI); document any extra environment variables in `docs/`.

## Coding Style & Naming Conventions
- Default to TypeScript where possible; use 2-space indentation, single quotes, and trailing commas in multi-line structures.
- Run Prettier and ESLint before pushing; keep configs in repo (`.prettierrc`, `.eslintrc.*`) and align editor settings to them.
- Name files by responsibility (`magnet-parser.ts`, `torrent-writer.ts`); export narrow public APIs from index files to keep imports clean.

## Testing Guidelines
- Write unit tests alongside each module; prefer Jest or Vitest with file names `*.test.ts`.
- Cover parsing, error handling, and integration paths for magnet input → torrent output; include regression tests for previously reported bugs.
- Keep tests deterministic: stub network and filesystem IO, use fixtures, and avoid timing-based assertions.

## Commit & Pull Request Guidelines
- Use Conventional Commits (`feat:`, `fix:`, `chore:`, `test:`) with imperative, present-tense subjects.
- For PRs: include a concise summary, linked issues, and before/after notes (logs or CLI output). Add screenshots only if UI behavior is affected.
- Ensure CI passes locally (`npm run lint && npm test && npm run build`) before requesting review; mention any skipped checks and why.
