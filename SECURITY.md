# Security Policy

## Supported Versions

Security updates are provided for the following versions of Rampage L1:

| Version | Supported          |
| ------- | ------------------ |
| main    | :white_check_mark: |
| < v1.5  | :x:                |

As the project is currently in pre-release (v1.5 prototype), all active development occurs on `main`. Security patches will be applied directly to the main branch and tagged accordingly.

## Reporting a Vulnerability

**DO NOT report security vulnerabilities through public GitHub issues.**

If you discover a security vulnerability in Rampage L1, please report it responsibly through one of the following private channels:

### Private Disclosure

1. **GitHub Security Advisories (Preferred):** Use the [Security Advisories](https://github.com/RampageRPM1/rampage/security/advisories/new) feature to privately report vulnerabilities directly to the maintainer.

2. **Direct Contact:** Reach out to the project maintainer, Shea Patrick Kastl, through the GitHub profile associated with this repository.

### What to Include in Your Report

Please include as much of the following information as possible:

- Type of vulnerability (e.g., buffer overflow, SQL injection, consensus bypass, double-spend, Sybil attack, etc.)
- Full paths of source file(s) related to the vulnerability
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact assessment: what an attacker could gain from exploiting this vulnerability
- Suggested remediation (optional but appreciated)

## Response Timeline

We are committed to the following response process:

- **Acknowledgment:** Within 72 hours of receiving your report
- **Initial Assessment:** Within 7 days
- **Remediation Plan:** Within 30 days for critical issues
- **Patch Release:** As quickly as feasible depending on severity

## Severity Classification

| Severity | Description | Examples |
|----------|-------------|----------|
| Critical | Consensus breaking, network halt, fund theft | Double-spend, validator key compromise |
| High | Significant impact to chain security or data integrity | Governance manipulation, block forging |
| Medium | Limited impact, requires specific conditions | DoS vectors, minor data leaks |
| Low | Minimal impact, informational | Configuration warnings, log exposure |

## Disclosure Policy

We follow a **coordinated disclosure** model:

1. Reporter submits vulnerability privately.
2. Maintainer acknowledges and begins investigation.
3. Fix is developed and tested privately.
4. Patch is released.
5. A security advisory is published crediting the reporter (unless anonymity is requested).
6. A minimum of 90 days will pass before full public disclosure for critical/high severity issues.

## Bug Bounty

Rampage L1 does not currently operate a formal bug bounty program. However, significant vulnerability discoveries may be acknowledged publicly (with permission) and recognized in the project's NOTICE file.

## Scope

### In Scope

- Consensus mechanism vulnerabilities
- Governance bypass or manipulation
- Validator key management issues
- IBC security vulnerabilities in Rampage's implementation
- Smart contract / module vulnerabilities
- Denial of service attacks against the network
- Data integrity issues

### Out of Scope

- Third-party library vulnerabilities not directly exploitable in Rampage context
- Social engineering attacks
- Physical security issues
- Issues in forked or derivative projects not maintained by this repository

## Security Best Practices for Contributors

- Never commit private keys, mnemonics, or API secrets to the repository
- Use environment variables for all sensitive configuration
- All cryptographic implementations must use audited, well-established libraries
- New consensus-critical code requires explicit review from the core maintainer before merge

---

*This security policy is maintained by Shea Patrick Kastl. Last reviewed: 2025.*
