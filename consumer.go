package main

// import (
// 	"fmt"
// 	"log"
// 	"strings"
// 	"sync"

// 	"github.com/IBM/sarama"
// 	"github.com/kr/pretty"
// 	"github.com/nxczje/froxy/parser"
// )

// // Node là một nút trong cây nhị phân
// type Node struct {
// 	URL      string
// 	Param    []string
// 	Children map[string]*Node
// }

// // Tree là cấu trúc cây nhị phân
// type Tree struct {
// 	Root *Node
// 	mu   *sync.Mutex // Sử dụng mutex để đồng bộ hóa truy cập cây
// }

// // Thêm đường dẫn vào cây nhị phân với goroutines
// func (t *Tree) AddPath(url, param string) {

// 	t.mu.Lock()
// 	defer t.mu.Unlock()

// 	components := strings.Split(url, "/")
// 	current := t.Root

// 	for _, component := range components {
// 		if component == "" {
// 			continue
// 		}
// 		if current.Children == nil {
// 			current.Children = make(map[string]*Node)
// 		}

// 		// Kiểm tra xem nút con có tồn tại chưa
// 		child, exists := current.Children[component]
// 		if !exists {
// 			// Nếu không tồn tại, thêm một nút mới
// 			child = &Node{URL: component}
// 			current.Children[component] = child
// 		}

// 		// Di chuyển xuống nút con
// 		current = child
// 	}

// 	// Lưu trữ thông tin tham số tại nút lá
// 	current.Param = append(current.Param, param)
// }

// // Hiển thị cây nhị phân
// func (t *Tree) Display(node *Node, level int) {
// 	if node != nil {
// 		for i := 0; i < level; i++ {
// 			fmt.Print("  ")
// 		}
// 		fmt.Printf("%s\n", node.URL)
// 		for _, record := range node.Param {
// 			for i := 0; i < level+1; i++ {
// 				fmt.Print("  ")
// 			}
// 			fmt.Printf("Param: %s\n", record)
// 		}
// 		for _, child := range node.Children {
// 			t.Display(child, level+1)
// 		}
// 	}
// }

// func main() {
// 	target := "nothinnn.oob.nncg.uk"       //change this
// 	Addr := []string{"10.14.140.135:9092"} //change this
// 	Topic := "nothing"                     //change this
// 	config := sarama.NewConfig()
// 	config.Consumer.Offsets.Initial = sarama.OffsetOldest
// 	mu := sync.Mutex{}
// 	trees := make(map[string]Tree, 0)
// 	consumer, err := sarama.NewConsumer(Addr, config)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	partitionConsumer, err := consumer.ConsumePartition(Topic, 0, sarama.OffsetNewest)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	//  := &Tree{Root: &Node{URL: string(target)}}
// 	pretty.Println("Waiting for messages....", target)
// 	for {
// 		select {
// 		case msg := <-partitionConsumer.Messages():
// 			req, err := parser.ParseRequest(string(msg.Value))
// 			if err != nil {
// 				// pretty.Println(err)
// 				continue
// 			}
// 			mu.Lock()
// 			tree, exsit := trees[req.Host]
// 			if exsit {
// 				tree.AddPath(req.RequestURI, req.URL.RawPath)
// 			} else {
// 				trees[req.Host] = Tree{mu: &sync.Mutex{}, Root: &Node{URL: req.Host}}
// 			}
// 			// pretty.Println(req)

// 			// t.AddPath(req.URL.Host, req.URL.RawQuery)
// 			// if target == req.Host {
// 			//work somthing
// 			// binary.X8(req.URL.Scheme + "://" + req.URL.Host + req.URL.Path)
// 			// }
// 			mu.Unlock()
// 			tree.Display(tree.Root, 0)
// 			fmt.Println("----------------------------------------")

// 		case err := <-partitionConsumer.Errors():
// 			pretty.Println("Received errors", err)
// 		}
// 		// t.Display(t.Root, 0)
// 	}
// }
