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

1. Fork the repository
2. Clone your fork locally
3. Install Go (version specified in `go.mod`)
4. Install Terraform CLI
5. Run `go mod tidy` to install dependencies

## Making Changes

1. Create a new branch for your changes
2. Make your changes following our coding standards
3. Add or update tests as needed
4. Update documentation if applicable
5. Ensure all tests pass
6. Commit your changes with clear, descriptive messages

## Testing

- Run `go test ./...` to execute the test suite
- Add unit tests for new functionality
- Consider adding acceptance tests for new resources/data sources
- Ensure your changes don't break existing functionality

## Documentation

- Update the relevant documentation for any new features
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
- Follow the project's code of conduct
- Help others and share knowledge
- Report any issues or concerns to maintainers

## Questions?

If you have questions about contributing, feel free to:

- Open an issue for discussion
- Review existing issues and pull requests
- Reach out to maintainers

Thank you for helping make this project better!
