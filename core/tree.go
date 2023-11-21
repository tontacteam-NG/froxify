package core

import (
	"strings"
	"sync"
)

type Node struct {
	URL      string
	Param    []string
	Children map[string]*Node
}

// tree path
type Tree struct {
	Root *Node
	Mu   *sync.Mutex // Sử dụng mutex để đồng bộ hóa truy cập cây
}

func (t *Tree) AddPath(url, param string) {

	t.Mu.Lock()
	defer t.Mu.Unlock()

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

func (n *Node) SearchNode(targetURL string) *Node {

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
