import "./libs/historypp/src/1.0.6/index.js";
import "./libs/dompp/src/1.1.1/index.js";

import listHandler from "./domain/stories/handler.js";
import showHandler from "./domain/story/handler.js";

document.addEventListener("click", (e) => {
  const a = e.target.closest("[data-link]");
  if (!a) return;

  const url = new URL(a.href);
  if (url.origin !== location.origin) return;

  e.preventDefault();

  const path = normalize(url.pathname);
  history.navigatePush(path);
});

const outlet = document.getElementById("outlet");

history.router("/", listHandler(outlet));
history.router("/:slug", showHandler(outlet));

history.notFound((ctx) => {});

const initialPath = normalize(location.pathname);
history.navigateReplace(initialPath + location.search);

function normalize(path) {
  if (path !== "/" && path.endsWith("/")) {
    return path.slice(0, -1);
  }
  return path;
}