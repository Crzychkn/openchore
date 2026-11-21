# Deployment boundary

The live Ubuntu checkout at `/home/arthur/openchore` holds real household
data. The M1 work uses a separate development machine or clone with
synthetic data. Keep the live checkout, running services, SQLite database,
Docker volumes, credentials and Cloudflare configuration outside the
development workflow.

The new Compose file is for development. PR checks run in GitHub Actions.
Merging a customization PR into `custom` has no application deployment
step in the repository workflows. The inherited image and Pages workflows
have their own `main` triggers. Any automation configured outside the
repository is outside this audit.

## Future production releases

An operator should authorize a separate release and select a reviewed commit
from `custom`. Verify the candidate in the isolated stack with a synthetic
upgrade, document the migration and rollback requirements, and choose
immutable build artifacts for that revision. Before an approved production
upgrade, establish a consistent backup and verify its restore procedure in an
isolated environment. Include SQLite and persistent media and keep production
secrets in their existing operational boundary. A rollback decision must
account for database schema compatibility as well as application versions.

M1 makes no changes to production and establishes no automatic production
deployment.
