package web_stories

type Story struct {
	ID          string   `json:"id"`
	Slug        string   `json:"slug"`
	Title       string   `json:"title"`
	Category    string   `json:"category"`
	Excerpt     string   `json:"excerpt"`
	Content     []string `json:"content"`
	PublishedAt string   `json:"publishedAt"`
	ReadingTime int      `json:"readingTime"`
}

type StoryList struct {
	Stories []Story
}
