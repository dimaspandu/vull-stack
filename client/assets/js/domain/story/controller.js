import createStoryRepository from "./repository.js";
import createStoryView from "./view.jsx";
import createStoryViewModel from "./viewModel.js";

import { mountOutlet } from "../../utils/mount.js";
import { setPageMeta } from "../../utils/meta.js";

const _CACHE = {
  viewModel: null,
};

// Mirrors the server-side template: the detail page titles itself after the
// story it renders, and falls back to a not-found title when there is none.
function syncPageMeta(viewModel) {
  const story = viewModel.story;

  if (!story) {
    setPageMeta(
      "Story not found",
      "The story you are looking for does not exist or has been moved."
    );
    return;
  }

  setPageMeta(
    `${story.title} — Short Stories`,
    `${story.title} — ${story.excerpt}`
  );
}

export default async function showController(ctx) {
  const isServerRendered = ctx.firstArrive && ctx.isInitialRoute;

  // Back/forward or a self navigation to the same route: reuse the cached
  // viewModel so its listeners keep the already mounted view in sync.
  if (ctx.isSameRoute || ctx.isSameURL || ctx.isSelfNavigate) {
    if (isServerRendered) {
      _CACHE.viewModel._hydrate(
        _CACHE.viewModel.story
      );
    } else {
      await _CACHE.viewModel._load(ctx.params.slug);
    }
    syncPageMeta(_CACHE.viewModel);
    return;
  }

  const repository = createStoryRepository();
  const viewModel = createStoryViewModel(repository);
  _CACHE.viewModel = viewModel;

  // The initial request is server-rendered, so seed the viewModel from the
  // markup instead of refetching it. Every later route change reloads.
  if (isServerRendered) {
    viewModel._hydrate(repository.readInitialStory());
  } else {
    await viewModel._load(ctx.params.slug);
  }

  syncPageMeta(viewModel);

  const view = createStoryView(viewModel, history.navigatePop);

  mountOutlet(view, { hydrate: isServerRendered });
}