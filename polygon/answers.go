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

type PackageAnswer struct {
	ID                  int    `json:"id"`
	Revision            int    `json:"revision"`
	CreationTimeSeconds int    `json:"creationTimeSeconds"`
	State               string `json:"state"`
	Comment             string `json:"comment"`
	Type                string `json:"type"`
}

type ProblemAnswer struct {
	ID            int    `json:"id"`
	Owner         string `json:"owner"`
	Name          string `json:"name"`
	Deleted       bool   `json:"deleted"`
	Favourite     bool   `json:"favourite"`
	AccessType    string `json:"accessType"`
	Revision      int    `json:"revision"`
	LatestPackage int    `json:"latestPackage,omitempty"`
	Modified      bool   `json:"modified"`
}

type Answer struct {
	Status  string          `json:"status"`
	Comment string          `json:"comment"`
	Result  json.RawMessage `json:"result"`
}
