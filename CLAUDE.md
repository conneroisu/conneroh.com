# YOU ARE THE ORCHESTRATOR

You are Claude Code with a 200k context window, and you ARE the orchestration system. You manage the entire project, create todo lists, and delegate individual tasks to specialized subagents.

## 🎯 Your Role: Master Orchestrator

You maintain the big picture, create comprehensive todo lists, and delegate individual todo items to specialized subagents that work in their own context windows.

## 🚨 YOUR MANDATORY WORKFLOW

When the user gives you a project:

### Step 1: ANALYZE & PLAN (You do this)
1. Understand the complete project scope
2. Break it down into clear, actionable todo items
3. **USE TodoWrite** to create a detailed todo list
4. Each todo should be specific enough to delegate

### Step 2: DELEGATE TO SUBAGENTS (One todo at a time)
1. Take the FIRST todo item
2. Invoke the **`coder`** subagent with that specific task (Never trust that the `coder` agent will complete the task correctly always verify, test, and investigate changes)
3. The coder works in its OWN context window
4. Wait for coder to complete and report back

### Step 3: TEST THE IMPLEMENTATION
1. Take the coder's completion report
2. Invoke the **`tester`** subagent to verify
3. Tester uses Playwright MCP in its OWN context window
4. Wait for test results

### Step 4: HANDLE RESULTS
- **If tests pass**: Mark todo complete, move to next todo
- **If tests fail**: Invoke **`stuck`** agent for human input
- **If coder hits error**: They will invoke stuck agent automatically

### Step 5: ITERATE
1. Update todo list (mark completed items)
2. Move to next todo item
3. Repeat steps 2-4 until ALL todos are complete

## 🛠️ Available Subagents

### coder
**Purpose**: Implement one specific todo item

- **When to invoke**: For each coding task on your todo list
- **What to pass**: ONE specific todo item with clear requirements
- **Context**: Gets its own clean context window
- **Returns**: Implementation details and completion status
- **On error**: Will invoke stuck agent automatically

### tester
**Purpose**: Visual verification with Playwright MCP

- **When to invoke**: After EVERY coder completion
- **What to pass**: What was just implemented and what to verify
- **Context**: Gets its own clean context window
- **Returns**: Pass/fail with screenshots
- **On failure**: Will invoke stuck agent automatically

### stuck
**Purpose**: Human escalation for ANY problem

- **When to invoke**: When tests fail or you need human decision
- **What to pass**: The problem and context
- **Returns**: Human's decision on how to proceed
- **Critical**: ONLY agent that can use AskUserQuestion

## 🚨 CRITICAL RULES FOR YOU

**YOU (the orchestrator) MUST:**
1. ✅ Create detailed todo lists with TodoWrite
2. ✅ Delegate ONE todo at a time to coder
3. ✅ Test EVERY implementation with tester
4. ✅ Track progress and update todos
5. ✅ Maintain the big picture across 200k context
6. ✅ **ALWAYS create pages for EVERY link in headers/footers** - NO 404s allowed!

**YOU MUST NEVER:**
1. ❌ Implement code yourself (delegate to coder)
2. ❌ Skip testing (always use tester after coder)
3. ❌ Let agents use fallbacks (enforce stuck agent)
4. ❌ Lose track of progress (maintain todo list)
5. ❌ **Put links in headers/footers without creating the actual pages** - this causes 404s!

## 📋 Example Workflow

```
User: "Build a React todo app"

YOU (Orchestrator):
1. Create todo list:
   [ ] Set up React project
   [ ] Create TodoList component
   [ ] Create TodoItem component
   [ ] Add state management
   [ ] Style the app
   [ ] Test all functionality

2. Invoke coder with: "Set up React project"
   → Coder works in own context, implements, reports back

3. Invoke tester with: "Verify React app runs at localhost:3000"
   → Tester uses Playwright, takes screenshots, reports success

4. Mark first todo complete

5. Invoke coder with: "Create TodoList component"
   → Coder implements in own context

6. Invoke tester with: "Verify TodoList renders correctly"
   → Tester validates with screenshots

... Continue until all todos done
```

## 🔄 The Orchestration Flow

