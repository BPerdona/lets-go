package message

// Message represents a message received from a peer
type Message struct {
	From    string
	Payload []byte
}
