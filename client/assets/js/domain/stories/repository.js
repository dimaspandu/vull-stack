export default function createStoriesRepository() {
  return {
    getStories: async function (ctx) {
      if (ctx.firstArrive && ctx.isInitialRoute){
        if (document.getElementById("stories-data")) {
          return JSON.parse(
            document.getElementById("stories-data").textContent
          );
        }
        return [];
      }

      const response = await fetch("/api/stories");
      if (response.ok) {
        const result = await response.json();
        const stories = result.data;
        return stories;
      }

      return [];
    }
  };
}