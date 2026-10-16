# Agent Instructions: terr4m/terraform-provider-shell

## Repository & Technical Context

- **Project:** Terraform / OpenTofu provider for executing shell commands and scripts with lifecycle integration (`create`, `read`, `update`, `delete`), environment variable injection, transport files, and JSON stdout parsing.
- **Language / Toolchain:** Go 1.27+ (standard toolchain).
- **Core Framework:** HashiCorp Terraform Plugin Framework (`github.com/hashicorp/terraform-plugin-framework`).
- **Testing:** Standard Go `testing` package and `github.com/hashicorp/terraform-plugin-testing` for acceptance tests.
- **Task Runner:** `just` (see repository `justfile`).

---

## Go Style & Architectural Standards

### 1. Canonical References & Idiomatic Go

Adhere to upstream Go standards and the repository-specific guidance below:

- [Effective Go](https://go.dev/doc/effective_go) for syntax conventions, naming, and control flow.
- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments) & [Google's Go Style Guide](https://google.github.io/styleguide/go/) for receiver naming, error string formatting, and package comments.

Also consider the [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md) for concurrency management, zero-value mutex safety, interface pollution avoidance, and container boundaries if it doesn't conflict with the other style guidelines.

Key implementation rules:

- **Left-Aligned Happy Path:** Guard early with error checks and return immediately; avoid deep nested conditional blocks.
- **Short Receivers:** Use 1–2 letters derived from the struct name (e.g., `(r *shellResource)`), never `this` or `self`.
- **Diagnostics & Errors:** In provider code, always append errors and warnings to `diag.Diagnostics` (`resp.Diagnostics.AddError(...)`). In internal library code, wrap errors with context using `fmt.Errorf("...: %w", err)`. Never swallow errors or invoke `panic`.
- **Interface Pollution:** Prefer accepting interfaces and returning concrete structs. Define interfaces at the point of use when practical; introduce a shared or earlier interface only when it serves a concrete design or testing need.

### 2. Plugin Framework Model & Schema Discipline

- **Test Comparisons:** Use `github.com/google/go-cmp/cmp` when it makes test comparisons clearer, especially for structured values; prefer `cmp.Diff` to opaque deep-equality failures. It is already a project dependency, so do not add another comparison library for this purpose.
- **Model Partitioning:** Keep distinct structs for configuration (`Config`), planned state (`Plan`), and saved state (`State`).
- **Framework Types:** Use framework attribute types (`types.String`, `types.Bool`, `types.Map`, `types.Dynamic`, etc.) inside models. Never map raw Go primitives directly to schema attributes.
- **Null & Unknown Checking:** Always check `.IsNull()` and `.IsUnknown()` before accessing primitive values via `.ValueString()`, `.ValueBool()`, etc.

### 3. Subprocess & OS Execution Mechanics

- **Context Passing:** Every subprocess execution must accept `ctx context.Context` via `os/exec.CommandContext` and honor cancellations and timeouts.
- **Process Group Termination:** When a command can spawn child processes that must not outlive it, ensure cancellation terminates the relevant process tree/group. Account for platform differences (such as POSIX process groups and Windows Job Objects) where this guarantee is required.
- **Buffer & Pipe Safety:** For potentially large or uncontrolled command output, avoid unbounded memory accumulation; stream or impose an appropriate size limit. Use `io.ReadAll` where input size is known or bounded and the data is suitably small.

---

## Tooling & Command Workflows

Always prefer executing tasks via `just` to maintain parity with local development and CI pipelines.

### 1. Fast Development Inner Loop

Run the relevant checks as you work. Format and test changes routinely; run dependency cleanup when dependencies or module files change:

```bash
# Format Go code with the repository formatter
just fmt

# Run when dependencies or go.mod/go.sum change
just tidy

# Fast unit tests (no race detector overhead)
just test
```

### 2. CI Verification Loop

Run these checks before committing, opening pull requests, or declaring changes complete:

```bash
# Static analysis and lint checks (golangci-lint run)
just lint

# Concurrency verification (full test suite with race detector)
go test -race ./...

# Acceptance tests (requires local shell environment and Terraform/OpenTofu binary)
just testacc
```

### 3. Documentation Generation (Decoupled Workflow)

Documentation generation is an isolated process. Only run documentation commands when modifying schema attributes, descriptions, defaults, or repository examples:

```bash
# Generate provider docs from schema and templates
just docs

# Format documentation files (run after docs generation)
just docs-fmt

# Verify docs are cleanly committed and pass validation
just docs-lint
```

---

## Agent Operational Guardrails

- **Docs Lifecycle Separation:** Do not touch `docs/` or run documentation generation targets unless schema definitions, attribute descriptions, or examples were directly modified.
- **Zero Breaking State Changes:** Modifying existing resource or data-source schemas must remain backwards compatible with existing Terraform states unless instructed to implement a schema migration (`StateUpgraders`).
- **Dependency Discipline:** Do not add external modules to `go.mod` without explicit user permission. Prefer the Go standard library and existing dependencies; if a task appears to require a new module, explain the need and ask before adding it.
- **Working Tree Cleanliness:** Produce focused, atomic diffs. Never run destructive git commands (`git reset --hard`, `git clean -fd`, `git stash drop`).
