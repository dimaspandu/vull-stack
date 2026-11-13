export default function showHandler(ctx, outlet) {
  import("./controller.js")
    .then(({ default: showController }) =>
      showController(ctx, outlet));
}