// Router-agnostic state for a single story. The controller supplies the story
// data; this only fans out changes to subscribers.
export default function createStoryViewModel(repository) {
  const listeners = new Set();

  const viewModel = {
    story: null,
    error: null,
    notFound: false,

    // Seed state from an already-resolved story (for example, the
    // server-rendered markup on the initial request).
    _hydrate(story) {
      this.error = null;
      this.notFound = false;
      this.story = story;
      this.notifyListener();
    },

    // Reload state through the repository on every route change. The controller
    // passes a plain slug instead of the router context. A 404 means the story
    // does not exist (not-found state), anything else is a server error.
    async _load(slug) {
      this.error = null;
      this.notFound = false;

      try {
        this.story = await repository.getStory(slug);
      } catch (err) {
        if (err.status === 404) {
          this.notFound = true;
          this.story = null;
        } else {
          this.error = err.message;
          this.story = null;
        }
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