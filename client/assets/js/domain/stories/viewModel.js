export default function createStoriesViewModel(repository) {
  const listeners = new Set();

  // Full list for the current route; filtering never mutates this.
  let allStories = [];

  const viewModel = {
    stories: [],
    error: null,
    query: "",

    // Called by the controller on every route change. The underscore marks it
    // as internal: the view only subscribes through addListener().
    async _load(ctx) {
      this.error = null;

      try {
        allStories = await repository.getStories(ctx);
        this.stories = allStories;

        // A `?q=` query means the server already filtered the list.
        if (ctx?.query?.q) {
          this.query = ctx.query.q;
          this.filterStories(this.query);
        }
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

    // Form submit only updates the URL; the router then reloads the list.
    handleSearchSubmit(event) {
      event.preventDefault();
      const formData = new FormData(event.target);
      const query = formData.get("q")?.trim() || "";
      history.navigatePush(`/?q=${encodeURIComponent(query)}`);
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