```
USER gives project
    ↓
YOU analyze & create todo list (TodoWrite)
    ↓
YOU invoke coder(todo #1)
    ↓
    ├─→ Error? → Coder invokes stuck → Human decides → Continue
    ↓
CODER reports completion
    ↓
YOU invoke tester(verify todo #1)
    ↓
    ├─→ Fail? → Tester invokes stuck → Human decides → Continue
    ↓
TESTER reports success
    ↓
YOU mark todo #1 complete
    ↓
YOU invoke coder(todo #2)
    ↓
... Repeat until all todos done ...
    ↓
YOU report final results to USER
```

## 🎯 Why This Works

**Your 200k context** = Big picture, project state, todos, progress
**Coder's fresh context** = Clean slate for implementing one task
**Tester's fresh context** = Clean slate for verifying one task
**Stuck's context** = Problem + human decision

Each subagent gets a focused, isolated context for their specific job!

## 💡 Key Principles

1. **You maintain state**: Todo list, project vision, overall progress
2. **Subagents are stateless**: Each gets one task, completes it, returns
3. **One task at a time**: Don't delegate multiple tasks simultaneously
4. **Always test**: Every implementation gets verified by tester
5. **Human in the loop**: Stuck agent ensures no blind fallbacks

## 🚀 Your First Action

When you receive a project:

1. **IMMEDIATELY** use TodoWrite to create comprehensive todo list
2. **IMMEDIATELY** invoke coder with first todo item
3. Wait for results, test, iterate
4. Report to user ONLY when ALL todos complete

## ⚠️ Common Mistakes to Avoid

❌ Implementing code yourself instead of delegating to coder
❌ Skipping the tester after coder completes
❌ Delegating multiple todos at once (do ONE at a time)
❌ Not maintaining/updating the todo list
❌ Reporting back before all todos are complete
❌ **Creating header/footer links without creating the actual pages** (causes 404s)
❌ **Not verifying all links work with tester** (always test navigation!)

## ✅ Success Looks Like

- Detailed todo list created immediately
- Each todo delegated to coder → tested by tester → marked complete
- Human consulted via stuck agent when problems occur
- All todos completed before final report to user
- Zero fallbacks or workarounds used
- **ALL header/footer links have actual pages created** (zero 404 errors)
- **Tester verifies ALL navigation links work** with Playwright

---

**You are the conductor with perfect memory (200k context). The subagents are specialists you hire for individual tasks. Together you build amazing things!** 🚀
Note that if you are waiting for an action to complete you should not return, you must call a Bash(sleep {best estimate of seconds to sleep until complete}).
See documents for libaries and frameworks used in .claude/contexts/
<project>
Platform leveraging Nix/NixOS for orchestrating AI coding agents in reproducible, containerized development environments powered by Nix.

Projects:
* ./src/     - The Tanstack backend and frontend web console.
* ./docs/    - AstroJS MDX Starlight Documentation site.
* ./tests/   - Bun Tests & End-to-end tests for the web console and API powered by Playwright.
* ./sidecar/ - Cmake-based C++ build system for the /nix/ k8s sidecar.
* ./api/     - Backend Golang API that powers the web console.
* ./db/      - One and only collection of Postgres migrations and schema.
* ./proto/   - Protobuf definitions for the API (Generated with `nix run .#proto`)

./Procfile - Overmind configuration for running all the services locally.

# Constitution Enforcement

Key principles:

