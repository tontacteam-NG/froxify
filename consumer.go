package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/IBM/sarama"
	"github.com/kr/pretty"
	"github.com/nxczje/froxy/parser"
)

// Record chứa thông tin của một bản ghi
type Record struct {
	URL   string
	Param string
}

// Node là một nút trong cây nhị phân
type Node struct {
	URL      string
	Records  []Record
	Children map[string]*Node
}

// Tree là cấu trúc cây nhị phân
type Tree struct {
	Root *Node
}

// Thêm đường dẫn vào cây nhị phân
func (t *Tree) AddPath(url string, param string) {
	components := strings.Split(url, "/")
	current := t.Root

	for _, component := range components {
		if current.Children == nil {
			current.Children = make(map[string]*Node)
		}
		fmt.Println(current.Children)
		// Kiểm tra xem nút con có tồn tại chưa
		child, exists := current.Children[component]
		if !exists {
			// Nếu không tồn tại, thêm một nút mới
			child = &Node{URL: component}
			current.Children[component] = child
		}

		// Di chuyển xuống nút con
		current = child
	}

	// Lưu trữ thông tin tham số tại nút lá
	current.Records = append(current.Records, Record{URL: url, Param: param})
}

// Hiển thị cây nhị phân
func (t *Tree) Display(node *Node, level int) {
	if node != nil {
		for i := 0; i < level; i++ {
			fmt.Print("  ")
		}
		fmt.Printf("%s\n", node.URL)
		for _, record := range node.Records {
			for i := 0; i < level+1; i++ {
				fmt.Print("  ")
			}
			fmt.Printf("Param: %s\n", record.Param)
		}
		for _, child := range node.Children {
			t.Display(child, level+1)
		}
	}
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

	//  := &Tree{Root: &Node{URL: string(target)}}
	pretty.Println("Waiting for messages....", target)
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			// pretty.Println(string(msg.Value))
			req, err := parser.ParseRequest(string(msg.Value))
			if err != nil {
				pretty.Println(err)
			}
			pretty.Println(req)
			// t.AddPath(req.URL.Host, req.URL.RawQuery)
			// if target == req.Host {
			//work somthing
			// binary.X8(req.URL.Scheme + "://" + req.URL.Host + req.URL.Path)
			// }

		case err := <-partitionConsumer.Errors():
			pretty.Println("Received errors", err)
		}
		// t.Display(t.Root, 0)
	}
}
