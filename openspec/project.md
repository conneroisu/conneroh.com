# Project Context

## Purpose
Conner Ohnesorge’s personal portfolio site that showcases projects, long-form posts, speaking history, and employment highlights. The Go application renders pages server-side, keeps navigation fast with HTMX partial swaps, and provides a searchable, content-rich hub for prospective collaborators or employers.

## Tech Stack
- Backend: Go 1.24, standard library HTTP server, Bun ORM on SQLite (`master.db`)
- Rendering: `templ` components in `cmd/conneroh`, HTMX for progressive enhancement, Alpine.js (anchor/intersect plugins) for client interactivity, MathJax for rich math content
- Content pipeline: Markdown stored under `internal/data/**`, processed via `cmd/update` + `copygen` into SQLite, media uploaded to S3 through AWS SDK
- Asset pipeline: Tailwind CSS (CLI) with `twerge` pruning, Bun bundler for TypeScript/ESM (`index.ts` → `_static/dist/`), templ code generation (`templ generate`)
- Tooling: Nix flake (`nix develop`) for reproducible dev shells, Doppler for secrets, Fly.io deployment packages, GitHub Actions CI/CD, treefmt-nix (gofumpt/golines/alejandra)

## Project Conventions

### Code Style
- Go: always run `templ generate` before `go test`; format with `gofmt`, `gofumpt`, and `golines` via `nix fmt`/treefmt; keep handlers small and return `error` for centralized logging in `internal/routing`
- templ: colocate page/layout/component templates beneath `cmd/conneroh/`; prefer strongly-typed props and helper components in `components/`
- TypeScript: ESM modules bundled by Bun; keep DOM work declarative (HTMX triggers, Alpine stores) and emit TypeScript types when extending `window`
- CSS: Tailwind utilities first; non-trivial styles stay in `input.css` so `twerge` can fold unused classes into `classes.go`

### Architecture Patterns
- `cmd/` holds binaries (`conneroh` web server, `update`, `update-css`, `update-js`); `internal/` provides reusable domains (`assets`, `routing`, `logger`)
- Runtime uses a read-mostly SQLite file; web handlers query Bun models, compose templ components, and HTMX-aware wrappers (`routing.MorphableHandler`) decide between partial and full responses
- Markdown ingestion (`cmd/update`) hashes directories, upserts structured content, and pushes static assets to S3; UI generation (`cmd/update-css`) renders representative components to keep Tailwind class extraction deterministic
- HTMX navigation swaps the `#bodiody*` targets, while layouts ensure a consistent shell (Page vs Morpher) and Alpine powers minor client behaviour

### Testing Strategy
- `nix run .#runTests` (or `nix develop -c tests`) regenerates templ outputs then runs `go test ./...`; ensure new packages include unit tests
- Prefer table-driven Go tests; for pagination or routing utilities see `internal/routing/pagination_test.go` as reference
- Integration/browser tests are planned via Playwright (deps already vendored) and should live under `tests/` once added, executed through Nix targets
- CI blocks deploys: PRs trigger `pr-preview` workflow (tests + preview), `main` pushes run `cd.yml` (tests followed by Fly deploy)

### Git Workflow
- Work in short-lived feature branches branched from `main`; open PRs against `main`
- Keep commits scoped and imperative (e.g., “Add hero scroll handling”); large features require an OpenSpec change proposal before coding
- After review, merging to `main` kicks off the production CI/CD workflow; PRs automatically get Fly.io preview apps and status comments

## Domain Context
- Core entities: Posts, Projects, Tags, Employments; relationships stored in SQLite allow cross-linking (e.g., tags landing pages show posts/projects/employments)
- Content authored in Markdown with YAML front matter for metadata (slugs, created/updated dates, associations, banner paths)
- Site prioritises accessibility and fast navigation: HTMX handles partial swaps, Tailwind classes pre-generated for consistent styling, Math-heavy posts rely on MathJax
- Contact form submissions are captured server-side via templ forms and `gorilla/schema` parsing

## Important Constraints
- `cmd/update` requires `BUCKET_NAME` and AWS credentials (via Doppler) to push assets to S3 before DB refresh; run it before committing new content
- Generated artifacts (`cmd/conneroh/classes/classes.go`, `_static/dist/*`, `master.db`) are source-controlled; regenerate via `generate-all` or component-specific scripts when templates/content change
- Application expects read-only access to `master.db` at runtime; do not mutate content in handlers
- Keep Nix flake lockfiles (`flake.lock`, `bun.lock`) in sync with tooling updates; run `bun2nix` after dependency changes (`postinstall` script)
- HTMX targets (`#bodiody*`) and morph layouts assume consistent component IDs—changing them requires auditing morph handlers

## External Dependencies
- Fly.io for app hosting and automated PR previews/production deploys
- AWS S3 (Tigris interface) for static asset storage
- Doppler for secrets management in local scripts and CI
- GitHub Actions for CI/CD orchestration
- MathJax CDN for rendering mathematical content in posts
