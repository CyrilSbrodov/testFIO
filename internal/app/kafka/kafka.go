package kafka

import (
	"github.com/segmentio/kafka-go"
	"testFIO/internal/storage/model"
)

type Kafka interface {
	WriteMessages(message ...kafka.Message) error
	NewMessage(b []byte) kafka.Message
}

type Consume interface {
	SendToProducer(u model.User) error
	SendToBD(u model.User) error
}
