package zetka

import "github.com/paramus/go-zetka/message"

type Message struct {
	ID message.Snowflake `json:"id"`
}
