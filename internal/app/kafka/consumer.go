package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"testFIO/cmd/config"
	"testFIO/cmd/loggers"
	"testFIO/internal/app/api"
	"testFIO/internal/storage"
	"testFIO/internal/storage/model"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	service  storage.Service
	external api.External
	reader   *kafka.Reader
	logger   loggers.Logger
	producer *Producer
}

func NewConsumer(service storage.Service, external api.External, cfg config.ServerConfig, logger loggers.Logger) (*Consumer, error) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{cfg.Kafka.Addr},
		GroupID:        "servers",
		Topic:          cfg.Kafka.TopicConsumer,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
	})
	producer, err := NewProducer(cfg, logger)
	if err != nil {
		logger.LogErr(err, "")
		return nil, err
	}
	return &Consumer{
		service:  service,
		external: external,
		reader:   r,
		logger:   logger,
		producer: producer,
	}, nil
}

func (c *Consumer) Read() error {
	for {
		m, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			c.logger.LogErr(err, "failed to read messages")
			break
		}
		var u model.User
		if err = json.Unmarshal(m.Value, &u); err != nil {
			c.logger.LogErr(err, "failed to unmarshal messages")
			break
		}
		if u.Name == "" || u.Surname == "" {
			u.Error = fmt.Errorf("first or last name is not filled in")
			b, err := json.Marshal(u)
			if err != nil {
				c.logger.LogErr(err, "failed to marshal messages")
				break
			}
			if err := c.producer.WriteMessages(c.producer.NewMessage(b)); err != nil {
				c.logger.LogErr(err, "failed to send messages")
				break
			}
			c.logger.LogInfo("wrong name", "error", "first or last name is not filled in")
		} else {
			if u, err = c.external.AddAge(u); err != nil {
				c.logger.LogErr(err, "failed to get age")
				break
			}
			if u, err = c.external.AddGender(u); err != nil {
				c.logger.LogErr(err, "failed to get gender")
				break
			}
			if u, err = c.external.AddNational(u); err != nil {
				c.logger.LogErr(err, "failed to get national")
				break
			}
			if err = c.service.Collect(u); err != nil {
				c.logger.LogErr(err, "failed to collect user")
				break
			}
		}
	}
	if err := c.reader.Close(); err != nil {
		c.logger.LogErr(err, "failed to close reader")
		return err
	}
	return nil
}
