export default function createStoryRepository() {
  return {
    // Pure data access: resolve a single story by slug. The controller supplies
    // the slug, so no router context leaks in here.
    getStory: async function (slug) {
      const response = await fetch(`/api/stories/${encodeURIComponent(slug)}`);

      if (!response.ok) {
        const err = new Error(`Response status: ${response.status}`);
        err.status = response.status;
        throw err;
      }

      const result = await response.json();
      return result.data;
    },

    // The initial request is server-rendered, so the story is already in the
    // DOM. The controller reads it once from #story-data.
    readInitialStory() {
      const data = document.getElementById("story-data");
      if (!data) {
        return null;
      }
      return JSON.parse(data.textContent);
    },
  };
}