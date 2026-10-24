export default function listHandler(outlet) {
  return {
    onMeet: async (ctx) => {
      if (!ctx.isInitialRoute) {
        const { default: view } = await import("../stories/view.jsx");

        try {
          const response = await fetch("/api/stories");
          if (!response.ok) {
            throw new Error(`Response status: ${response.status}`);
          }

          const result = await response.json();
          console.log(result);
          outlet
            .setAttributes({ class: "story-page" })
            .setChildren(
              view({
                query: "",
                stories: result.data,
              })
            );
        } catch (error) {
          console.error(error.message);
        }
      } else {}
    },
  };
  // return {
      // const { default: view } = await import("./view.jsx");
      // const stories = JSON.parse(
      //   document.getElementById("stories-data").textContent
      // );
      // const {
      //   renderStoryListItem,
      // } = view();
      // const storySearchView = document.getElementById("story-search");
      // const storyListItemsView = document.getElementById("story-list-items");

      // storySearchView.setEvents({
      //   submit(e) {
      //     e.preventDefault();

      //     const formData = new FormData(storySearchView);
      //     const query = formData.get("q");
      //     const filteredStories = stories.filter((story) =>
      //       ["slug", "title", "category", "excerpt"].some((key) =>
      //         story[key]?.toLowerCase().includes(query.toLowerCase())
      //       )
      //     );
      //     const storyLisItemView = filteredStories.map(
      //       story => renderStoryListItem(story)
      //     );

      //     storyListItemsView.setChildren(storyLisItemView);
      //   }
      // });
      // const storyLisItemView = stories.map(
      //   story => renderStoryListItem(story)
      // );
      // console.log("storyLisItemView", storyLisItemView);
      // console.log("storyListItemsView", storyListItemsView);

      // storyListItemsView
      // // // .setState({ stories })
      // .setChildren(() => 
      //   stories.map(
      //     story => renderStoryListItem(story)
      //   )
      // );
    // },
  // };
}

async function handleOnMeet(outlet) {
  if (!ctx.isInitialRoute) {
    const { default: view } = await import("./view.jsx");

    try {
      const response = await fetch("/api/stories");
      if (!response.ok) {
        throw new Error(`Response status: ${response.status}`);
      }

      const result = await response.json();
      console.log(result);
      outlet
        .setAttributes({ class: "story-page" })
        .setChildren(
          view({
            query: "",
            stories: result.data,
          })
        );
    } catch (error) {
      console.error(error.message);
    }
  } else {}
}