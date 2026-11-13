export default function createStoriesViewModel(repository) {
  const listeners = new Set();

  const viewModel = {
    stories: [],
    loading: false,
    error: null,
    query: "",

    async _load(ctx) {
      this.loading = true;
      this.error = null;
      this.notifyListener();

      try {
        // Ambil query dari ctx router (misal ctx.query.q)
        if (ctx?.query?.q) {
          this.query = ctx.query.q;
        } else {
          this.query = "";
        }

        // Repository akan otomatis mendeteksi apakah ini firstArrive (SSR/DOM) atau fetch API biasa
        this.stories = await repository.getStories(ctx);
      } catch (err) {
        this.error = err.message || "Failed to load stories";
        this.stories = [];
      } finally {
        this.loading = false;
        this.notifyListener();
      }
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