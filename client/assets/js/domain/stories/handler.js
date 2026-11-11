export default function listHandler(ctx, outlet) {
  import("./controller.js")
    .then(({ default: listController }) => 
      listController(ctx, outlet));
}