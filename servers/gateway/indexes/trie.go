package indexes

import (
	// "log"
	"sort"
	"sync"
)

//TODO: implement a trie data structure that stores
//keys of type string and values of type int64

type int64set map[int64]struct{}

type Node struct {
	Key      rune
	Vals     int64set
	Children map[rune]*Node
	Parent   *Node
	mx       sync.RWMutex
}

type Trie struct {
	Root *Node
	size int
}

func NewTrie() *Trie {
	trie := &Node{
		// Vals:     int64set{},
		Children: make(map[rune]*Node),
	}
	return &Trie{
		Root: trie,
		size: 0,
	}
}

//Add adds a key and value to the trie.
func (t *Trie) Add(key string, value int64) {
	t.Root.mx.Lock()
	if len(key) != 0 {
		currNode := t.Root
		for _, letter := range key {
			if _, ok := currNode.Children[letter]; !ok { // letter didn't exist in trie before
				newNode := &Node{
					Vals:     make(int64set),
					Children: make(map[rune]*Node),
					Parent:   currNode,
				}
				newNode.Key = letter
				currNode.Children[letter] = newNode
			}
			currNode = currNode.Children[letter]
		}
		currNode.Vals.add(value)
		t.size++
	}
	t.Root.mx.Unlock()
}

type runeSlice []rune

func (rs runeSlice) Len() int           { return len(rs) }
func (rs runeSlice) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }
func (rs runeSlice) Less(i, j int) bool { return rs[i] < rs[j] }

//Find finds `max` values matching `prefix`. If the trie
//is entirely empty, or the prefix is empty, or max == 0,
//or the prefix is not found, this returns a nil slice.
func (t *Trie) Find(prefix string, max int) []int64 {
	t.Root.mx.RLock()
	defer t.Root.mx.RUnlock()
	if t.Root == nil || len(prefix) == 0 || max == 0 {
		return nil
	}
	firstLetter := []rune(prefix)
	if _, ok := t.Root.Children[firstLetter[0]]; !ok {
		return nil
	}

	ans := []int64{}
	endNode := t.Root.traverse(max, 0, prefix, &ans)
	if ans == nil {
		return ans
	}

	if len(ans) < max && len(endNode.Children) != 0 { // need to traverse down children with breadth first search
		keys := make(runeSlice, 0, len(endNode.Children))

		for k := range endNode.Children {
			keys = append(keys, k)
		}
		if len(endNode.Children) != 1 {
			sort.Sort(keys)
		}

		for _, key := range keys {
			endNode.Children[key].traverseForExtraChild(max, &ans)
		}
	}

	return ans
}

func (n *Node) traverse(max int, index int, prefix string, vals *[]int64) *Node {
	prefixRune := []rune(prefix)

	if index == len(prefix) { // found endNode
		for id := range n.Vals {
			if len(*vals) != max {
				*vals = append(*vals, id)
			}
		}
		return n
	} else if n == nil || len(n.Children) == 0 {
		*vals = nil
		return nil
	} else {
		childNode, ok := n.Children[prefixRune[index]]
		if !ok {
			*vals = nil
			return nil
		}
		return childNode.traverse(max, index+1, prefix, vals)
	}
}

func (n *Node) traverseForExtraChild(max int, ans *[]int64) {
	if len(*ans) != max && n != nil {
		if len(n.Children) != 0 {
			keys := make(runeSlice, 0, len(n.Children))

			for k := range n.Children {
				keys = append(keys, k)
			}
			// if len(n.Children) != 1 {
			sort.Sort(keys)
			// }

			for _, key := range keys {
				n.Children[key].traverseForExtraChild(max, ans)
			}
		} else if len(n.Vals) != 0 { // add all of the ids from node to answer
			for id := range n.Vals {
				if len(*ans) != max {
					*ans = append(*ans, id)
				}
			}
		}
	}
	return
}

//Remove removes a key/value pair from the trie
//and trims branches with no values.
func (t *Trie) Remove(key string, value int64) {
	t.Root.mx.RLock()
	defer t.Root.mx.RUnlock()
	ans := []int64{}
	endNode := t.Root.traverse(1, 0, key, &ans)
	if endNode != nil {
		if _, ok := endNode.Vals[value]; ok {
			endNode.Vals.remove(value)
			t.size--
			if len(endNode.Children) == 0 && len(endNode.Vals) == 0 {
				endNode.trimBranch(len(key)-1, key, t.Root)
			}
		}
	}
}

func (n *Node) trimBranch(index int, prefix string, root *Node) {
	prefixRune := []rune(prefix)
	if n != root && n.Key == prefixRune[index] {
		parent := n.Parent
		if len(n.Children) == 0 && len(n.Vals) == 0 {
			delete(parent.Children, n.Key)
		}
		parent.trimBranch(index-1, prefix, root)
	}
	return
}

//add adds a value to the set and returns
//true if the value didn't already exist in the set.
func (s int64set) add(value int64) bool {
	if exist := s.has(value); !exist {
		s[value] = struct{}{}
		return !exist
	}
	return false
}

//remove removes a value from the set and returns
//true if that value was in the set, false otherwise.
func (s int64set) remove(value int64) bool {
	if exist := s.has(value); exist {
		delete(s, value)
		return exist
	}
	return false
}

//has returns true if value is in the set,
//or false if it is not in the set.
func (s int64set) has(value int64) bool {
	_, exist := s[value]
	return exist
}