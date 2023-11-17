package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"
	"sync"

	"github.com/IBM/sarama"
	"github.com/kr/pretty"
	"github.com/nxczje/froxy/parser"
)

// Node là một nút trong cây nhị phân
type Node struct {
	URL      string
	Param    []string
	Children map[string]*Node
}

// Tree là cấu trúc cây nhị phân
type Tree struct {
	Root *Node
	mu   *sync.Mutex // Sử dụng mutex để đồng bộ hóa truy cập cây
}

var fu_channel = make(chan string, 2)

// Thêm đường dẫn vào cây nhị phân với goroutines
func (t *Tree) AddPath(url, param string) {

	t.mu.Lock()
	defer t.mu.Unlock()

	components := strings.Split(url, "/")
	current := t.Root

	for _, component := range components {
		if component == "" {
			continue
		}
		if current.Children == nil {
			current.Children = make(map[string]*Node)
		}

		// Kiểm tra xem nút con có tồn tại chưa
		child, exists := current.Children[component]
		if !exists {
			// Nếu không tồn tại, thêm một nút mới
			child = &Node{URL: component, Param: make([]string, 0)}
			current.Children[component] = child
		}

		// Di chuyển xuống nút con
		current = child
	}
	if param == "" {
		return
	}
	// Lưu trữ thông tin tham số tại nút lá
	current.Param = append(current.Param, param)
}

// func SearchNode(node *Node, targetURL string) *Node {
// 	if node == nil {
// 		return nil
// 	}

// 	if node.URL == targetURL {
// 		return node
// 	}

// 	for _, child := range node.Children {
// 		if found := SearchNode(child, targetURL); found != nil {
// 			if slices.Contains(s) {
// 				return found
// 			}
// 		}
// 	}

//		return nil
//	}
func (n *Node) SearchNode1(targetURL string) *Node {

	components := strings.Split(targetURL, "/")
	current := n

	// pretty.Println(components)
	for _, component := range components {
		if component == "" {
			continue
		}
		child, exists := current.Children[component]
		if !exists {
			// Nếu không tồn tại, thêm một nút mới
			return nil
		}
		current = child
		// Di chuyển xuống nút con
	}
	return current
}

// Hiển thị cây nhị phân
func (t *Tree) Display(node *Node, level int) {
	if node != nil {
		for i := 0; i < level; i++ {
			fmt.Print("  ")
		}
		fmt.Printf("%s\n", node.URL)
		for _, record := range node.Param {
			for i := 0; i < level+1; i++ {
				fmt.Print("  ")
			}
			fmt.Printf("Param: %s\n", record)
		}
		for _, child := range node.Children {
			t.Display(child, level+1)
		}
	}
}

func main() {
	target := "nothinnn.oob.nncg.uk"       //change this
	Addr := []string{"10.14.140.245:9092"} //change this
	Topic := "nothing"                     //change this
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	mu := sync.Mutex{}
	trees := make(map[string]Tree, 0)
	consumer, err := sarama.NewConsumer(Addr, config)
	if err != nil {
		log.Fatalln(err)
	}

	partitionConsumer, err := consumer.ConsumePartition(Topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalln(err)
	}
	go runSqlmap()
	//  := &Tree{Root: &Node{URL: string(target)}}
	pretty.Println("Waiting for messages....", target)
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
				n := tree.Root.SearchNode1(req.URL.Path)
				if n != nil && slices.Contains(n.Param, req.URL.RawQuery) {
					// tree.Display(tree.Root, 0)
					mu.Unlock()
					continue

				}
				mu.Unlock()
				tree.AddPath(req.URL.Path, req.URL.RawQuery)
				saveRequest(msg.Value)
			} else {
				tree = Tree{mu: &sync.Mutex{}, Root: &Node{URL: req.Host}}
				mu.Unlock()
				tree.AddPath(req.URL.Path, req.URL.RawQuery)
				saveRequest(msg.Value)
				trees[req.Host] = tree
			}
			// pretty.Println(req)

			// t.AddPath(req.URL.Host, req.URL.RawQuery)
			// if target == req.Host {
			//work somthing
			// binary.X8(req.URL.Scheme + "://" + req.URL.Host + req.URL.Path)
			// }
			// tree.Display(tree.Root, 0)

		case err := <-partitionConsumer.Errors():
			pretty.Println("Received errors", err)
		}
		// t.Display(t.Root, 0)
	}
}

func saveRequest(req []byte) {
	file, err := os.CreateTemp("/tmp/sqlmap", "sqlmap")
	defer file.Close()
	if err != nil {
		log.Println(err.Error())
	}
	_, err = file.Write(req)
	if err != nil {
		log.Println(err.Error())
	}
	go pushFuzz(file.Name())

}

func runSqlmap() {
	var max = make(chan string, 5)
	for i := range fu_channel {
		max <- "123"
		cmd := exec.Command("sqlmap", "--batch", "--level", "1", "--technique=BET", "--skip=\"cookie,referer,host\"", "--parse-errors", "-r", i)
		go func(cmd exec.Cmd) {
			fmt.Println(cmd.String())
			err := cmd.Start()
			if err != nil {
				fmt.Println(err.Error())
			}
			err = cmd.Wait()
			if err != nil {
				fmt.Println(err.Error())
			}
			<-max
		}(*cmd)

	}
}

func pushFuzz(filename string) {
	fu_channel <- filename
	fmt.Println("pushed")
}
