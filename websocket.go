package zetka

import (
	"compress/zlib"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/valyala/fastjson"

	"github.com/paramusio/go-zetka/intent"
	"github.com/paramusio/go-zetka/message"

	"github.com/paramusio/go-zetka/opcode"

	"github.com/gorilla/websocket"
)

type GatewayEvent struct {
	Type     string          `json:"t"`
	Sequence int64           `json:"s"`
	OpCode   opcode.Opcode   `json:"op"`
	Data     json.RawMessage `json:"d"`
}

func (c *Client) Receive(gwuri string, results chan *GatewayEvent) error {
	dialer := &websocket.Dialer{
		EnableCompression: true,
	}

	conn, _, err := dialer.Dial(gwuri, nil)
	if err != nil {
		return err
	}
	WebSocketConnectionCounter.Inc()

	for {
		select {
		// On restart, let's close the current connection and re-dial.

		// TODO(tobbbles): This is naive and could result in endless recursion.
		// We should add some logic here for maximum retries
		case <-c.restartch:
			_ = conn.Close()
			WebSocketConnectionCounter.Dec()

			return c.Receive(gwuri, results)

		default:
			mt, reader, err := conn.NextReader()
			if err != nil {
				return err
			}

			event, err := c.decode(mt, reader)
			if err != nil {
				return err
			}

			if err := c.parse(conn, event); err != nil {
				return err
			}

			results <- event
		}
	}
}

func (c *Client) decode(mt int, reader io.Reader) (*GatewayEvent, error) {
	msg := &GatewayEvent{}

	switch mt {
	case websocket.TextMessage:
		if err := json.NewDecoder(reader).Decode(msg); err != nil {
			return nil, err
		}
		return msg, nil

	case websocket.BinaryMessage:
		r, err := zlib.NewReader(reader)
		if err != nil {
			return nil, err
		}

		if err := json.NewDecoder(r).Decode(msg); err != nil {
			return nil, err
		}

		return msg, nil

	default:
		return nil, errors.New("there has been a serious programming error")
	}
}

func (c *Client) parse(conn *websocket.Conn, msg *GatewayEvent) error {
	switch msg.OpCode {
	case opcode.Heartbeat:
		// TODO: Start a timer to measure differenence between heartbeat rcv -> send -> ack
		c.sendHeartbeat(conn, c.sequence.Load().(int64))

	case opcode.Hello:
		// Parse Interval
		interval := fastjson.GetInt(msg.Data, "heartbeat_interval")

		// Start heartbeat dispatcher
		go c.heartbeat(conn, time.Duration(interval), c.sequence.Load().(int64))

		// Send authentication packets
		if err := c.auth(conn); err != nil {
			return err
		}

	case opcode.Dispatch:
		c.sequence.Store(msg.Sequence)

	case opcode.InvalidSession:
		// Send a signal to re-start the receive loop
		c.restartch <- struct{}{}

	case opcode.HeartbeatACK:

	default:
		fmt.Printf("Received unknown event code %d\n", msg.OpCode)
		fmt.Println(string(msg.Data))
	}

	return nil
}

// heartbeat starts a go routine to send heart beats at the given interval
// TODO(tobbbles) look at an sync atomic value for the interval in case it will change dynamically
func (c *Client) heartbeat(conn *websocket.Conn, interval time.Duration, sequence int64) {
	for {
		<-time.After(interval * time.Millisecond)

		c.sendHeartbeat(conn, sequence)
	}
}
func (c *Client) sendHeartbeat(conn *websocket.Conn, sequence int64) {
	if err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(`{"op": 1,"d":%d}`, sequence))); err != nil {
		c.errs <- err
	}
}

func (c *Client) auth(conn *websocket.Conn) error {
	// Send auth
	a := &message.Authenticate{
		Token: c.token,
		Properties: message.Properties{
			OS:      "linux",
			Browser: "paramus",
			Device:  "paramus",
		},
		Compress: c.compress,
		Shard:    [2]int{0, 1},
		Intents:  intent.All,
	}

	authbuf, err := json.Marshal(a)
	if err != nil {
		return err
	}

	buf, err := json.Marshal(&GatewayEvent{OpCode: opcode.Identify, Data: authbuf})
	if err != nil {
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, buf)
}
