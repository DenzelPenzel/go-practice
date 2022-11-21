package main

import (
	"sort"
	"strings"
)

type Node struct {
	children map[string]*Node
	content  string
	isFile   bool
}

type FileSystem struct {
	root *Node
}

func Constructor() FileSystem {
	return FileSystem{
		root: &Node{
			children: make(map[string]*Node),
		},
	}
}

func getPath(path string) []string {
	paths := strings.Split(path, "/")
	var res []string
	for i := 0; i < len(paths); i++ {
		if paths[i] != "" {
			res = append(res, paths[i])
		}
	}
	return res
}

func (fs *FileSystem) Ls(path string) []string {
	paths := getPath(path)
	node := fs.root

	for _, key := range paths {
		child, ok := node.children[key]
		if !ok {
			return []string{}
		}
		node = child
	}

	if node.isFile {
		return []string{paths[len(paths)-1]}
	}

	var res []string

	for name := range node.children {
		res = append(res, name)
	}

	sort.Strings(res)
	return res
}

func (fs *FileSystem) Mkdir(path string) {
	paths := getPath(path)
	node := fs.root

	for _, key := range paths {
		child, ok := node.children[key]
		if !ok {
			child = &Node{children: make(map[string]*Node)}
			node.children[key] = child
		}
		node = child
	}
}

func (fs *FileSystem) AddContentToFile(filePath string, content string) {
	paths := getPath(filePath)
	node := fs.root

	for i, key := range paths {
		child, ok := node.children[key]
		if !ok {
			if len(paths) == i+1 {
				child = &Node{isFile: true, content: ""}
				node.children[key] = child
			} else {
				child = &Node{children: make(map[string]*Node)}
				node.children[key] = child
			}
		}
		node = child
	}

	if node.isFile {
		node.content += content
	}

}

func (fs *FileSystem) ReadContentFromFile(filePath string) string {
	paths := getPath(filePath)
	node := fs.root

	for _, key := range paths {
		node = node.children[key]
	}
	if node.isFile {
		return node.content
	}
	return ""
}
