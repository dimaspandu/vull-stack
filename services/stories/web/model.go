package web_stories

type Story struct {
	ID          string `json:"id"`
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Category    string `json:"category"`
	Excerpt     string `json:"excerpt"`
	PublishedAt string `json:"publishedAt"`
	ReadingTime int    `json:"readingTime"`
}

type StoryListResponse struct {
	Data []Story `json:"data"`
}

type StoryList struct {
	Stories []Story `json:"stories"`
}

type StoryShowResponse struct {
	Data StoryShow `json:"data"`
}

type StoryShow struct {
	Story
	Content []string `json:"content"`
}
