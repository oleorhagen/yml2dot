package renderer

import (
	"fmt"
	"os"
	"strings"

	"github.com/lucasepe/dot"
	"gopkg.in/yaml.v2"
)

var nodes = map[string]*dot.Node{}

var nodeNames [2]string

// Render returns a GraphViz representation of a YAML tree.
func Render(v yaml.MapSlice) *dot.Graph {
	g := dot.NewGraph(dot.Undirected)
	g.Attr("nodesep", "4")
	// g.Attr("rankdir", "LR")
	g.Attr("pad", "0.5")
	g.Attr("ranksep", "0.5")
	g.Attr("fontname", "Fira Mono")
	g.Attr("fontsize", "11")

	g.NodeBaseAttrs().
		Attr("fontname", "Fira Mono").
		Attr("fontsize", "10").
		Attr("margin", "0.3,0.1").
		Attr("fillcolor", "#fafafa").
		Attr("shape", "box").
		Attr("penwidth", "2.0").
		Attr("style", "rounded,filled")

	// fmt.Fprintf(os.Stderr, "Render: v: %v\n", v)
	for _, el := range v {
		// fmt.Fprintf(os.Stderr, "Render: El: %v", el)
		topLevelNode = fmt.Sprintf("%v", el.Key)
		// fmt.Fprintf(os.Stderr, "Render: topLevelNode: %s\n", topLevelNode)
		fmt.Fprintf(os.Stderr, "Render: subvals: %[1]T: %[1]v\n", el.Value)
		for _, val := range el.Value.(yaml.MapSlice) {
			nodeNames[0] = topLevelNode
			renderMapItem(val, g, nil, 0)
		}
	}

	return g
}

func NewNode(name string, g *dot.Graph, depth int, withAttrs ...func(*dot.AttributesMap)) *dot.Node {
	if name == "release_component" {
		return nil
	}
	if name == "independent_component" {
		return nil
	}
	if name == "non_core_component" {
		return nil
	}
	if strings.Contains(name, "true") || strings.Contains(name, "false") {
		return nil
	}
	if name == "" {
		return nil
	}
	fmt.Fprintf(os.Stderr, "NewNode: name: %s\n", name)
	fmt.Fprintf(os.Stderr, "NewNode: topLevelName: %s\n", topLevelNode)
	fmt.Fprintf(os.Stderr, "NewNode: nodeSuffixes: %s depth: %d\n", nodeNames, depth)
	// if !strings.Contains(name, ":") {
	// 	name = name + ":" + topLevelNode
	// }
	if depth == 1 || depth == 0 {
		name = name + ":" + nodeNames[0]
	} else if depth == 4 {

		name = name + ":" + nodeNames[1]
	} else {
		panic("Logic error!")
	}
	if n, ok := nodes[name]; ok {
		fmt.Fprintf(os.Stderr, "Returning old: name: %s -> %v\n", name, nodes[name])
		return n
	}
	fmt.Fprintf(os.Stderr, "Producing new node: %s\n", name)
	child := g.Node()
	child.Attr("label", name)
	nodes[name] = child
	return child
}

var topLevelNode string

