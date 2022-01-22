package lib

import "time"

type PushRequest struct {
	Streams []PushStream `json:"streams"`
}

type PushStream struct {
	Stream Label   `json:"stream"`
	Values []Value `json:"values"`
}

type Label map[string]string

type Value [2]string

type LabelInfos map[string]Label

type LogInfos map[string]Logs

type Logs []Log

type Log struct {
	Time time.Time
	Log  string
}
