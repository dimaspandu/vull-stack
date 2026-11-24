// Mounting a view into the outlet is a router concern, so it lives here
// instead of in every controller. The controller decides what to render and
// hands back a view; this decides how that view replaces or hydrates the
// outlet.
//
// Two distinct behaviors are needed:
//
//   hydrate — the server already rendered this markup, so adopt the existing
//             children instead of replacing them. Uses the props form of
//             createElement, which adopts the outlet node and attaches the
//             new children on top of the server ones.
//
//   replace — a client-side navigation already swapped the route, so the
//             previous outlet children must be dropped before the new view is
//             mounted. setChildren() reconciles and removes the stale nodes,
//             which the props form of createElement never does.
export function mountOutlet(view, { hydrate = false } = {}) {
  const outlet = document.getElementById("outlet");
  if (!outlet) {
    return null;
  }

  if (hydrate) {
    document.createElement(outlet, { class: "story-page" }, view);
    return outlet;
  }

  document.createElement(outlet)
    .setAttributes({ class: "story-page" })
    .setChildren(view);

  return outlet;
}