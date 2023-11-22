package main

import (
	"log"
	"slices"
	"sync"

	"github.com/IBM/sarama"
	"github.com/fatih/color"
	"github.com/kr/pretty"
	"github.com/nxczje/froxy/core"
	"github.com/nxczje/froxy/parser"
)

func main() {
	target := "nothinnn.oob.nncg.uk"       //change this
	Addr := []string{"10.14.140.245:9092"} //change this
	Topic := "nothing"                     //change this
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	mu := sync.Mutex{}
	trees := make(map[string]core.Tree, 0)
	consumer, err := sarama.NewConsumer(Addr, config)
	if err != nil {
		log.Fatalln(err)
	}

	partitionConsumer, err := consumer.ConsumePartition(Topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalln(err)
	}
	go core.Run()
	//  := &Tree{Root: &Node{URL: string(target)}}
	color.Green("[INFO] %s", target)
	for {
		select {
		case msg := <-partitionConsumer.Messages():

			req, err := parser.ParseRequest(string(msg.Value))
			if err != nil {
				// pretty.Println(err)
				continue
			}

			mu.Lock()
			tree, exists := trees[req.Host]
			if exists {
				n := tree.Root.SearchNode(req.URL.Path)
				if n != nil && slices.Contains(n.Param, req.URL.RawQuery) {
					// tree.Display(tree.Root, 0)
					mu.Unlock()
					continue

				}
				mu.Unlock()
				tree.AddPath(req.URL.Path, req.URL.RawQuery)
				if err := core.SaveRequest(msg.Value); err != nil {
					log.Println(err.Error())
				}
			} else {
				tree = core.Tree{Mu: &sync.Mutex{}, Root: &core.Node{URL: req.Host}}
				mu.Unlock()
				tree.AddPath(req.URL.Path, req.URL.RawQuery)
				if err := core.SaveRequest(msg.Value); err != nil {
					log.Println(err.Error())
				}
				trees[req.Host] = tree
			}
		case err := <-partitionConsumer.Errors():
			pretty.Println("Received errors", err)
		}
		// t.Display(t.Root, 0)
	}
}
