# Contributing to Rampage L1

Thank you for your interest in contributing to Rampage L1. This document outlines the process for contributing to this project and the standards we expect from contributors.

## Important: Intellectual Property Notice

By submitting any contribution to this repository, you agree that:

1. Your contribution is your original work or you have the right to submit it.
2. You grant Shea Patrick Kastl and the Rampage project a perpetual, worldwide, non-exclusive, royalty-free license to use, reproduce, modify, and distribute your contribution as part of the Rampage project under the Apache License 2.0.
3. You understand that Rampage L1, its architecture, consensus mechanism, governance model, and associated trademarks remain the proprietary intellectual property of Shea Patrick Kastl.
4. Your contribution does not grant you any rights to the Rampage name, brand, or trademarks.

## Contributor License Agreement (CLA)

All contributors must agree to the Rampage Contributor License Agreement before any pull request can be merged. By opening a pull request, you affirm that you have read and agree to the CLA terms described in this document and the NOTICE file.

## How to Contribute

### Reporting Issues

- Use the GitHub Issues tracker to report bugs or request features.
- Search existing issues before opening a new one to avoid duplicates.
- Use the provided issue templates and fill them out completely.
- For security vulnerabilities, do NOT open a public issue. See SECURITY.md.

### Development Workflow

1. Fork the repository.
2. Create a feature branch from `main`: `git checkout -b feat/your-feature-name`
3. Make your changes following the coding standards below.
4. Write or update tests to cover your changes.
5. Ensure all tests pass: `make test`
6. Commit your changes with clear, descriptive commit messages.
7. Push to your fork and open a Pull Request against `main`.

### Branch Naming Conventions

- `feat/` – new features
- `fix/` – bug fixes
- `docs/` – documentation changes
- `chore/` – maintenance, refactoring, tooling
- `test/` – test additions or corrections
- `security/` – security patches

### Commit Message Format

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <short summary>

[optional body]

[optional footer]
```

Example:
```
feat(governor): add vote weight decay for inactive validators

Implements exponential decay on vote weight for validators that have
not participated in the last 100 blocks.

Closes #42
```

### Coding Standards

- **Language:** Go (primary), follow standard Go formatting (`gofmt`)
- **Linting:** Run `golangci-lint run` before submitting
- **Testing:** Maintain or improve test coverage; unit tests required for all new logic
- **Documentation:** All exported functions, types, and packages must have GoDoc comments
- **Error handling:** Always handle errors explicitly; no silent failures
- **Security:** Avoid introducing dependencies with known CVEs; minimize external dependencies

### Pull Request Requirements

- PRs must reference an open issue (except minor typo/doc fixes)
- All CI checks must pass
- At least one review approval required before merge
- Keep PRs focused and minimal in scope
- Include test cases demonstrating the change
- Update relevant documentation

## What We Will Not Accept

- Changes that alter or weaken the core Rampage L1 consensus mechanism without explicit approval from the core team
- Changes to governance parameters that have not gone through the on-chain governance process
- Code that introduces backdoors, surveillance capabilities, or censorship mechanisms
- Contributions that violate the Apache 2.0 license or third-party license terms
- Contributions that infringe on third-party intellectual property

## Code of Conduct

All contributors are expected to follow our [Code of Conduct](CODE_OF_CONDUCT.md). Harassment, discrimination, or disrespectful behavior will result in permanent ban from the project.

## Questions

For general questions about contributing, open a Discussion on GitHub. For licensing or intellectual property questions, contact the project maintainer through the repository.

---

*Rampage L1 is maintained by Shea Patrick Kastl. Contributions are welcome under the terms described above.*
