package handler

import (
	"github.com/IBM/sarama"
	"github.com/kr/pretty"
)

func HandlerConsumer(addKafka string, TopicKafka string) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumer, err := sarama.NewConsumer([]string{addKafka}, config)
	if err != nil {
		pretty.Println(err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			panic(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition(TopicKafka, 0, sarama.OffsetOldest)
	if err != nil {
		pretty.Println(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			panic(err)
		}
	}()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			pretty.Println("Received messages", string(msg.Value))
		case err := <-partitionConsumer.Errors():
			pretty.Println("Received errors", err)
		}
	}

}
