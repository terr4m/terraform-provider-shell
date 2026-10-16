# Contributing to Terraform Provider Shell

Thank you for your interest in contributing to the _Terr4m_ Shell TF provider! We welcome contributions in all forms including documentation, issues, and code.

## Getting Started

Before contributing, please take a moment to review this guide to understand our contribution process.

## Types of Contributions

We welcome the following types of contributions:

- **Bug reports and fixes**
- **Documentation improvements**
- **New features and enhancements**
- **Code quality improvements**
- **Test coverage improvements**

## Before You Start

### For Feature Requests

Before making significant changes, please **open an issue first** to discuss:

- The motivation for the change
- The proposed approach
- Impact on existing functionality
- Migration path for users (if applicable)

## Development Setup

1. Fork the repository and clone your fork.
2. Install the Go version specified in `go.mod` and Terraform CLI. The acceptance tests use the locally installed Terraform binary and need no cloud credentials.
3. Install `just` if you want to use the repository's task recipes. Install `golangci-lint` and `rumdl` to run the same Go and Markdown checks as CI.
4. Download dependencies with `go mod download` (and run it in `tools/` before generating docs). Use `go mod tidy` only when changing dependencies, then review the resulting `go.mod` and `go.sum` changes.

## Making Changes

1. Create a new branch for your changes
2. Make your changes following our coding standards
3. Add or update tests as needed
4. Update documentation if applicable
5. Ensure all tests pass
6. Commit your changes with clear, descriptive messages

## Testing and validation

From the repository root:

- `go test ./...` runs the Go test suite, including provider acceptance tests that exercise Terraform configurations and local scripts.
- `go test -race ./...` runs tests with the race detector where supported.
- `golangci-lint run --timeout=10m` runs the configured Go linters.
- `golangci-lint fmt` formats Go code using the configured formatters.
- `rumdl check .` checks Markdown.
- `cd tools && go mod download && go generate ./...` formats Terraform examples and regenerates Registry documentation from the provider schemas, examples, and templates. Edit the relevant schema descriptions or files under `templates/`, then regenerate docs; do not hand-edit generated `docs/` pages. Review and commit generated changes when relevant.
- `go build ./...` verifies that all packages build.

Add or update unit and acceptance coverage for behavior changes. Before submitting, run the checks relevant to your changes and ensure generated documentation is current.

## Documentation

- Update provider schema descriptions and/or documentation templates as appropriate, then regenerate Registry docs; do not hand-edit generated pages.
- Ensure code comments are clear and helpful
- Follow the existing documentation style and format

## Submitting Your Contribution

1. Push your changes to your fork
2. Open a pull request against the main branch
3. Provide a clear description of your changes
4. Reference any related issues
5. Ensure CI checks pass

## Code Review Process

- All submissions require review before merging
- Reviewers may request changes or ask questions
- Please be responsive to feedback and willing to iterate
- Maintainers will merge approved changes

## Community Guidelines

- Be respectful and constructive in all interactions
- Follow the project's [Code of Conduct](https://github.com/terr4m/.github/blob/main/CODE_OF_CONDUCT.md).
- Help others and share knowledge
- Report any issues or concerns to maintainers

## Questions?

If you have questions about contributing, feel free to:

- Open an issue for discussion
- Review existing issues and pull requests
- Reach out to maintainers

Thank you for helping make this project better!
