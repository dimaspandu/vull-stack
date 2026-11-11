export default async function listHandler(ctx, outlet) {
  if (ctx.firstArrive && ctx.isInitialRoute) {
    const { default: view } = await import("./view.jsx");

    const stories = JSON.parse(
      document.getElementById("stories-data").textContent
    );

    elementBuilder(outlet).setChildren(
      view({ query: "", stories })
    );
  } else {
    const { default: view } = await import("./view.jsx");

    try {
      const response = await fetch("/api/stories");
      if (!response.ok) {
        throw new Error(`Response status: ${response.status}`);
      }
      const result = await response.json();
      const stories = result.data;

      elementBuilder(outlet)
        .reset()
        .setAttributes({ class: "story-page" })
        .setChildren(view({ query: "", stories }));
    } catch (error) {
      elementBuilder(outlet)
        .reset()
        .setAttributes({ class: "story-page" })
        .setChildren(
          elementBuilder("h1", null, "Internal Server Error.")
        );
    }
  }
}