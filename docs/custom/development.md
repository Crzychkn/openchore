# Isolated Docker development

Use Docker with Compose V2 on a development machine and a separate clone.
Use the standalone [Compose file](../../compose.custom-dev.yaml) with
synthetic household data.

```bash
mkdir -p "$HOME/src"
git clone --branch custom https://github.com/Crzychkn/openchore.git "$HOME/src/openchore-custom-dev"
cd "$HOME/src/openchore-custom-dev"

dev_compose() {
  docker compose --env-file /dev/null -p openchore-custom-dev -f compose.custom-dev.yaml "$@"
}
dev_compose config
dev_compose up -d --build --wait
dev_compose ps
curl --fail --silent --show-error http://127.0.0.1:18080/api/users
```

Open http://127.0.0.1:18080 and create a synthetic household in the setup
wizard. Check profile selection and Today, then run `dev_compose restart`
and verify that the profiles remain. Rebuild with `dev_compose up -d --build
--wait` after source changes. `dev_compose logs api web` shows startup
and application logs. `dev_compose down` retains the development database.

- The explicit Compose project owns its network, volume and generated
  container and image names. For a parallel checkout, choose a distinct
  development project name in the helper and a free host port.
- The only published port is web on loopback. Override `OPENCHORE_DEV_PORT`
  to change the default `18080`. Nginx accesses API port `8080` over the
  Compose network. The API health check gates initial web startup.
- `openchore-custom-dev_dev-data` holds SQLite at `/data/openchore.db`.
  The API working directory is `/`, so its relative chore-photo and TTS
  paths also resolve under `/data`. The container's config path expects
  optional seed data and defaults to the setup flow with a new database.
  For a demo seed, a local Compose override may mount
  `./config/config.example.yaml:/app/config/config.yaml:ro` on the API.
- The helper skips automatic `.env` loading. The API environment is
  explicit, with a development time zone override via `OPENCHORE_DEV_TZ`.
  There is no host config mount or production credential passthrough in
  the default development file. AI, TTS and Cloudflare are unnecessary
  for this stack.

For the full source test set, use Go and Node versions from the workflow.
Run E2E from a disposable development checkout. Playwright removes the
checkout-local test database during server startup. Use free ports and
`CI=1` to prevent Playwright from reusing a running server.

```bash
go test ./...
(cd web && npm ci && npm test && npm run build)
(cd e2e && npm ci && npx playwright install chromium)
CI=1 E2E_API_PORT=18081 E2E_WEB_PORT=15173 make test-e2e
```
