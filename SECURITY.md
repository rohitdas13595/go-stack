# Security policy

## Supported versions

GoStack is under active development. Security fixes are applied to the default
branch (`main` or `master`) as they are confirmed and fixed. There is no
long-term support (LTS) line yet; if that changes, this file will list supported
tags.

## Reporting a vulnerability

**Please do not report security issues through public GitHub issues.**

Instead, use one of these options:

1. **GitHub private vulnerability reporting** (if enabled on the repository):
   open the repository on GitHub and use **Security → Advisories → Report a
   vulnerability**.
2. **Email** the maintainers at a dedicated security contact if one is published
   in the repository or organization profile. If none is listed, use the
   repository owner’s contact on their GitHub profile for a private message.

Include:

- A short description of the issue and its impact
- Steps to reproduce (proof-of-concept, requests, or code paths)
- Affected versions or commits, if known
- Your suggestion for a fix (optional)

We aim to acknowledge reports within a few business days and to coordinate
disclosure once a fix is available.

## Scope

In scope:

- The GoStack framework code in this repository (including `cmd/`, packages  under the module root, and maintained `examples/`).
- The `gostack` CLI as shipped here.

Out of scope (report to the relevant upstream instead):

- Vulnerabilities only in third-party dependencies — use `govulncheck` and
  follow advisories for those projects.
- Applications built *with* GoStack unless the flaw is in framework code itself.

## Disclosure

We prefer responsible disclosure: please allow time for a patch before public
details. We will credit reporters who wish to be named in advisory text unless
they prefer to stay anonymous.
