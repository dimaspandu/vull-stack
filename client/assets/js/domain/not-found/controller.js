import createNotFoundView from "./view.jsx";

import { mountOutlet } from "../../utils/mount.js";
import { setPageMeta } from "../../utils/meta.js";

export default function notFoundController() {
  const view = createNotFoundView(history.navigatePush);

  setPageMeta("404 — Page not found", "404 — Page not found");

  // No route matched, so there is no server-rendered markup to hydrate.
  mountOutlet(view, { hydrate: false });
}