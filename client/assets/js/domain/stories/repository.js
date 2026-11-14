export default function createStoriesRepository() {
  return {
    getStories: async function (ctx) {
      // The initial request is server-rendered, so the list is already in the DOM.
      if (ctx.firstArrive && ctx.isInitialRoute) {
        const data = document.getElementById("stories-data");
        if (data) {
          return JSON.parse(data.textContent);
        }
        return [];
      }

      const query = ctx.query?.q;
      const remoteUrl = query
        ? `/api/stories?q=${encodeURIComponent(query)}`
        : "/api/stories";

      const response = await fetch(remoteUrl);
      if (response.ok) {
        const result = await response.json();
        return result.data;
      }

      return [];
    }
  };
}