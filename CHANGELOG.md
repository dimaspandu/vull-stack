# Changelog

## [1.0.1] - 2026-11-28

### Added
- `client/500.html`: server error page, served with status 500 by `services/stories/web.InternalServerError`.
- `services/admin/`: placeholder domain for content management (backoffice).
- `ErrorState` on the stories list view, rendered when `viewModel.error` is set.

### Changed
- Renamed `client/[id].html` to `client/[slug].html` to match the `/{slug}` route.
- Story detail 404 now renders `EmptyState` ("Story not found") instead of `ErrorState`; only true server errors render `ErrorState`.
- Extracted `formatPublishedAt` helper in both `services/stories/api` and `services/stories/web` usecases.
- `StoryLayout` now wraps detail/error/empty states with `story-detail-page__container` + back navigation for consistent presentation.
- `repository.getStory` now carries the HTTP status on error so the client can distinguish 404 from 500.

## [1.0.0] - 2026-11-24

### Added
- Initial release of `vull-stack`: SSR-first Go application that hydrates into a
  single page application with no bundler, no build step, and no `node_modules`.
- `internal/jsxtranspiler`: on-the-fly `.jsx` → JS compiler supporting the
  `@jsx`, `@jsxFrag`, and `@jsxRuntime` pragmas, block/line comments, and
  variadic `CompileWithFactory`.
- `services/stories/web` and `services/stories/api`: HTML and JSON layers over
  the same MySQL table, each with `handler → usecase → repository`, returning
  404 for unknown paths/slugs instead of 500.
- Client SPA layer with `historypp` (router) and `dompp` (JSX runtime), one
  domain folder per route (`stories`, `story`, `not-found`) with
  `handler → controller → repository → viewModel → view`.
- `utils/mount.js`: centralized `mountOutlet(view, {hydrate})` that picks
  hydration (server-rendered markup) vs full replace (client navigation).
- `components/loading`: singleton reference-counted indicator in `#portal`,
  wrapped by `withLoading(handler)`; initial SSR route is skipped.
- Scroll-to-top (`window.scrollTo(0, 0)`) on every client-side navigation.
- Search: server-side `?q=` filtering plus live client-side filtering over the
  route's already-loaded stories.
- Page metadata sync (`utils/meta.js`) across client-side navigations.
- All client assets served with `Cache-Control: no-store` so refactors never
  leave stale modules in the browser cache.

### Changed
- Refactored the client MVVM layer to be cleanly layered: controllers no longer
  touch the DOM directly, repositories no longer read router context, views and
  viewModels no longer import `history`.
- README rewritten with mermaid architecture, MVVM, and sequence diagrams.

### Removed
- Dropped `hashttp` and `ngapack` git submodules; `.gitmodules` now lists only
  `historypp` and `dompp`.

### License
MIT, Copyright (c) 2026 Dimas Pandu Pratama.