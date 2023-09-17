package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
	"testFIO/cmd/config"
	"testFIO/cmd/loggers"
)

type Producer struct {
	conn   *kafka.Conn
	cfg    config.ServerConfig
	Msg    []kafka.Message
	Batch  kafka.Batch
	logger loggers.Logger
}

func (p *Producer) WriteMessages(message ...kafka.Message) error {
	if _, err := p.conn.WriteMessages(message...); err != nil {
		p.logger.LogErr(err, "failed to write messages")
		return err
	}
	return nil
}

func NewProducer(cfg config.ServerConfig, logger loggers.Logger) (*Producer, error) {
	var msg []kafka.Message
	var b kafka.Batch
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", cfg.Kafka.TopicProducer, 0)
	if err != nil {
		logger.LogErr(err, "failed to dial leader")
		return nil, err
	}
	return &Producer{
		conn:  conn,
		cfg:   cfg,
		Msg:   msg,
		Batch: b,
	}, nil
}

func (p *Producer) NewMessage(b []byte) kafka.Message {
	msg := kafka.Message{
		Topic: p.cfg.Kafka.TopicProducer,
		Value: b,
	}
	return msg
}
