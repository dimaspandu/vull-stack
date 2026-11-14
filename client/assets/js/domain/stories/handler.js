export default async function listHandler(ctx) {
  return import("./controller.js")
    .then(({ default: listController }) =>
      listController(ctx));
}