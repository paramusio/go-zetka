package zetka

import "github.com/paramusio/go-zetka/message"

type Message struct {
	ID message.Snowflake `json:"id"`
}
