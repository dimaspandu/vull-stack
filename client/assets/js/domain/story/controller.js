import createStoryRepository from "./repository.js";
import createStoryView from "./view.jsx";
import createStoryViewModel from "./viewModel.js";

export default async function showController(ctx, outlet) {
  const repository = createStoryRepository();
  const viewModel = createStoryViewModel(repository);

  // Setiap kali route diakses / berubah, _load langsung dipanggil dengan ctx dari router
  await viewModel._load(ctx);

  const view = createStoryView(viewModel);

  if (ctx.firstArrive && ctx.isInitialRoute){
    elementBuilder(outlet).setChildren(view);
  } else {
    elementBuilder(outlet)
      .reset()
      .setAttributes({ class: "story-detail-page" })
      .setChildren(view);
  }
}