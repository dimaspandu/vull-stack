// Router-agnostic state for the story list. It owns the loaded stories, the
// active keyword, and the listener fan-out; the controller supplies data and
// the navigation callback, so nothing in here imports historypp or the DOM.
export default function createStoriesViewModel(repository, navigate) {
  const listeners = new Set();

  // Full list for the current route; filtering never mutates this.
  let allStories = [];

  const viewModel = {
    stories: [],
    error: null,
    query: "",

    // Populate state from an already-resolved data set. The controller decides
    // whether that set came from the server-rendered markup or the API.
    _hydrate(initialStories, query) {
      this.error = null;
      allStories = initialStories;
      this.query = query || "";
      this.filterStories(this.query);
      this.notifyListener();
    },

    // Reload state through the repository on every route change. The controller
    // passes a plain query value instead of the router context.
    async _load(query) {
      this.error = null;

      try {
        allStories = await repository.getStories(query);
        this.query = query || "";
        this.filterStories(this.query);
      } catch (err) {
        this.error = err.message;
        this.stories = [];
      } finally {
        this.notifyListener();
      }
    },

    // Client-side search over the stories loaded for this route.
    searchStories(keyword) {
      this.query = keyword;
      this.filterStories(keyword);
      this.notifyListener();
    },

    // Form submit only updates the URL; the router then reloads the list. The
    // navigation callback is injected so the viewModel stays router-free.
    handleSearchSubmit(event) {
      event.preventDefault();
      const formData = new FormData(event.target);
      const query = formData.get("q")?.trim() || "";
      navigate(`/?q=${encodeURIComponent(query)}`);
    },

    filterStories(keyword) {
      if (!keyword.trim()) {
        this.stories = allStories;
        return;
      }
      const lower = keyword.toLowerCase();
      this.stories = allStories.filter(story =>
        story.title.toLowerCase().includes(lower) ||
        story.excerpt.toLowerCase().includes(lower) ||
        story.category.toLowerCase().includes(lower)
      );
    },

    addListener(listener) {
      listeners.add(listener);
      return () => listeners.delete(listener);
    },

    notifyListener() {
      listeners.forEach((listener) => listener(viewModel));
    },
  };

  return viewModel;
}