export default async function createStoriesViewModel(ctx, repository) {
  const listeners = new Set();
  const stories = await repository.getStories(ctx);

  return {
    stories,

    query: "",

    fetchStories() {},

    search(query) {
      this.query = query;
      // this.stories.filter(story => {
      //   if (
      //     story.slug === query ||
      //     story.title === query ||
      //     story.excerpt === query
      //   ) {
      //     return story;
      //   }
      // });
    },

    addListener(listener) {
      listeners.add(listener);
      return () => listeners.delete(listener);
    },

    notifyListener() {
      listeners.forEach((listener) => listener(this));
    },
  };
}