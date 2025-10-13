<!-- OPENSPEC:START -->
# OpenSpec Instructions

These instructions are for AI assistants working in this project.

Always open `@/openspec/AGENTS.md` when the request:
- Mentions planning or proposals (words like proposal, spec, change, plan)
- Introduces new capabilities, breaking changes, architecture shifts, or big performance/security work
- Sounds ambiguous and you need the authoritative spec before coding

Use `@/openspec/AGENTS.md` to learn:
- How to create and apply change proposals
- Spec format and conventions
- Project structure and guidelines

Keep this managed block so 'openspec update' can refresh the instructions.

<!-- OPENSPEC:END -->

Personal website for Conner Ohnesorge written with Go, Tailwind, and go-templ.
.
├── AGENTS.md
├── bun.lock
├── bun.nix
├── CLAUDE.md
├── cmd
│   ├── conneroh
│   │   ├── classes
│   │   ├── components
│   │   ├── doc.go
│   │   ├── errors.go
│   │   ├── handlers.go
│   │   ├── layouts
│   │   ├── README.md
│   │   ├── root.go
│   │   ├── routes.go
│   │   ├── _static
│   │   └── views
│   ├── update
│   │   └── main.go
│   └── update-css
│       └── main.go
├── flake.lock
├── flake.nix
├── GEMINI.md
├── go.mod
├── go.sum
├── index.ts
├── input.css
├── internal
│   ├── assets
│   │   ├── doc.go
│   │   ├── emp.go
│   │   ├── errors.go
│   │   ├── hash.go
│   │   ├── markdown.go
│   │   ├── paths.go
│   │   ├── README.md
│   │   ├── s3.go
│   │   ├── static.go
│   │   ├── types.go
│   │   ├── upsert.go
│   │   └── validation.go
│   ├── cache
│   │   ├── docs.hash
│   │   └── templ.hash
│   ├── copygen
│   │   ├── copygen.go
│   │   ├── setup.go
│   │   └── setup.yml
│   ├── data
│   │   ├── assets
│   │   ├── employments
│   │   ├── posts
│   │   ├── projects
│   │   └── tags
│   ├── logger
│   │   ├── dev.go
│   │   ├── doc.go
│   │   └── prod.go
│   └── routing
│       ├── doc.go
│       ├── handlers.go
│       ├── pagination.go
│       ├── pagination_test.go
│       ├── README.md
│       ├── slugs.go
│       └── targets.go
├── main.go
├── master.db
├── openspec
│   ├── AGENTS.md
│   ├── changes
│   ├── project.md
│   └── specs
├── package.json
├── README.md
├── tailwind.config.js
└── tsconfig.json
