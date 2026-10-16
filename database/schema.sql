CREATE TABLE stories (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  slug VARCHAR(191) NOT NULL,
  title VARCHAR(255) NOT NULL,
  category VARCHAR(50) NOT NULL,
  excerpt TEXT NOT NULL,
  content JSON NOT NULL,
  published_at DATETIME NULL DEFAULT NULL,
  reading_time TINYINT UNSIGNED NOT NULL DEFAULT 1,

  PRIMARY KEY (id),
  UNIQUE KEY uq_stories_slug (slug),
  CONSTRAINT chk_stories_reading_time
    CHECK (reading_time > 0)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;