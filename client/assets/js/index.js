import "./libs/historypp/src/1.0.6/index.js";
import "./libs/dompp/src/1.1.1/index.js";

import listHandler from "./domain/stories/handler.js";
import showHandler from "./domain/story/handler.js";
import {
  normalizePath
} from "./utils/common.js";

const outlet = document.getElementById("outlet");

history.router("/", listHandler(outlet));
history.router("/:slug", showHandler(outlet));

history.notFound((ctx) => {});

const initialPath = normalizePath(location.pathname);
history.navigateReplace(initialPath + location.search);