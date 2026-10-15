# Security Policy

> **📋 Team Contact:** For security-related questions about team responsibilities, see [TEAM_STRUCTURE.md](./TEAM_STRUCTURE.md).

## Supported versions

Security fixes are backported to the versions below. Versions not listed no longer receive updates — please upgrade before reporting an issue against them.

| Version | Supported |
|---|---|
| 2.x.x   | ✅ |
| 1.x.x   | ✅ (critical fixes only) |
| < 1.0   | ❌ |

## Reporting a vulnerability

**Please do not open a public issue for security vulnerabilities.** Public issues are visible to everyone immediately, including anyone who might exploit the report before a fix ships.

Instead, report privately using one of these methods:

1. **GitHub Private Vulnerability Reporting (preferred)** — go to the **Security** tab of this repository → **Report a vulnerability**. This opens a private draft advisory visible only to maintainers, where we can discuss and coordinate a fix directly with you.
2. **Email** — send details to **security@your-domain.com** (update this with your actual security contact). Include:
   - A description of the vulnerability and its potential impact
   - Steps to reproduce, or a proof-of-concept if available
   - Affected version(s)
   - Any suggested remediation, if you have one

## What to expect

- **Acknowledgment** within 3 business days of your report.
- **Initial assessment** (severity, affected versions) within 7 days.
- We'll keep you updated as a fix is developed and coordinate a disclosure timeline with you — we generally aim to patch and publish within 90 days of the initial report, sooner for critical issues.
- Once a fix is released, we'll publish a GitHub Security Advisory crediting you (unless you'd prefer to remain anonymous) and referencing the CVE if one was assigned.

## Disclosure policy

We follow **coordinated disclosure**: details are kept private until a patch is available, at which point we publish an advisory alongside the fix. We ask reporters not to publicly disclose the vulnerability before then.

## Scope

This policy covers the code in this repository. Vulnerabilities in third-party dependencies should be reported to the maintainers of those projects directly, though we'd appreciate a heads-up so we can track and update the dependency on our end.
