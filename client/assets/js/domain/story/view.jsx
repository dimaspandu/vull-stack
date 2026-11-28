function StoryLayout({ back }, children) {
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

      {children}
    </div>
  );
}

function StoryDetail({
  category,
  title,
  content,
  publishedAt,
  readingTime,
  back
}) {
  return (
    <StoryLayout back={back}>
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
    </StoryLayout>
  );
}

function StoryAlert({ role, title, description, icon, back }) {
  return (
    <div class="story-alert" role={role}>
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
          {icon}
        </svg>
      </div>

      <div class="story-alert__content">
        <h2 class="story-alert__title">
          {title}
        </h2>

        <p class="story-alert__description">
          {description}
        </p>

        <a
          class="story-alert__link"
          href="/"
          data-link
          onClick={e => {
            e.preventDefault();
            back();
          }}
        >
          ← Back to stories
          <span aria-hidden="true">→</span>
        </a>
      </div>
    </div>
  );
}

function EmptyState({ back }) {
  return (
    <StoryLayout back={back}>
      <StoryAlert
        role="status"
        title="Story not found"
        description="The story you are looking for does not exist or has been moved."
        icon={
          <>
            <path d="M12 11v5" />
            <path d="M12 8h.01" />
          </>
        }
        back={back}
      />
    </StoryLayout>
  );
}

function ErrorState({ back }) {
  return (
    <StoryLayout back={back}>
      <StoryAlert
        role="alert"
        title="Something went wrong"
        description="We couldn't load the story due to a server error. Please try again in a moment."
        icon={
          <>
            <path d="M12 8v4" />
            <path d="M12 16h.01" />
          </>
        }
        back={back}
      />
    </StoryLayout>
  );
}

export default function createStoryView(viewModel, back) {
  const storyView = document.createComponent(() => {
    if (viewModel.notFound) {
      return EmptyState({ back });
    }

    if (viewModel.error) {
      return ErrorState({ back });
    }

    if (!viewModel.story) {
      return EmptyState({ back });
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