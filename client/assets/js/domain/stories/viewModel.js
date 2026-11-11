export default async function createStoriesViewModel(ctx, repository) {
  const listeners = new Set();
  const stories = await repository.getStories(ctx);

  return {
    stories,

    query: "",

    addListener(listener) {
      listeners.add(listener);
      return () => listeners.delete(listener);
    },

    notifyListener() {
      listeners.forEach((listener) => listener(this));
    },
  };
}