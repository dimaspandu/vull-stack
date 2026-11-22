import createNotFoundView from "./view.jsx";

import { setPageMeta } from "../../utils/meta.js";

export default function notFoundController(ctx) {
  const outlet = document.getElementById("outlet");
  const view = createNotFoundView();

  setPageMeta("404 — Page not found", "404 — Page not found");

  // No route matched, so there is no server-rendered markup to hydrate.
  document.createElement(outlet)
    .setAttributes({ class: "story-page" })
    .setChildren(view);
}