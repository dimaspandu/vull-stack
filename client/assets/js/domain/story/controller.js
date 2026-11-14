import createStoryRepository from "./repository.js";
import createStoryView from "./view.jsx";
import createStoryViewModel from "./viewModel.js";

const _CACHE = {
  repository: null,
  viewModel: null,
  view: null,
};

export default async function showController(ctx, outlet) {
  if (
    ctx.isSameRoute ||
    ctx.isSameURL ||
    ctx.isSelfNavigate
  ) {
    await _CACHE.viewModel._load(ctx);
  } else {
    const repository = createStoryRepository();
    _CACHE.repository = repository;
    const viewModel = createStoryViewModel(repository);
    _CACHE.viewModel = viewModel;

    // Setiap kali route diakses / berubah, _load langsung dipanggil dengan ctx dari router
    await viewModel._load(ctx);

    const view = createStoryView(viewModel);
    _CACHE.view = view;

    if (ctx.firstArrive && ctx.isInitialRoute){
      elementBuilder(outlet).setChildren(view);
    } else {
      elementBuilder(outlet)
        .reset()
        .setAttributes({ class: "story-detail-page" })
        .setChildren(view);
    }
  }
}