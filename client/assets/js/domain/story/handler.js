export default function showHandler() {
  return {
    onMeet: async (ctx) => {
      console.log("ctx", ctx);
      const { default: view } = await import("../story/view.jsx");
    }
  };
}