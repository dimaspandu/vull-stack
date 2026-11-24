import createStoriesRepository from "./repository.js";
import createStoriesView from "./view.jsx";
import createStoriesViewModel from "./viewModel.js";

import { mountOutlet } from "../../utils/mount.js";
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

  const isServerRendered = ctx.firstArrive && ctx.isInitialRoute;

  // Back/forward or a self navigation to the same route: reuse the cached
  // viewModel so its listeners keep the already mounted view in sync. Only the
  // keyword can change, and the server-rendered markup already carries it.
  if (ctx.isSameRoute || ctx.isSameURL || ctx.isSelfNavigate) {
    const query = ctx.query?.q;
    if (isServerRendered) {
      _CACHE.viewModel._hydrate(
        _CACHE.viewModel.stories,
        query
      );
    } else {
      await _CACHE.viewModel._load(query);
    }
    return;
  }

  const repository = createStoriesRepository();
  const viewModel = createStoriesViewModel(repository, history.navigatePush);
  _CACHE.viewModel = viewModel;

  // The initial request is server-rendered, so seed the viewModel from the
  // markup instead of refetching it. Every later route change reloads.
  if (isServerRendered) {
    viewModel._hydrate(
      repository.readInitialStories(),
      ctx.query?.q
    );
  } else {
    await viewModel._load(ctx.query?.q);
  }

  const view = createStoriesView(viewModel, history.navigatePush);

  mountOutlet(view, { hydrate: isServerRendered });
}