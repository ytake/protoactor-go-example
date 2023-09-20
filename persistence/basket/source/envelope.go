package source

import (
	"encoding/json"
	"log"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type envelope struct {
	Type       string          `json:"type"`
	Message    json.RawMessage `json:"event"`
	EventIndex int             `json:"eventIndex"`
	DocType    string          `json:"doctype"`
}

func newEnvelope(message proto.Message, doctype string, eventIndex int) *envelope {
	typeName := proto.MessageName(message)
	bytes, err := json.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}
	envelope := &envelope{
		Type:       string(typeName),
		Message:    bytes,
		EventIndex: eventIndex,
		DocType:    doctype,
	}
	return envelope
}

func (envelope *envelope) message() proto.Message {
	mt, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(envelope.Type))
	if err != nil {
		log.Fatal(err)
	}

	pm := mt.New().Interface()
	err = json.Unmarshal(envelope.Message, pm)
	if err != nil {
		log.Fatal(err)
	}
	return pm
}
