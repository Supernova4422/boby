package service

import (
	"fmt"
	"image"
)

// A Message is sent using a Sender.
type Message struct {
	URL         string
	Title       string
	Description string
	Fields      []MessageField
	Image       image.Image
}

// A MessageField stores a field and value pair.
type MessageField struct {
	Field  string
	Value  string
	URL    string
	Inline bool
}

func (m *MessageField) toString() string {
	return fmt.Sprintf("%s,%s,%s,%v", m.Field, m.Value, m.URL, m.Inline)
}

func (m *Message) ToString() string {
	fieldsString := ""
	for i, field := range m.Fields {
		if i > 0 {
			fieldsString += ","
		}
		fieldsString += field.toString()
	}

	return fmt.Sprintf("%s,%s,%v,{%s}", m.Title, m.Description, m.Image == nil, fieldsString)
}
