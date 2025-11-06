# Copilot Instructions

These instructions define how GitHub Copilot should assist with this Go project. The goal is to ensure consistent, high-quality code generation aligned with Go idioms, the chosen architecture, and our team's best practices.

## 🧠 Context

- **Project Type**: Terraform Provider
- **Language**: Go
- **Framework / Libraries**: `github.com/hashicorp/terraform-plugin-framework`

This project is a Terraform provider that supports executing arbitrary shell scripts.

### 📁 File Structure

Use this structure as a guide when creating or updating files:

```text
internal/
  provider/
  shell/
  tfdynamic/
go.mod
go.sum
main.go
```

## ⚙️ General Guidelines

- Follow idiomatic Go conventions (<https://go.dev/doc/effective_go>).
- Use named functions over long anonymous ones.
- Organize logic into small, composable functions.
- Prefer interfaces for dependencies to enable mocking and testing.
- Use `golangci-lint` to enforce formatting.
- Avoid unnecessary abstraction; keep things simple and readable.
- Use `context.Context` for request-scoped values and cancellation.

## 🛠️ Tools

- For linting, use `golangci-lint run`
- For formatting use `golangci-lint fmt`

## 🧶 Patterns

### ✅ Patterns to Follow

- Use **Clean Architecture** and **Repository Pattern**.
- Implement input validation using Go structs and validation tags.
- Use custom error types for wrapping and handling business logic errors.
- Logging should be handled via `log/slog`.
- Use dependency injection via constructors (avoid global state).
- Keep `main.go` minimal with `main` calling `run` before delegating to `internal`.

### 🚫 Patterns to Avoid

- Don't use global state unless absolutely required.
- Don't hardcode config—use environment variables or config files.
- Don't panic or exit in library code; return errors instead.
- Avoid embedding business logic in HTTP handlers.

## 🧪 Testing Guidelines

- Use `testing`, [matryer/is](https://github.com/matryer/is) & [google/go-cmp](https://github.com/google/go-cmp) for assertions.
- Add tests to `_test.go` files.
- Use separate namespace for black box tests.
- Mock external services using interfaces and mocks for unit tests.
- Include table-driven tests for functions with many input variants.
- Follow TDD for core business logic.

## 🔁 Iteration & Review

- Review Copilot output before committing.
- Refactor generated code to ensure readability and testability.
- Use comments to give Copilot context for better suggestions.
- Regenerate parts that are unidiomatic or too complex.

## 📚 References

- [Go Style Guide](https://google.github.io/styleguide/go/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Validator](https://github.com/go-playground/validator)
