package polygon

import "encoding/json"

type TestAnswer struct {
	Index           int     `json:"index"`
	Group           string  `json:"group"`
	Points          float32 `json:"points"`
	UseInStatements bool    `json:"useInStatements"`
}

type GroupAnswer struct {
	Name           string   `json:"name"`
	PointsPolicy   string   `json:"pointsPolicy"`
	FeedbackPolicy string   `json:"feedbackPolicy"`
	Dependencies   []string `json:"dependencies"`
}

type Answer struct {
	Status  string          `json:"status"`
	Comment string          `json:"comment"`
	Result  json.RawMessage `json:"result"`
}
