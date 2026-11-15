# vull-stack

A small Go application that renders pages on the server and then takes over as a
single page application in the browser: server-side rendering for the first
paint, client-side routing for every navigation after it.

No bundler, no transpiler step, no `node_modules`. The Go binary is the only
build artifact; `.jsx` files are compiled on demand while the assets are served.

## Features

- **SSR first.** `/` and `/{slug}` are rendered with `html/template` from MySQL,
  so the first paint is complete HTML.
- **SPA after hydration.** The same HTML is reused as the initial state; later
  navigations are handled in the browser by `historypp` and never hit the
  document.
- **JSX without a build step.** `internal/jsxtranspiler` compiles `.jsx` on the
  fly when the asset is requested and serves it as plain JavaScript.
- **MVVM per domain.** Each route owns a `handler → controller → repository →
  viewModel → view` chain, with the view split into `componentBuilder` slices.
- **Route transition indicator.** A reference-counted loading indicator is
  mounted once into `#portal` and driven by the router.
- **Search.** Server-side filtering through `?q=`, plus live client-side
  filtering over the stories already loaded for the route.

## Tech stack

| Layer | Choice |
| --- | --- |
| Language | Go 1.25 |
| HTTP | `net/http` (`ServeMux` patterns, `CrossOriginProtection`) |
| Database | MySQL (`go-sql-driver/mysql`) |
| Config | `.env` via `joho/godotenv` |
| Templates | `html/template` |
| Client router | [`historypp`](client/assets/js/libs/historypp) |
| Client DOM/JSX runtime | [`dompp`](client/assets/js/libs/dompp) |
| Tests | `go test` |

`historypp` and `dompp` are git submodules, so they are versioned separately
from this repository and must not be edited here.

## Architecture

```
browser
  │  GET /                GET /{slug}
  ▼
main.go ──► services/stories/web   (handler → usecase → repository)  ──► MySQL
  │
  ├─► html/template: client/index.html, client/[id].html
  └─► /assets/*: static files, .jsx compiled by internal/jsxtranspiler

browser (after hydration)
  index.js ──► historypp router ──► domain/{stories,story,not-found}/handler.js
                                        │
                                        └─► controller.js ──► repository.js ──► /api/stories
                                              │
                                              └─► viewModel.js ──► view.jsx
```

The web and API layers are separate on purpose: `services/stories/web` returns
HTML, `services/stories/api` returns JSON, and both share the same table and the
same `handler → usecase → repository` shape.

## Project structure

```
main.go                       asset serving, wiring, HTTP server
internal/jsxtranspiler/       JSX → JS compiler used at request time
services/stories/web/         HTML use case (list, show)
services/stories/api/         JSON use case (list, show)
database/schema.sql           table definition
database/seed.sql             sample rows
client/index.html             list template (SSR)
client/[id].html              detail template (SSR)
client/assets/js/index.js     router bootstrap
client/assets/js/domain/      one folder per route: handler, controller,
                              repository, viewModel, view
client/assets/js/components/  reusable UI (loading indicator)
client/assets/js/utils/       shared helpers
client/assets/css/            styles
client/assets/js/libs/        dompp, historypp (submodules)
```

## Getting started

Requirements: Go 1.25+ and a reachable MySQL server.

```bash
cp .env.example .env
```

Create the schema and load the sample data:

```bash
mysql -u root -p < database/schema.sql
mysql -u root -p vss < database/seed.sql
```

Run the server:

```bash
go run .
```

Then open <http://127.0.0.1:6060>.

On Windows with WDAC or Smart App Control enabled, `go run` can be blocked by
policy. Build the binary and run it directly instead:

```powershell
go build -o main.exe .
.\main.exe
```

## Configuration

| Variable | Default | Description |
| --- | --- | --- |
| `DB_HOST` | `127.0.0.1` | MySQL host |
| `DB_PORT` | `3306` | MySQL port |
| `DB_USER` | `root` | MySQL user |
| `DB_PASS` | empty | MySQL password |
| `DB_NAME` | `vss` | Database name |
| `APP_PORT` | `6060` | HTTP port |

`.env` is gitignored; `.env.example` is the template.

## Routes

| Method | Path | Type | Description |
| --- | --- | --- | --- |
| `GET` | `/` | HTML | Story list, optional `?q=` filter |
| `GET` | `/{slug}` | HTML | Story detail |
| `GET` | `/api/stories` | JSON | Story list, optional `?q=` filter |
| `GET` | `/api/stories/{slug}` | JSON | Story detail |
| `GET` | `/assets/...` | static | CSS, images, plain JS |
| `GET` | `/assets/**/*.jsx` | JS | `.jsx` compiled on the fly |
| `GET` | `/favicon.ico` | 404 | Not served |

API responses are wrapped: `{"data": ...}`.

## How the SPA layer works

**Bootstrap.** `client/assets/js/index.js` registers the routes, then replays the
current URL with `history.navigateReplace(...)`. The first run is flagged
`firstArrive && isInitialRoute`, which tells every layer that the server already
rendered this page.

**SSR handoff.** The server embeds the data it used for the render in a
`<script type="application/json">` block (`#stories-data` on the list,
`#story-data` on the detail). On the initial route the repository reads that
block instead of calling the API, so the first paint costs no extra request.

**Mounting.** `#outlet` is the view container and is replaced on client-side
navigation. `#portal` lives outside it, so the loading indicator survives a route
change.

**View state.** `viewModel` is the only mutable state; `view.jsx` subscribes
through `addListener` and each `componentBuilder` slice re-renders on its own
notification.

**Loading.** `components/loading` owns a singleton viewModel rendered into
`#portal`. `withLoading(handler)` wraps every route handler, shows the indicator
before the handler runs and hides it when the handler settles. Show/hide are
reference counted so overlapping navigations do not hide each other's spinner.
The initial server-rendered route is skipped, because there is no transition to
indicate.

**JSX pipeline.** A request for `/assets/**/*.jsx` is read from disk, compiled by
`internal/jsxtranspiler` with `elementBuilder` as the factory, and returned as
`application/javascript` with `Cache-Control: no-cache`. Everything else under
`/assets/` is served as a static file.

## Tests

```bash
go test ./...
```

The transpiler is the only unit-tested package today
(`internal/jsxtranspiler`), including JSX comments inside opening tags.

## Notes

- The transpiler is a development convenience: it runs on every request for a
  `.jsx` file and intentionally sends `Cache-Control: no-cache`.
- `client/assets/js/factories/elementBuilder.js` and the `bindLinkNavigation`
  helper are currently unused; the runtime builder comes from `dompp`.
- `/api/stories/{slug}` answers 500 for an unknown slug instead of 404.
