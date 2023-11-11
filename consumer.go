package main

import (
	"github.com/IBM/sarama"
	"github.com/kr/pretty"
	"github.com/nxczje/froxy/binary"
	"github.com/nxczje/froxy/parser"
)

func main() {
	target := "nothinnn.oob.nncg.uk"   //change this
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
	pretty.Println("Waiting for messages....", target)
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			req, err := parser.ParseRequest(string(msg.Value))
			if err != nil {
				pretty.Println(err)
			}
			pretty.Println(req.Host)
			if target == req.Host {
				//work somthing
				binary.X8(req.URL.Scheme + "://" + req.URL.Host + req.URL.Path)
			}
			// pretty.Println("Received messages", string(msg.Value))

		case err := <-partitionConsumer.Errors():
			pretty.Println("Received errors", err)
		}
	}

}