func renderMapItem(v yaml.MapItem, g *dot.Graph, parent *dot.Node, depth int) {
	name := fmt.Sprintf("%v", v.Key)
	// Skip the mid-level meta-data nodes:
	if name == "docker_image" || name == "git" || name == "docker_container" {
		nodeNames[1] = name
		renderVal(v.Value, g, parent, depth+1)
		nodeNames[1] = ""
		return
	}
	child := NewNode(name, g, depth)
	if child == nil {
		return
	}
	if parent != nil {
		// TODO
		// Name the child after either the top level, or if the it is a
		// child, suffix the parent name.
		//
		if v.Value == nil {
			fmt.Fprintf(os.Stderr, "renderMapItem: nodeNames: %v\n", nodeNames)
			fmt.Fprintf(os.Stderr, "renderMapItem: Leaf node: %v\n", v)
		} else {
			fmt.Fprintf(os.Stderr, "renderMapItem: Not a leaf node: %v\n", v)
		}
		// child.Attr("label", name)

		// if len(g.FindEdges(*parent, *child)) == 0 {
		fmt.Fprintf(os.Stderr, "renderMapItem: Edge: %v:%s -> %v:%s depth: %d\n", parent, nodeNames[0], child, nodeNames[1], depth)
		link := g.Edge(parent, child)
		link.Attr("arrowhead", "none")
		link.Attr("penwidth", "2.0")
		// }
	} else {
		// Top level (git, docker_image, docker_container ) should not be rendered
		fmt.Fprintf(os.Stderr, "renderMapItem: Top level node %v\n", v)
		fmt.Fprintf(os.Stderr, "renderMapItem: Top level name %s\n", topLevelNode)
		// child.Attr("label", name+":"+topLevelNode)
		// child.Attr("label", dot.HTML(fmt.Sprintf("<b>%v</b>", name)))
		// child.Attr("shape", "plaintext")
		// child.Attr("style", "")
	}

	renderVal(v.Value, g, child, depth)
}

func renderVal(v interface{}, g *dot.Graph, parent *dot.Node, depth int) {

	switch v.(type) {
	case []interface{}:
		renderSlice(v.([]interface{}), g, parent, depth+1)
	case yaml.MapSlice:
		for _, el := range v.(yaml.MapSlice) {
			renderMapItem(el, g, parent, depth+1)
		}
	case map[string]interface{}:
		renderMap(v.(map[string]interface{}), g, parent, depth)
	default:
		if parent != nil {
			name := fmt.Sprintf("%v", v)
			// name := fmt.Sprintf("%v%s", v, parent.Attrs().Value("label").(string))
			child := NewNode(name, g, depth)
			if child == nil {
				return
			}
			fmt.Fprintf(os.Stderr, "renderVal: Leaf node: %v\n", v)
			fmt.Fprintf(os.Stderr, "renderMapItem: nodeNames: %v\n", nodeNames)
			// child.Attr("label", name)
			fmt.Fprintf(os.Stderr, "renderVal: Leaf node name: %s\n", name)
			// child.Attr("label", name+":"+parent.Attrs().Value("label").(string))
			// if len(g.FindEdges(*parent, *child)) == 0 {
			fmt.Fprintf(os.Stderr, "renderVal: Adding edge from: %s -> %s depth: %d\n", parent.Attrs().Value("label").(string), name, depth)

			fmt.Fprintf(os.Stderr, "renderVal: Edge: %v:%s -> %v:%s\n", parent, nodeNames[0], child, nodeNames[1])
			link := g.Edge(parent, child)
			link.Attr("arrowhead", "none")
			link.Attr("penwidth", "2.0")
			// }
		} else {
			fmt.Fprintf(os.Stderr, "renderVal: Top level node %v\n", v)
		}
	}
}

func renderMap(m map[string]interface{}, g *dot.Graph, parent *dot.Node, depth int) {
	for k, v := range m {
		name := fmt.Sprintf("%v", k)
		child := NewNode(name, g, depth, dot.WithLabel(name))
		if child == nil {
			return
		}
		if parent != nil {
			// if len(g.FindEdges(*parent, *child)) == 0 {
			fmt.Fprintf(os.Stderr, "renderMap: Adding Edge: %v:%s -> %v:%s depth: %d\n", parent, nodeNames[0], child, nodeNames[1], depth)
			link := g.Edge(parent, child)
			link.Attr("arrowhead", "none")
			link.Attr("penwidth", "2.0")
			// }
		} else {
			fmt.Fprintf(os.Stderr, "renderMap: Top level node %v\n", v)
		}
		renderVal(v, g, child, depth+1)
	}
}

func renderSlice(slc []interface{}, g *dot.Graph, parent *dot.Node, depth int) {
	for _, v := range slc {
		renderVal(v, g, parent, depth+1)
	}
}
