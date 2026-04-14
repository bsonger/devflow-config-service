# devflow-config-service deploy

This deployment expects the image built from `Dockerfile.staging` so `/app/config-repo` exists inside the container.

`POST /api/v1/app-configs/{id}/sync-from-repo` now runs `git pull --ff-only origin main` before freezing a revision.

Because the fixed repo is `git@github.com:bsonger/devflow-config-repo.git`, the runtime container also needs SSH credentials that can read GitHub, for example:
- a mounted `/home/devuser/.ssh` with private key + `known_hosts`
- or an injected `SSH_AUTH_SOCK`/agent-compatible setup

Without runtime GitHub credentials, sync will fail at the pull step with `424 failed_precondition`.
