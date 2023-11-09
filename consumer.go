package main

import (
	"github.com/IBM/sarama"
	"github.com/kr/pretty"
)

func main() {
	Addr := []string{"127.0.0.1:9092"} //change this
	Topic := "nothing"                 //change this
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumer, err := sarama.NewConsumer(Addr, config)
	if err != nil {
		pretty.Println(err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			panic(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition(Topic, 0, sarama.OffsetOldest)
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
