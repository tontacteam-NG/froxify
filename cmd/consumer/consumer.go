package main

import (
	"log"
	"slices"
	"sync"

	"github.com/IBM/sarama"
	"github.com/kr/pretty"
	"github.com/nxczje/froxy/core"
	"github.com/nxczje/froxy/parser"
)

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

// var fu_channel = make(chan string, 2)

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
// 			child = &Node{URL: component, Param: make([]string, 0)}
// 			current.Children[component] = child
// 		}

// 		// Di chuyển xuống nút con
// 		current = child
// 	}
// 	if param == "" {
// 		return
// 	}
// 	// Lưu trữ thông tin tham số tại nút lá
// 	current.Param = append(current.Param, param)
// }

// // func SearchNode(node *Node, targetURL string) *Node {
// // 	if node == nil {
// // 		return nil
// // 	}

// // 	if node.URL == targetURL {
// // 		return node
// // 	}

// // 	for _, child := range node.Children {
// // 		if found := SearchNode(child, targetURL); found != nil {
// // 			if slices.Contains(s) {
// // 				return found
// // 			}
// // 		}
// // 	}

// //		return nil
// //	}
// func (n *Node) SearchNode(targetURL string) *Node {

// 	components := strings.Split(targetURL, "/")
// 	current := n

// 	// pretty.Println(components)
// 	for _, component := range components {
// 		if component == "" {
// 			continue
// 		}
// 		child, exists := current.Children[component]
// 		if !exists {
// 			// Nếu không tồn tại, thêm một nút mới
// 			return nil
// 		}
// 		current = child
// 		// Di chuyển xuống nút con
// 	}
// 	return current
// }

// Hiển thị cây nhị phân
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

// func saveRequest(req []byte) {
// 	file, err := os.CreateTemp("/tmp/sqlmap", "sqlmap")
// 	defer file.Close()
// 	if err != nil {
// 		log.Println(err.Error())
// 	}
// 	_, err = file.Write(req)
// 	if err != nil {
// 		log.Println(err.Error())
// 	}
// 	go pushFuzz(file.Name())

// }

// func runSqlmap() {
// 	var max = make(chan string, 5)
// 	for i := range fu_channel {
// 		max <- "Job"
// 		cmd := exec.Command("osmedeus", "scan", "-f", "froxy", "-t", i)
// 		go func(cmd exec.Cmd) {
// 			fmt.Println(cmd.String())

// 			err := cmd.Start()
// 			if err != nil {
// 				fmt.Println(err.Error())
// 			}
// 			err = cmd.Wait()
// 			if err != nil {
// 				fmt.Println(err.Error())
// 			}
// 			<-max
// 		}(*cmd)

// 	}
// }

// func pushFuzz(filename string) {
// 	fu_channel <- filename
// 	fmt.Println("pushed")
// }
