package qavideo

type Feed struct {
	Content         []Entry `json:"content"`
	QuestionAnswers []Entry `json:"question_answers"`
	Comments        []Entry `json:"comments"`
}

type Entry struct {
	Url        string `json:"url"`
	Username   string `json:"username,omitempty"`
	Text       string `json:"text"`
	AvatarFile string `json:"avatar_file,omitempty"`
	Role       string `json:"role,omitempty"`
	Datetime   string `json:"datetime"`
	DataID     string `json:"data_id,omitempty"`
	ParentID   string `json:"parent_id,omitempty"`
	Type       string `json:"type"`
	Position   int    `json:"position"`
	Chunk      int    `json:"chunk,omitempty"`
}
