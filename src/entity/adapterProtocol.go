package entity

type TargetProtocol struct {
	ID     int         `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params,omitempty"`
}

type TargetParams struct {
	ID       int           `json:"id,omitempty"`
	Message  interface{}   `json:"message,omitempty"`
	TargetId string        `json:"targetId,omitempty"`
	Edits    []interface{} `json:"edits,omitempty"`
}

type IRange struct {
	StartLine   int
	StartColumn int
	EndLine     int
	EndColumn   int
}

type IDisabledStyle struct {
	Content  string
	CssRange IRange
}
