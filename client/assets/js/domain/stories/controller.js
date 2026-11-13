import createStoriesRepository from "./repository.js";
import createStoriesView from "./view.jsx";
import createStoriesViewModel from "./viewModel.js";

export default async function listController(ctx, outlet) {
  const repository = createStoriesRepository();
  const viewModel = createStoriesViewModel(repository);

  // Setiap kali route diakses / berubah, _load langsung dipanggil dengan ctx dari router
  await viewModel._load(ctx);

  const view = createStoriesView(viewModel);

  if (ctx.firstArrive && ctx.isInitialRoute){
    elementBuilder(outlet).setChildren(view);
  } else {
    elementBuilder(outlet)
      .reset()
      .setAttributes({ class: "story-page" })
      .setChildren(view);
  }
}