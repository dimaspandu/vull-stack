export default function createStoryRepository() {
  return {
    getStory: async function (ctx) {
      if (ctx.firstArrive && ctx.isInitialRoute){
        if (document.getElementById("story-data")) {
          return JSON.parse(
            document.getElementById("story-data").textContent
          );
        }
        return null;
      }

      const response = await fetch(`/api/stories/${ctx.params.slug}`);

      if (!response.ok) {
        throw new Error(`Response status: ${response.status}`);
      }

      const result = await response.json();
      return result.data;
    }
  };
}