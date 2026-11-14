export function normalizePath(path) {
  if (path !== "/" && path.endsWith("/")) {
    return path.slice(0, -1);
  }
  return path;
}

export function bindLinkNavigation() {
  document.addEventListener("click", (e) => {
    const a = e.target.closest("[data-link]");
    if (!a) return;

    const url = new URL(a.href);
    if (url.origin !== location.origin) return;

    e.preventDefault();

    const path = normalizePath(url.pathname);
    history.navigatePush(path);
  });
}