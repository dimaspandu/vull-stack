export default function createStoryViewModel(repository) {
  const listeners = new Set();

  const viewModel = {
    story: null,
    error: null,

    // Called by the controller on every route change. The underscore marks it
    // as internal: the view only subscribes through addListener().
    async _load(ctx) {
      this.error = null;

      try {
        this.story = await repository.getStory(ctx);
      } catch (err) {
        this.error = err.message;
        this.story = null;
      } finally {
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