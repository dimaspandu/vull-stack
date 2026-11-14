export default function createStoryRepository() {
  return {
    getStory: async function (ctx) {
      // The initial request is server-rendered, so the story is already in the DOM.
      if (ctx.firstArrive && ctx.isInitialRoute) {
        const data = document.getElementById("story-data");
        if (data) {
          return JSON.parse(data.textContent);
        }
        return null;
      }

      const response = await fetch(`/api/stories/${encodeURIComponent(ctx.params.slug)}`);

      if (!response.ok) {
        throw new Error(`Response status: ${response.status}`);
      }

      const result = await response.json();
      return result.data;
    }
  };
}