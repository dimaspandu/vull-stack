/** @jsx elementBuilder */
function StoryListItem({
  slug,
  category,
  title,
  excerpt,
  publishedAt,
  readingTime
}) {
  return (
    <li class="story-list__item">
      <article class="story-card">
        <a
          class="story-card__link"
          href={`/${slug}`}
          data-link
          onClick={e => {
            e.preventDefault();
            history.navigatePush(`/${slug}`);
          }}
        >
          <div class="story-card__content">
            <p class="story-card__category">
              {category}
            </p>

            <h3 class="story-card__title">
              {title}
            </h3>

            <p class="story-card__excerpt">
              {excerpt}
            </p>

            <div class="story-card__meta">
              <time datetime="2026-10-06">
                {publishedAt}
              </time>

              <span aria-hidden="true">·</span>

              <span>
                {readingTime} min read
              </span>
            </div>
          </div>
        </a>
      </article>
    </li>
  );
}

export default function view({
  query,
  stories,
}) {
  return (
    <div class="story-page__container">

      <header class="story-page__header">
        <h1 class="story-page__title">
          Short Stories
        </h1>

        <p class="story-page__description">
          Small stories for a few quiet minutes.
        </p>
      </header>

      <form
        class="story-search"
        id="story-search"
        action="/"
        method="get"
        role="search"
        onSubmit={(e) => {
          e.preventDefault();
        }}
      >
        <div class="story-search__field">
          <input
            class="story-search__input"
            name="q"
            type="search"
            placeholder="Search stories..."
            autocomplete="off"
            value={query}
          />

          <button
            class="story-search__button"
            type="submit"
          >
            Search
          </button>
        </div>
      </form>

      <section
        class="story-list"
        aria-labelledby="story-list-title"
      >
        <header class="story-list__header">
          <h2
            class="story-list__title"
            id="story-list-title"
          >
            Latest stories
          </h2>

          <p class="story-list__count">
            {stories.length} stories
          </p>
        </header>

        {stories.length > 0 ? (
          <ul class="story-list__items" id="story-list-items">
            {stories.map(story => <StoryListItem { ...story } />)}
          </ul>
        ) : (
          <div class="story-alert" role="status">
            <div class="story-alert__icon" aria-hidden="true">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.7"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <circle cx="12" cy="12" r="9" />
                <path d="M12 11v5" />
                <path d="M12 8h.01" />
              </svg>
            </div>

            <div class="story-alert__content">
              <h3 class="story-alert__title">
                No stories found
              </h3>

              <p class="story-alert__description">
                We could not find any stories matching your search.
                Try a different keyword.
              </p>

              <a class="story-alert__link" href="/" data-link>
                Browse all stories
                <span aria-hidden="true">→</span>
              </a>
            </div>
          </div>
        )}

      </section>

    </div>
  );
}