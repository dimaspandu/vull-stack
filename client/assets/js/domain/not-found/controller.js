import createNotFoundView from "./view.jsx";

export default function notFoundController(_) {
  const outlet = document.getElementById("outlet");
  const view = createNotFoundView();

  // No route matched, so there is no server-rendered markup to hydrate.
  document.createElement(outlet)
    .setAttributes({ class: "story-page" })
    .setChildren(view);
}