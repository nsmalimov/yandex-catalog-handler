package entity

type Result struct {
	Cause   string         `db:"cause"`
	Results []ResultByFile `db:"result"`
}

type ResultByFile struct {
	FileName   string `json:"filename"`
	Was        int    `json:"was"`
	Now        int    `json:"now"`
	WasRemove  int    `json:"was_remove"`
	ErrorCause string `json:"error_cause"`
}
