package model

type Formula struct {
	Formula       string                 `json:"formula"`
	Description   string                 `json:"description"`
	Concepts      map[string]interface{} `json:"concepts"`
	SourceChapter string                 `json:"source_chapter"`
}
