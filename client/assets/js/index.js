import "./libs/historypp/src/1.0.7/index.js";
import "./libs/dompp/src/index.js";

import listHandler from "./domain/stories/handler.js";
import showHandler from "./domain/story/handler.js";
import {
  normalizePath
} from "./utils/common.js";

const outlet = document.getElementById("outlet");

history.router("/", (ctx) => listHandler(ctx, outlet));
history.router("/:slug", (ctx) => showHandler(ctx, outlet));

history.notFound((ctx) => {});

const initialPath = normalizePath(location.pathname);
history.navigateReplace(initialPath + location.search);