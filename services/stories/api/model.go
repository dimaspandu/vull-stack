package api_stories

type Story struct {
	ID          int      `json:"id"`
	Slug        string   `json:"slug"`
	Title       string   `json:"title"`
	Category    string   `json:"category"`
	Excerpt     string   `json:"excerpt"`
	Content     []string `json:"content"`
	PublishedAt string   `json:"publishedAt"`
	ReadingTime int      `json:"readingTime"`
}

type StoryListResponse struct {
	Data []Story `json:"data"`
}

type StoryShowResponse struct {
	Data *Story `json:"data"`
}
