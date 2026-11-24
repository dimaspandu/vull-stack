function StoryDetail({
  category,
  title,
  content,
  publishedAt,
  readingTime,
  back
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
            back();
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
            <time datetime={publishedAt}>
              {publishedAt}
            </time>

            <span aria-hidden="true">·</span>

            <span>
              {readingTime}{` `}min read
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
  return <h1>Internal Server Error.</h1>;
}

function EmptyState() {
  return <h1>Story not found.</h1>;
}

export default function createStoryView(viewModel, back) {
  const storyView = document.createComponent(() => {
    if (viewModel.error) {
      return ErrorState();
    }

    if (!viewModel.story) {
      return EmptyState();
    }

    return (
      <StoryDetail
        category={viewModel.story.category}
        title={viewModel.story.title}
        content={viewModel.story.content}
        publishedAt={viewModel.story.publishedAt}
        readingTime={viewModel.story.readingTime}
        back={back}
      />
    );
  });

  viewModel.addListener(() => storyView.rerender());

  return storyView;
}