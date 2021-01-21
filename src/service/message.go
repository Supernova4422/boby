package service

// A Message is sent using a Sender.
type Message struct {
	URL         string
	Title       string
	Description string
	Fields      []MessageField
	Footer      string
}

// A MessageField stores a field and value pair.
type MessageField struct {
	Field string
	Value string
	URL   string
}
