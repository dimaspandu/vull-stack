/** @jsx elementBuilder */
export default function createNotFoundView() {
  return (
    <div class="story-alert" role="alert">

      <div class="story-alert__content">
        <h1 class="story-alert__title">
          404 — Page not found
        </h1>

        <p class="story-alert__description">
          The page you are looking for does not exist or has been moved.
        </p>

        <a
          class="story-alert__link"
          href="/"
          data-link
          onClick={e => {
            e.preventDefault();
            history.navigatePush("/");
          }}
        >
          ← Back to stories
        </a>
      </div>

    </div>
  );
}