export default function showHandler(ctx) {
  return import("./controller.js")
    .then(({ default: showController }) =>
      showController(ctx));
}