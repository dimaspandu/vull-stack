export default function createStoriesViewModel(repository) {
  const listeners = new Set();
  let allStories = [];

  const viewModel = {
    stories: [],
    loading: false,
    error: null,
    query: "",

    async _load(ctx) {
      this.loading = true;
      this.notifyListener();

      try {
        allStories = await repository.getStories(ctx);
        this.stories = allStories;
        
        if (ctx?.query?.q) {
          this.query = ctx.query.q;
          this.filterStories(this.query);
        }
      } catch (err) {
        this.error = err.message;
        this.stories = [];
      } finally {
        this.loading = false;
        this.notifyListener();
      }
    },

    // 1. Core logic pencarian & filter string
    searchStories(keyword) {
      this.query = keyword;
      this.filterStories(keyword);
      this.notifyListener();
    },

    // 2. Event handler khusus form submit (bersih dari urusan view)
    handleSearchSubmit(event) {
      event.preventDefault();
      const formData = new FormData(event.target);
      const query = formData.get("q")?.trim() || "";
      history.navigatePush(`/?q=${query}`);
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