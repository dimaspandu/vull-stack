export default async function showHandler(ctx, outlet) {
  if (ctx.firstArrive && ctx.isInitialRoute) {
    const { default: view } = await import("../story/view.jsx");

    const story = JSON.parse(
      document.getElementById("story-data").textContent
    );

    elementBuilder(outlet).setChildren(
      view({
        category: story.category,
        title: story.title,
        content: story.content,
        publishedAt: story.publishedAt,
        readingTime: story.readingTime,
      })
    );
  } else {
    const { default: view } = await import("./view.jsx");

    try {
      const response = await fetch(`/api/stories/${ctx.params.slug}`);
      if (!response.ok) {
        throw new Error(`Response status: ${response.status}`);
      }
      const result = await response.json();
      const story = result.data;

      elementBuilder(outlet)
        .reset()
        .setAttributes({ class: "story-detail-page" })
        .setChildren(
          view({
            category: story.category,
            title: story.title,
            content: story.content,
            publishedAt: story.publishedAt,
            readingTime: story.readingTime,
          })
        );
    } catch (error) {
      elementBuilder(outlet)
        .reset()
        .setAttributes({ class: "story-detail-page" })
        .setChildren(
          elementBuilder("h1", null, "Internal Server Error.")
        );
    }
  }
}