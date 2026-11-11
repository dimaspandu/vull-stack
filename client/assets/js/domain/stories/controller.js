import createStoriesRepository from "./repository.js";
import createStoriesView from "./view.jsx";
import createStoriesViewModel from "./viewModel.js";

export default async function listController(ctx, outlet) {
  const repository = createStoriesRepository();
  const viewModel = await createStoriesViewModel(ctx, repository);
  const {
    mainView,
    storyListView,
  } = createStoriesView(viewModel);

  viewModel.addListener(() => storyListView.rerender());

  if (ctx.firstArrive && ctx.isInitialRoute){
    elementBuilder(outlet).setChildren(mainView);
  } else {
    elementBuilder(outlet)
      .reset()
      .setAttributes({ class: "story-page" })
      .setChildren(mainView);
  }
}