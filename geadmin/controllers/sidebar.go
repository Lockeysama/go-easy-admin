package geacontrollers

import (
	"reflect"
	"sort"
	"strings"
)

// SideNode 侧栏节点
type SideNode struct {
	Path  string
	Title string
	Icon  string
	Child []*SideNode
}

// Tree 侧栏树
var Tree map[string]map[string]SideNode

func init() {
	Tree = make(map[string]map[string]SideNode)
}

// RegisterSideTree 控制器注册侧栏节点
func RegisterSideTree(c ControllerRolePolicy) {
	reflectVal := reflect.ValueOf(c)
	ct := reflect.Indirect(reflectVal).Type()
	prefix := strings.ToLower(c.Prefix()[1:])
	if _, ok := Tree[prefix]; !ok {
		Tree[prefix] = make(map[string]SideNode)
		parent := new(SideNode)
		parent.Path = prefix
		parent.Title = c.PrefixAlias()
		if parent.Title == "" {
			parent.Title = strings.TrimSuffix(ct.Name(), "Controller")
		}
		parent.Icon = c.PrefixIcon()
		Tree[prefix]["__ParentNode"] = *parent
	}
	child := new(SideNode)
	child.Path = c.AdminName()
	if child.Path == "" {
		child.Path = strings.ToLower(strings.TrimSuffix(ct.Name(), "Controller"))
	}
	child.Title = c.AdminNameAlias()
	if child.Title == "" {
		child.Title = strings.TrimSuffix(ct.Name(), "Controller")
	}
	child.Icon = c.AdminIcon()

	Tree[prefix][child.Path] = *child
}

// SideTree 获取授权侧栏树
func SideTree(path map[string][]string) *[]SideNode {
	// TODO 速度优化
	trees := new([]SideNode)
	keys := reflect.ValueOf(path).MapKeys()
	keysOrder := func(i, j int) bool { return keys[i].Interface().(string) < keys[j].Interface().(string) }
	sort.Slice(keys, keysOrder)
	for _, prefix := range keys {
		paths := path[prefix.Interface().(string)]
		if nodesMap, ok := Tree[prefix.Interface().(string)]; ok {
			parent := nodesMap["__ParentNode"]
			for _, _path := range paths {
				child := nodesMap[_path]
				parent.Child = append(parent.Child, &child)
			}
			*trees = append(*trees, parent)
		}
	}
	return trees
}
