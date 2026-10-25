import { bindLinkNavigation } from "../../utils/common.js";

export default function showHandler() {
  return {
    onMeet: async (ctx) => {
      if (!ctx.isInitialRoute) {
        const { default: view } = await import("../story/view.jsx");

        try {
          const response = await fetch(`/api/stories/${ctx.params.slug}`);
          if (!response.ok) {
            throw new Error(`Response status: ${response.status}`);
          }
          const result = await response.json();

          outlet
            .setAttributes({ class: "story-detail-page" })
            .setChildren(
              view(result.data)
            );
        } catch (error) {
          console.error(error.message);
        }
      } else {}

      bindLinkNavigation();
    }
  };
}