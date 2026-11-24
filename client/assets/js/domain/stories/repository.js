export default function createStoriesRepository() {
  // The initial request is server-rendered, so the list is already in the DOM.
  // The controller reads it once from #stories-data; the repository itself is
  // route-agnostic and only knows how to fetch the story list.
  function readInitialStories() {
    const data = document.getElementById("stories-data");
    if (!data) {
      return [];
    }
    return JSON.parse(data.textContent);
  }

  return {
    // Pure data access: given a keyword (or none), resolve the story list.
    // No router context leaks in here; the controller normalizes it.
    getStories: async function (query) {
      const remoteUrl = query
        ? `/api/stories?q=${encodeURIComponent(query)}`
        : "/api/stories";

      const response = await fetch(remoteUrl);
      if (response.ok) {
        const result = await response.json();
        return result.data;
      }

      return [];
    },

    readInitialStories,
  };
}