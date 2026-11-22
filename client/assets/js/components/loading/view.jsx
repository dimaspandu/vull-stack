export default function createLoadingView(viewModel) {
  const loadingView = document.createComponent(() => (
    <div
      class="app-loading"
      role="status"
      aria-live="polite"
      hidden={!viewModel.visible}
    >
      <span class="app-loading__spinner" aria-hidden="true"></span>

      <span class="app-loading__label">
        Loading...
      </span>
    </div>
  ));

  viewModel.addListener(() => loadingView.rerender());

  return loadingView;
}