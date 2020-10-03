package message

import "github.com/paramusio/go-zetka/intent"

type Shard [2]int

type Authenticate struct {
	Token      string        `json:"token"`
	Properties Properties    `json:"properties"`
	Compress   bool          `json:"compress"`
	Shard      Shard         `json:"shard"`
	Intents    intent.Intent `json:"intents"`
}

type Properties struct {
	OS      string `json:"$os"`
	Browser string `json:"$browser"`
	Device  string `json:"$device"`
}
