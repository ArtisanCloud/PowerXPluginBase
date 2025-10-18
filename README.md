# PowerX Plugin Lifecycle Resources

Key references for getting a new plugin repository lifecycle-ready:

- Start with the detailed [Lifecycle Overview](docs/lifecycle/overview.md)
- Follow the [Bootstrap Quickstart](docs/lifecycle/quickstart.md) for the abbreviated checklist
- Complete the [Bootstrap Compliance Checklist](docs/lifecycle/checklists/bootstrap-checklist.md)
- Sync docs for reviewers with `make sync-lifecycle-docs`

## Security & Compliance Additions

This plugin ships the privacy, ToolGrant, audit, and vulnerability response capabilities described in `specs/004-security-compliance/`:

- Follow the [Security & Compliance Quickstart](specs/004-security-compliance/quickstart.md) to bootstrap databases, masking rules, and advisory drills.
- Review the [Vulnerability Response Playbook](docs/security/vulnerability-response.md) for intake channels, SLA expectations, and packaging workflows.
- Consult `docs/security/audit-logs.md` for audit export tooling and retention policies.

For the full documentation tree, see [docs/readme.md](docs/readme.md).
