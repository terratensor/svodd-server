package qaparser

type Entry struct {
	Url        string `json:"url"`
	Username   string `json:"username"`
	Text       string `json:"text"`
	AvatarFile string `json:"avatar_file"`
	Role       string `json:"role"`
	Datetime   string `json:"datetime"`
	DataID     string `json:"data_id,omitempty"`
	ParentID   string `json:"parent_id"`
	Type       string `json:"type"`
	Position   int    `json:"position"`
	Chunk      int    `json:"chunk,omitempty"`
}