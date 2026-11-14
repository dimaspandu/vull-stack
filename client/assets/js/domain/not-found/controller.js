import createNotFoundView from "./view.jsx";

export default function notFoundController(ctx) {
  const outlet = document.getElementById("outlet");
  const view = createNotFoundView();

  // No route matched, so there is no server-rendered markup to hydrate.
  elementBuilder(outlet)
    .reset()
    .setAttributes({ class: "story-page" })
    .setChildren(view);
}