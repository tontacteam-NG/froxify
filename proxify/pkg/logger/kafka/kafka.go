package kafka

import (
	"context"
	"crypto/md5"
	"fmt"
	"regexp"
	"strings"

	"github.com/IBM/sarama"
	"github.com/go-redis/redis/v8"
	"github.com/nxczje/froxy/proxify/pkg/types"
)

// Options required for kafka
type Options struct {
	// Address for kafka instance
	Addr string `yaml:"addr"`
	// Topic to produce messages to
	Topic string `yaml:"topic"`
	// Redis client
	Redis *redis.Client `yaml:"redis"`
	// TLS passes tls configuration to elasticsearch
	Filter []string `yaml:"tls"`
}

// Client for Kafka
type Client struct {
	producer sarama.SyncProducer
	topic    string
	addr     string
	Redis    *redis.Client
	Filter   []string
}

// New creates and returns a new client for kafka
func New(option *Options) (*Client, error) {

	config := sarama.NewConfig()
	// Wait for all in-sync replicas to ack the message
	config.Producer.RequiredAcks = sarama.WaitForAll
	// Retry up to 10 times to produce the message
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{option.Addr}, config)
	if err != nil {
		return nil, err
	}

	return &Client{
		producer: producer,
		topic:    option.Topic,
		addr:     option.Addr,
		Redis:    option.Redis,
		Filter:   option.Filter,
	}, nil
}

// Store passes the message to kafka
func (c *Client) Save(data types.OutputData) error {

	method := strings.Split(data.DataString, " ")[0]
	if method == "CONNECT" {
		return nil
	}
	// fmt.Println(data.Userdata.Host)
	//process filter before saving
	for _, line := range c.Filter {
		matched, err := regexp.MatchString(line, data.Userdata.Host)
		if err != nil {
			fmt.Println("regexp.MatchString ERROR:", err)
		}
		if matched {
			return nil
		}
	}

	hash := CaculatorHash(data)
	exits, err := c.Redis.SIsMember(context.Background(), "hash_consumer", hash).Result()
	if err != nil {
		fmt.Println("Redis.SIsMember ERROR:", err)
		return err
	} else {
		if exits {
			return nil
		} else {
			c.Redis.SAdd(context.Background(), "hash_consumer", hash)
		}
	}
	msg := &sarama.ProducerMessage{
		Topic: c.topic,
		Value: sarama.StringEncoder(data.DataString),
	}

	_, _, err = c.producer.SendMessage(msg)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func CaculatorHash(data types.OutputData) []byte {
	hasher := md5.New()
	hasher.Write(data.Data)
	md5Hash := hasher.Sum(nil)
	fmt.Printf("MD5 Hash: %x\n", md5Hash)
	return md5Hash
}
