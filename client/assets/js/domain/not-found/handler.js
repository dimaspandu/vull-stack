export default function notFoundHandler(ctx) {
  return import("./controller.js")
    .then(({ default: notFoundController }) =>
      notFoundController(ctx));
}