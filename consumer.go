package main

import (
	"fmt"
	"log"
	"os"

	"github.com/IBM/sarama"
	"github.com/kr/pretty"
	"github.com/nxczje/froxy/parser"
)

type Node struct {
	target string
	data   string
	left   *Node
	right  *Node
}

type BinaryTree struct {
	root *Node
}

func (t *BinaryTree) insert(target string, data string) *BinaryTree {
	if t.root == nil {
		if _, err := os.Stat("output/" + target); os.IsNotExist(err) {
			file, err := os.Create("output/" + target)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			t.root = &Node{target: target, data: data, left: nil, right: nil}

			defer file.Close()
		} else {
			t.root = &Node{target: target, data: data, left: nil, right: nil}
		}
		err := os.WriteFile("output/"+target, []byte(data), 0644)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		t.root.insert(target, data)
	}
	return t
}

func (n *Node) insert(target string, data string) {
	if n == nil {
		return
	} else if data < n.data {
		if n.left == nil {
			n.left = &Node{target: target, data: data, left: nil, right: nil}
		} else {
			n.left.insert(target, data)
		}
	} else if data > n.data {
		if n.right == nil {
			n.right = &Node{data: data, left: nil, right: nil}
		} else {
			n.right.insert(target, data)
		}
	}
	// do nothing if data == n.data
	if data == n.data {
		return
	}
}

func print(node *Node) {
	if node == nil {
		return
	}
	print(node.left)
	fmt.Println(node.data)
	print(node.right)
}

func main() {
	target := "nothinnn.oob.nncg.uk"   //change this
	Addr := []string{"127.0.0.1:9092"} //change this
	Topic := "nothing"                 //change this
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumer(Addr, config)
	if err != nil {
		log.Fatalln(err)
	}

	partitionConsumer, err := consumer.ConsumePartition(Topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatalln(err)
	}

	t := &BinaryTree{}
	pretty.Println("Waiting for messages....", target)
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			req, err := parser.ParseRequest(string(msg.Value))
			if err != nil {
				pretty.Println(err)
			}
			url := req.URL.Host
			pretty.Println(req.URL)
			t.insert(target, url)
			// if target == req.Host {
			//work somthing
			// binary.X8(req.URL.Scheme + "://" + req.URL.Host + req.URL.Path)
			// }

		case err := <-partitionConsumer.Errors():
			pretty.Println("Received errors", err)
		}
		print(t.root)
	}
}