- **Test-First Development (TDD)**: Tests MUST be written before implementation (Red-Green-Refactor)
- **Zero Technical Debt**: No TODOs, no mock data, all code must pass strict linting
- **Performance Targets**: 95+ Lighthouse score, <100ms API responses
- **80% Test Coverage Minimum**: Critical paths must have comprehensive tests
- **Authentication and Authorization**: Use [better-auth](https://github.com/go-playground/better-auth) for authentication and authorization in the tanstack backend (Never in the go backend).
- **Personal Organization Plugin**: Custom better-auth plugin automatically creates personal organizations after email verification (see `./specs/005-personal-orgs-custom/` for implementation details).
- **Organization Projects Plugin**: Custom better-auth plugin that adds project management capabilities to organizations (see `./src/lib/plugins/projects/` for implementation).

## Validation Commands
```bash
# Full Linting
nix develop -c `lint`
# Full Testing
nix develop -c `tests`

# Individual checks
biome lint && cd .. # Console
golangci-lint run # API
```

## Quality Gates
- **Pre-push**: All tests pass, 80% coverage, type checking
- **Pull Request**: Integration tests, performance budgets, documentation
- **Merge**: E2E tests, no conflicts, changelog updated

## Better-Auth Plugins

### Organization Environment Plugin

**Location**: `./src/lib/plugins/organization-environment/`

The Organization Environment plugin adds Nix-powered development environment management capabilities to organizations. Environments are containerized development workspaces that AI agents and users interact with.

**Features**:
- **CRUD Operations**: Create, read, update, delete environments
- **Organization Integration**: Environments belong to organizations with member-based access control
- **Project Association**: Optional project association for environment-project relationships
- **Lifecycle Hooks**: Customizable before/after hooks for all operations
- **Permission System**: Organization membership for read, ownership for modify/delete
- **URI Management**: Flexible URI validation supporting custom Nix schemes (`nix+http://`)
- **Soft Deletion**: Support for soft delete with `deleted_at` timestamp
- **Build Protection**: Prevent deletion of environments with active builds
- **Paths Management**: Native Postgres TEXT[] for filesystem path arrays

**API Endpoints**:
- `POST /api/auth/environment/create` - Create new environment
- `GET /api/auth/environment/list` - List environments (filterable by org/project/user)
- `GET /api/auth/environment/get` - Get environment by ID
- `PATCH /api/auth/environment/update` - Update environment details
- `DELETE /api/auth/environment/delete` - Delete environment (checks for active builds)

**Configuration Options**:
```typescript
organizationEnvironment({
  // Maximum environments per organization (default: 50)
  maxEnvironmentsPerOrganization: 50,

  // Maximum environments per user across all orgs (default: 10)
  maxEnvironmentsPerUser: 10,

  // Enable URI validation (default: true)
  validateUri: true,

  // Custom lifecycle hooks
  environmentHooks: {
    beforeCreateEnvironment: async ({ environment, organization, user, project }) => {
      // Custom validation logic
    },
    afterCreateEnvironment: async ({ environment, organization, user, project }) => {
      // Post-creation actions (e.g., logging, notifications)
    },
    // ... other hooks
  },
})
```

**Usage with TanStack Query**:
```typescript
import { authClient } from "@/lib/auth-client";

// List environments in organization
const { data: environments } = useQuery({
  queryKey: ['environments', organizationId],
  queryFn: () => authClient.environment.list({ organizationId })
});

// Create new environment
const createMutation = useMutation({
  mutationFn: (data) => authClient.environment.create(data)
});

createMutation.mutate({
  organizationId: "org_123",
  name: "My Dev Environment",
  uri: "https://env.connix.io/abc",
  description: "Development environment for project X",
  paths: ["/src", "/tests"],
  projectId: "proj_456", // Optional
});
```

**Database Schema**:
- **Table**: `environment`
- **Columns**: id, name, description, paths (TEXT[]), uri, org_id, user_id, created_at, updated_at, deleted_at
- **Indexes**: org_id, user_id, deleted_at (partial)
- **Constraints**: Foreign keys to organization, user

**Related Files**:
- `src/lib/plugins/organization-environment/index.ts` - Server plugin
- `src/lib/plugins/organization-environment/client.ts` - Client plugin
- `src/lib/plugins/organization-environment/types.ts` - TypeScript types
- `src/lib/plugins/organization-environment/schema.ts` - Database schema
- `src/lib/plugins/organization-environment/routes/` - Endpoint implementations
- `db/migrations/00005_environment_plugin_updates.sql` - Plugin migration
- `db/migrations/

### Organization Deployment Plugin

**Location**: `./src/lib/plugins/deployment/`

The Organization Deployment plugin adds deployment management capabilities to environments. Deployments are running instances of builds deployed to Kubernetes pods, representing the live execution environment for Nix-built containers.

**Features**:
- **Deployment Lifecycle Management**: Track deployments through pending, deploying, active, stopping, stopped, and failed states
- **Kubernetes Integration**: Support for K8s namespaces, pod names, and deployment URLs
- **Build Association**: Link deployments to specific builds with cascade protection
- **Status Tracking**: Comprehensive state management with timestamps for created, deployed, and terminated times
- **Active Deployment Limits**: Configurable maximum active deployments per organization
- **Metadata Support**: Extensible JSONB metadata field for custom deployment data
- **Lifecycle Hooks**: Customizable before/after hooks for all deployment operations
- **Build Deletion Protection**: Prevent deletion of builds with active deployments (via ON DELETE RESTRICT)
- **Pod Name Uniqueness**: Ensure no namespace/pod name collisions for active deployments

**Error Handling**:
- Invalid resource quantity strings (CPU/memory requests or limits) fail fast during pod spec construction. The deployment status transitions to `failed` and `error_message` captures the offending field (e.g., `invalid cpu request: ...`).
- When resource metadata is omitted, safe defaults are applied (100m/1000m CPU, 256Mi/1Gi memory). Provide Kubernetes-formatted values when overriding defaults:

```json
{
  "resources": {
    "requests": {
      "cpu": "500m",
      "memory": "512Mi"
    },
    "limits": {
      "cpu": "2000m",
      "memory": "2Gi"
    }
  }
}
```

**API Endpoints**:
- `POST /api/auth/deployment/create` - Create new deployment
- `GET /api/auth/deployment/list` - List deployments (filterable by org/environment/status)
- `GET /api/auth/deployment/get` - Get deployment by ID
- `PATCH /api/auth/deployment/update` - Update deployment status/metadata (admin only)
- `POST /api/auth/deployment/terminate` - Terminate running deployment

**Configuration Options**:
```typescript
organizationDeployment({
  // Maximum active deployments per organization (default: 50)
  // Active = pending, deploying, or active status
  maxActiveDeploymentsPerOrg: 50,

  // Custom lifecycle hooks
  deploymentHooks: {
    beforeCreateDeployment: async ({ deployment, environment, build, organization, user }) => {
      // Custom validation logic
    },
    afterCreateDeployment: async ({ deployment, environment, build, organization, user }) => {
      // Post-creation actions (e.g., trigger K8s controller)
    },
    beforeUpdateDeployment: async ({ deployment, updates, user }) => {
      // Pre-update validation
    },
    afterUpdateDeployment: async ({ deployment, previousData, user }) => {
      // Post-update actions
    },
    beforeTerminateDeployment: async ({ deployment, environment, build, organization, user }) => {
      // Pre-termination checks
    },
    afterTerminateDeployment: async ({ deployment, environment, build, organization, user }) => {
      // Post-termination cleanup
    },
  },
})
```

**Usage with TanStack Query**:
```typescript
import {
  useDeployments,
  useDeployment,
  useCreateDeployment,
  useUpdateDeployment,
  useTerminateDeployment
} from "@/hooks/deployment-queries";

// List deployments in organization
const { data: deployments } = useDeployments("org_123");

// Filter by environment
const { data: envDeployments } = useDeployments("org_123", {
  envId: "env_456",
});

// Filter by status
const { data: activeDeployments } = useDeployments("org_123", {
  status: "active",
});

// Get single deployment
const { data: deployment } = useDeployment("deploy_123");

// Create new deployment
const createMutation = useCreateDeployment();
await createMutation.mutateAsync({
  envId: "env_123",
  buildId: "build_456",
  namespace: "connix-envs", // Optional, defaults to "connix-envs"
  metadata: { version: "1.0.0" }, // Optional
});

// Update deployment status (admin only)
const updateMutation = useUpdateDeployment();
await updateMutation.mutateAsync({
  deploymentId: "deploy_123",
  status: "active",
  podName: "pod-xyz-789",
  deploymentUrl: "https://deploy.connix.io/abc123",
});

// Terminate deployment
const terminateMutation = useTerminateDeployment();
await terminateMutation.mutateAsync({
  deploymentId: "deploy_123",
});
```

**Deployment Status Lifecycle**:
- **pending**: Deployment created, awaiting Kubernetes scheduling
- **deploying**: Deployment in progress, pod being created
- **active**: Deployment running and accessible (requires deployedAt timestamp)
- **stopping**: Deployment termination initiated
- **stopped**: Deployment successfully terminated
- **failed**: Deployment encountered an error

**Database Schema**:
- **Table**: `environment_deployment`
- **Columns**: id, envId, buildId, status (enum), podName, namespace, deploymentUrl, createdAt, deployedAt, terminatedAt, metadata (JSONB), errorMessage
- **Indexes**: envId (partial: active only), buildId, status (partial: active only), composite (envId, status)
- **Constraints**:
  - Foreign keys to environment (ON DELETE CASCADE), build (ON DELETE RESTRICT)
  - CHECK constraint: deployedAt must be NOT NULL when status is 'active'
  - UNIQUE constraint: namespace + podName (partial: active deployments only)

**Relationship to Environment and Build**:
- **Environment → Build → Deployment**: Deployments are instances of builds within environments
- **Cascade behavior**: Deleting an environment cascades to deployments; deleting a build is blocked if active deployments exist
- **Build protection**: Cannot delete builds with active deployments (prevents orphaned running pods)

**Related Files**:
- `src/lib/plugins/deployment/index.ts` - Server plugin
- `src/lib/plugins/deployment/client.ts` - Client plugin
- `src/lib/plugins/deployment/types.ts` - TypeScript types
- `src/lib/plugins/deployment/schema.ts` - Database schema and validation
- `src/lib/plugins/deployment/routes/` - Endpoint implementations
- `src/hooks/deployment-queries.ts` - React Query hooks
- `db/migrations/00007_deployments.sql` - Database migration

```txt
.
├── AGENTS.md
├── api
│   ├── build.go
│   ├── cache.go
│   ├── cmd
│   │   ├── doc.go
│   │   ├── root.go
│   │   ├── routes.go
│   │   ├── runner.go
│   │   ├── server.go
│   │   ├── shutdown.go
│   │   └── startup.go
│   ├── doc.go
│   ├── env.go
│   ├── init.go
│   └── internal
│       ├── api
│       ├── flakes
│       ├── nix
│       ├── operator
│       └── workers
├── internal
│       ├── kindest
│       ├── auth
│       └── servers
├── biome.jsonc
├── bunfig.toml
├── Cargo.toml
├── components.json
├── db
│   ├── account.sql.go
│   ├── apikey.sql.go
│   ├── auth.sql.go
│   ├── builds.sql.go
│   ├── chick (Clickhouse schema)
│   │   └── *.sql (Clickhouse schema migrations)
│   ├── db.go
│   ├── environments.sql.go
│   ├── invitation.sql.go
│   ├── member.sql.go
│   ├── migrations
│   │   └── *.sql (Migrations)
│   ├── models.go (Generated Sqlc Models)
│   ├── *.sql.go  (Generated Sqlc Queries)
│   ├── querier.go
│   ├── queries
│   │   ├── account.sql
│   │   ├── apikey.sql
│   │   ├── auth.sql
│   │   ├── builds.sql
│   │   ├── environments.sql
│   │   ├── invitation.sql
│   │   ├── member.sql
│   │   ├── oauthAccessToken.sql
│   │   ├── oauthApplication.sql
│   │   ├── oauthConsent.sql
│   │   ├── organization.sql
│   │   ├── session.sql
│   │   ├── subscription.sql
│   │   ├── twoFactor.sql
│   │   ├── user.sql
│   │   └── verification.sql
│   ├── session.sql.go
│   ├── subscription.sql.go
│   ├── twoFactor.sql.go
│   ├── user.sql.go
│   └── verification.sql.go
├── docker-compose.yml
├── Dockerfile
├── docs
│   ├── astro.config.mjs
│   ├── Dockerfile
│   ├── flake.nix
│   ├── fly.toml
│   ├── openapi-example-plugin.md
│   ├── package.json
│   ├── public
│   │   └── favicon.svg
│   ├── README.md
│   ├── src
│   │   ├── assets
│   │   ├── components
│   │   ├── content
│   │   └── content.config.ts
│   └── tsconfig.json
├── flake.nix
├── go.mod
├── go.sum
├── happydom.ts
├── infra
│   └── distributed-builders
│       └── configuration.nix
├── package.json
├── playwright.config.ts
├── Procfile
├── public
│   ├── favicon.ico
│   ├── logo192.png
│   ├── logo512.png
│   ├── manifest.json
│   └── robots.txt
├── publish.config.js
├── README.md
├── shell.nix
├── sidecar
│   ├── README.md
│   ├── INSPO.md
│   ├── k8s
│   │   ├── configmap.yaml
│   │   └── deployment.yaml
│   ├── src
│   └── test
├── sqlc.yaml
├── src
│   ├── api
│   │   ├── api.ts
│   │   ├── build.ts
│   │   ├── cache.ts
│   │   ├── environment.ts
│   │   ├── middleware.ts
│   │   └── types.ts
│   ├── components
│   │   ├── avatar-upload.tsx
│   │   ├── BuildLogsVirtualList.tsx
│   │   ├── dev
│   │   ├── falling-flakes.tsx
│   │   ├── form-fields.tsx
│   │   ├── forms
│   │   ├── hooks
│   │   ├── landing
│   │   ├── layout
│   │   ├── OrganizationCombobox.tsx
│   │   ├── presigned-avatar.tsx
│   │   ├── settings
│   │   ├── sidebar-dropdown.tsx
│   │   ├── sidebar.tsx
│   │   ├── stripe
│   │   └── ui
│   ├── emails
│   │   ├── notion-magic-link.tsx
│   │   ├── plaid-verify-identity.tsx
│   │   ├── static
│   │   ├── stripe-welcome.tsx
│   │   ├── vercel-invite-user.tsx
│   │   └── verification-email.tsx
│   ├── env.ts
│   ├── hooks
│   │   ├── auth-queries.ts
│   │   ├── form-hook.test.ts
│   │   ├── form-hook.ts
│   │   ├── use-animation-state.ts
│   │   ├── use-organization-context.ts
│   │   ├── use-reduced-motion.ts
│   │   └── use-theme.ts
│   ├── integrations
│   │   └── tanstack-query
│   ├── lib
│   │   ├── auth-client.ts
│   │   ├── auth.ts
│   │   ├── avatar-validation-client.ts
│   │   ├── avatar-validation-server.ts
│   │   ├── build-logs-collection.ts
│   │   ├── crop-utils.ts
│   │   ├── db.ts (pg db)
│   │   ├── email.ts
│   │   ├── icons
│   │   ├── logger.ts
│   │   ├── phone-validation.ts
│   │   ├── s3.ts
│   │   ├── seo.ts
│   │   ├── settings-tabs.ts
│   │   ├── stripe.ts
│   │   ├── theme-context.tsx
│   │   ├── user.ts
│   │   └── utils.ts
│   ├── logo.svg
│   ├── router.tsx
│   ├── routes
│   │   ├── api/
│   │   ├── organizations/
│   │   ├── about.tsx
│   │   ├── __root.tsx
│   │   ├── careers.tsx
│   │   ├── contact.tsx
│   │   ├── cookies.tsx
│   │   ├── dashboard.tsx
│   │   ├── forgot-password.tsx
│   │   ├── index.tsx
│   │   ├── offline.tsx
│   │   ├── pricing.tsx
│   │   ├── privacy.tsx
│   │   ├── reset-password.tsx
│   │   ├── security.tsx
│   │   ├── settings.tsx
│   │   ├── sign-in.tsx
│   │   ├── sign-up.tsx
│   │   ├── terms.tsx
│   │   └── verify-2fa.tsx
│   ├── routeTree.gen.ts
│   ├── styles
│   │   ├── animations.css
│   │   └── prose.css
│   └── styles.css
├── tests
│   ├── bun
│   │   ├── authentication.spec.tsx
│   │   └── subscription-status.spec.tsx
│   ├── integration
│   ├── performance
│   └── playwright
│       ├── accessibility.e2e.ts
│       └── basic-navigation.e2e.ts
├── tsconfig.check.json
├── tsconfig.json
├── tygo.yaml
├── typedoc.config.js
└── vite.config.ts
```

</project>

<!-- OPENSPEC:START -->
# OpenSpec Instructions

These instructions are for AI assistants working in this project.

Always open `@/openspec/AGENTS.md` when the request:
- Mentions planning or proposals (words like proposal, spec, change, plan)
- Introduces new capabilities, breaking changes (openspecs/changes/), architecture shifts, or big performance/security work
- Sounds ambiguous and you need the authoritative spec before coding

Use `@/openspec/AGENTS.md` to learn:
- How to create and apply change proposals
- Spec format and conventions
- Project structure and guidelines

Keep this managed block so 'openspec update' can refresh the instructions.

<!-- OPENSPEC:END -->
