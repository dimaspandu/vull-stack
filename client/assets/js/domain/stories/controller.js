import createStoriesRepository from "./repository.js";
import createStoriesView from "./view.jsx";
import createStoriesViewModel from "./viewModel.js";

export default async function listController(ctx, outlet) {
  const repository = createStoriesRepository();
  const viewModel = await createStoriesViewModel(ctx, repository);
  const view = createStoriesView(viewModel);

  viewModel.addListener(() => view.rerender());

  if (ctx.firstArrive && ctx.isInitialRoute){
    elementBuilder(outlet).setChildren(view.render());
  } else {
    elementBuilder(outlet)
      .reset()
      .setAttributes({ class: "story-page" })
      .setChildren(view.render());
  }
}