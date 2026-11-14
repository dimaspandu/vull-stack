import "./libs/historypp/src/1.0.7/index.js";
import "./libs/dompp/src/index.js";

import { withLoading } from "./components/loading/controller.js";

import listHandler from "./domain/stories/handler.js";
import showHandler from "./domain/story/handler.js";
import notFoundHandler from "./domain/not-found/handler.js";
import {
  normalizePath
} from "./utils/common.js";

// withLoading keeps the transition indicator in one place: every route shows
// it before the handler runs and hides it when the handler settles.
history.router("/", withLoading(listHandler));
history.router("/:slug", withLoading(showHandler));

history.notFound(withLoading(notFoundHandler));

// The server already rendered the first URL, so replay the router against it
// to hydrate the matching domain instead of leaving the SSR markup untouched.
const initialPath = normalizePath(location.pathname);
history.navigateReplace(initialPath + location.search);