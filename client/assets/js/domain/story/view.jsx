/** @jsx elementBuilder */
function StoryDetail({
  category,
  title,
  content,
  publishedAt,
  readingTime
}) {
  return (
    <div class="story-detail-page__container">

      <nav
        class="story-detail-page__navigation"
        aria-label="Story navigation"
      >
        <a
          class="story-detail-page__back"
          href="/"
          data-link
          onClick={e => {
            e.preventDefault();
            history.back();
          }}
        >
          ← Back to stories
        </a>
      </nav>

      <article class="story">

        <header class="story__header">

          <p class="story__category">
            {category}
          </p>

          <h1 class="story__title">
            {title}
          </h1>

          <div class="story__meta">
            <time datetime="2026-10-06">
              {publishedAt}
            </time>

            <span aria-hidden="true">·</span>

            <span>
              {readingTime} min read
            </span>
          </div>

        </header>

        <div class="story__body">

          {content.map(text => (<p>{text}</p>))}

        </div>

      </article>

    </div>
  );
}

function ErrorState() {
  return elementBuilder("h1", null, "Internal Server Error.");
}

function EmptyState() {
  return elementBuilder("h1", null, "Story not found.");
}

export default function createStoryView(viewModel) {
  const storyView = componentBuilder(() => {
    if (viewModel.error) {
      return ErrorState();
    }

    if (!viewModel.story) {
      return EmptyState();
    }

    const {
      category,
      title,
      content,
      publishedAt,
      readingTime
    } = viewModel.story;

    return (
      <StoryDetail
        category={category}
        title={title}
        content={content}
        publishedAt={publishedAt}
        readingTime={readingTime}
      />
    );
  });

  viewModel.addListener(() => storyView.rerender());

  return storyView.render();
}