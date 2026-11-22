import createStoriesRepository from "./repository.js";
import createStoriesView from "./view.jsx";
import createStoriesViewModel from "./viewModel.js";

import { setPageMeta } from "../../utils/meta.js";

const _CACHE = {
  viewModel: null,
};

export default async function listController(ctx) {
  // The list has one title and one description, so they are the same on every
  // load and can be set before the cached and fresh paths split.
  setPageMeta(
    "Short Stories",
    "A collection of short stories to read in a few quiet minutes."
  );

  // Back/forward or a self navigation to the same route: reuse the cached
  // viewModel so its listeners keep the already mounted view in sync.
  if (ctx.isSameRoute || ctx.isSameURL || ctx.isSelfNavigate) {
    await _CACHE.viewModel._load(ctx);
    return;
  }

  const outlet = document.getElementById("outlet");

  const repository = createStoriesRepository();
  const viewModel = createStoriesViewModel(repository);
  _CACHE.viewModel = viewModel;

  // Every route change reloads the state through the router context.
  await viewModel._load(ctx);

  const view = createStoriesView(viewModel);

  if (ctx.firstArrive && ctx.isInitialRoute) {
    // The server already rendered this markup, so only mount the view.
    document.createElement(
      outlet,
      { class: "story-page" },
      view
    );
  } else {
    // Client-side navigation: drop the previous DOM before mounting the new view.
    document.createElement(outlet)
      .setAttributes({ class: "story-page" })
      .setChildren(view);
  }
}