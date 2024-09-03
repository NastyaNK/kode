package models

type SpellCheckResult struct {
	Code int      `json:"code"`
	Pos  int      `json:"pos"`
	Row  int      `json:"row"`
	Col  int      `json:"col"`
	Len  int      `json:"len"`
	Word string   `json:"word"`
	S    []string `json:"s"`
}

type Position struct {
	Pos int `json:"pos"`
	Row int `json:"row"`
	Col int `json:"col"`
}

type SpellResultTO struct {
	Word        string   `json:"word"`
	Suggestions []string `json:"suggestions"`
	Position    Position `json:"position"`
}
