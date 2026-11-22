import createLoadingView from "./view.jsx";
import createLoadingViewModel from "./viewModel.js";

const viewModel = createLoadingViewModel();

// Mounted once into #portal, then driven entirely by the viewModel so any
// controller can toggle it without touching the DOM.
let mounted = false;

// Overlapping navigations must not hide the indicator early, so show/hide are
// reference counted instead of a plain boolean.
let pending = 0;

function mount() {
  if (mounted) return;

  const portal = document.getElementById("portal");
  if (!portal) return;

  mounted = true;

  document.createElement(portal).setChildren(createLoadingView(viewModel));
}

export function showLoading() {
  pending++;
  viewModel.show();
  mount();
}

export function hideLoading() {
  pending = Math.max(0, pending - 1);
  if (pending === 0) {
    viewModel.hide();
  }
}

// Route transition indicator, owned here instead of by every controller. The
// server already rendered the first route, so that load is not a transition and
// must not flash the indicator.
export function withLoading(handler) {
  return async (ctx) => {
    if (ctx.firstArrive && ctx.isInitialRoute) {
      return handler(ctx);
    }

    showLoading();

    try {
      return await handler(ctx);
    } finally {
      hideLoading();
    }
  };
}