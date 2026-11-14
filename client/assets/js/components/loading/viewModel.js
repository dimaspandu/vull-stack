export default function createLoadingViewModel() {
  const listeners = new Set();

  const viewModel = {
    visible: false,

    show() {
      if (this.visible) return;

      this.visible = true;
      this.notifyListener();
    },

    hide() {
      if (!this.visible) return;

      this.visible = false;
      this.notifyListener();
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