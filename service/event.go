package service

import (
	"account-api/model"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
)

type EventSender interface {
	Send(key model.FourEyesMessageKey, message any) error
}

type EventReceiver interface {
	Receive() ([]any, error)
}

type KafkaConfig struct {
	Url       string
	Topic     string
	Partition int
}

type KafkaSender struct {
	Conn *kafka.Conn
}

func NewKafkaSender(config KafkaConfig) EventSender {
	conn, err := kafka.DialLeader(context.Background(), "tcp", config.Url, config.Topic, config.Partition)
	if err != nil {
		return nil
	}
	return &KafkaSender{Conn: conn}
}

func (k *KafkaSender) Send(key model.FourEyesMessageKey, message any) error {
	bytesMsg, err := json.Marshal(message)
	if err != nil {
		return err
	}
	bytesKey, err := json.Marshal(key)
	if err != nil {
		return err
	}

	_, err = k.Conn.WriteMessages(kafka.Message{Key: bytesKey, Value: bytesMsg})
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := k.Conn.Close(); err != nil {
		return err
	}
	return nil
}
