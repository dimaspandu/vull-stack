import createStoriesRepository from "./repository.js";
import createStoriesView from "./view.jsx";
import createStoriesViewModel from "./viewModel.js";

const _CACHE = {
  repository: null,
  viewModel: null,
  view: null,
};

export default async function listController(ctx, outlet) {
  if (
    ctx.isSameRoute ||
    ctx.isSameURL ||
    ctx.isSelfNavigate
  ) {
    await _CACHE.viewModel._load(ctx);
  } else {
    const repository = createStoriesRepository();
    _CACHE.repository = repository;
    const viewModel = createStoriesViewModel(repository);
    _CACHE.viewModel = viewModel;

    // Setiap kali route diakses / berubah, _load langsung dipanggil dengan ctx dari router
    await viewModel._load(ctx);

    const view = createStoriesView(viewModel);
    _CACHE.view = view;

    if (ctx.firstArrive && ctx.isInitialRoute){
      // hydrate!
      elementBuilder(
        outlet,
        { class: "story-page" },
        view
      );
    } else {
      // remount from scratch!
      elementBuilder(outlet)
        .reset()
        .setAttributes({ class: "story-page" })
        .setChildren(view);
    }
  }
}