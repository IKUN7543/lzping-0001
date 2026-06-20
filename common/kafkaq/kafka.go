package kafkaq

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"go-zero-ecommerce/common/errx"
	"time"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireAll,
			Async:        false,
		},
	}
}

func (p *Producer) Send(ctx context.Context, topic string, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: data,
	}
	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		return errx.ErrKafkaSendFail
	}
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic string, groupId string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        brokers,
			Topic:          topic,
			GroupID:        groupId,
			MinBytes:       10e3,
			MaxBytes:       10e6,
			CommitInterval: time.Second,
		}),
	}
}

func (c *Consumer) Consume(ctx context.Context, handler func(ctx context.Context, key string, value []byte) error) error {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			return err
		}
		err = handler(ctx, string(msg.Key), msg.Value)
		if err == nil {
			c.reader.CommitMessages(ctx, msg)
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
