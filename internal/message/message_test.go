package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	msg := Message{
		From:    "127.0.0.1:8080",
		Payload: []byte("Hello, world!"),
	}

	assert.Equal(t, "127.0.0.1:8080", msg.From)
	assert.Equal(t, []byte("Hello, world!"), msg.Payload)
}
