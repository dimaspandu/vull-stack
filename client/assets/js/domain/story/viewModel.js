export default function createStoryViewModel(repository) {
  const listeners = new Set();

  const viewModel = {
    story: null,
    loading: false,
    error: null,

    async _load(ctx) {
      this.loading = true;
      this.error = null;
      this.notifyListener();

      try {
        this.story = await repository.getStory(ctx);
      } catch (err) {
        this.error = err.message;
        this.story = null;
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