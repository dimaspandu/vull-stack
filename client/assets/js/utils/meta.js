// Keeps the document title and the meta description in sync during client-side
// navigation, so a route change never leaves the previous page's metadata.
export function setPageMeta(title, description) {
  document.title = title;

  const metaDescription = document.querySelector('meta[name="description"]');
  if (metaDescription) {
    metaDescription.content = description;
  }
}