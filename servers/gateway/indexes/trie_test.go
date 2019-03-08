package indexes

import (
	// "log"
	// "reflect"
	// "sort"
	"testing"
)

func TestTrieAdd(t *testing.T) {
	cases := []struct {
		name string
		keys []string
		vals []int64
	}{
		{
			"Valid Trie from Slides",
			[]string{"go", "git", "gob", "go", "goal", "foo"},
			[]int64{1, 2, 3, 4, 5, 1},
		},
		{
			"Single Value Lower - 2 letters",
			[]string{"go"},
			[]int64{1},
		},
		{
			"Mutliple Value Lower",
			[]string{"neha", "yadav", "hates", "sushi"},
			[]int64{1, 2, 3, 4},
		},
		{
			"Key is Null",
			[]string{""},
			[]int64{1},
		},
	}
	for _, c := range cases {
		trie := NewTrie()
		for i, k := range c.keys {
			trie.Add(k, c.vals[i])
			currNode := trie.Root
			if c.name == "Key is Null" {
				if len(currNode.Vals) != 0 {
					t.Errorf("case %s: added a null key", c.name)
				}
			} else {
				for _, letter := range k {
					if _, ok := currNode.Children[letter]; ok {
						currNode = currNode.Children[letter]
					} else {
						t.Errorf("case %s: didn't add letter into trie", c.name)
					}
				}
				if _, ok := currNode.Vals[c.vals[i]]; !ok {
					t.Errorf("case %s: didn't add the right id into trie: expected %d, returned %v", c.name, c.vals[i], currNode.Vals)
				}
			}
		}
	}
}

func TestTrieRemove(t *testing.T) {
	cases := []struct {
		name            string
		inputKey        []string
		inputVal        []int64
		inputPrefixFind string
		keys            []string
		vals            []int64
		expectedVals    []int64
	}{
		{
			"Valid Trie from Slides w/ Trim",
			[]string{"goal"},
			[]int64{5},
			"go",
			[]string{"go", "git", "gob", "go", "goal", "foo"},
			[]int64{1, 2, 3, 4, 5, 1},
			[]int64{1, 3, 4,},
		},
		{
			"Single Value Lower",
			[]string{"neha"},
			[]int64{1},
			"n",
			[]string{"neha"},
			[]int64{1},
			[]int64{},
		},
		{
			"Remove IDs but don't Trim bc of Existing Children",
			[]string{"go", "go"},
			[]int64{1, 4},
			"go",
			[]string{"go", "git", "gob", "go", "goal", "foo"},
			[]int64{1, 2, 3, 4, 5, 1},
			[]int64{3, 5},
		},
		{
			"Remove IDs but don't Trim bc of Existing IDs",
			[]string{"go"},
			[]int64{1},
			"go",
			[]string{"go", "git", "gob", "go", "goal", "foo"},
			[]int64{1, 2, 3, 4, 5, 1},
			[]int64{4, 5},
		},
		{
			"Doesn't Remove bc of Wrong ID",
			[]string{"go"},
			[]int64{5},
			"go",
			[]string{"go", "git", "gob", "go", "goal", "foo"},
			[]int64{1, 2, 3, 4, 5, 1},
			[]int64{1, 4},
		},
		{
			"Doesn't Remove bc of Wrong Key",
			[]string{"goat"},
			[]int64{1},
			"go",
			[]string{"go", "git", "gob", "go", "goal", "foo"},
			[]int64{1, 2, 3, 4, 5, 1},
			[]int64{1, 4, 3, 5},
		},
	}
	for _, c := range cases {
		trie := NewTrie()
		for i, k := range c.keys {
			trie.Add(k, c.vals[i])
		}
		for i, inputKey := range c.inputKey { // more than one input key
			trie.Remove(inputKey, c.inputVal[i])

			results := trie.Find(c.inputPrefixFind, len(inputKey))
			
			if len(c.expectedVals) != len(results) {
				t.Errorf("case %s: incorrect length when comparing results(%v) and expectedVals(%v)",
					c.name, results, c.expectedVals)
			}
		}
	}
}
func TestTrieFind(t *testing.T) {
	cases := []struct {
		name         string
		inputKey     string
		inputVal     int
		keys         []string
		vals         []int64
		expectedVals []int64
		expectedErr  bool
	}{
		{
			"Valid Trie from Slides w/ N < IDS",
			"go",
			1,
			[]string{"go", "git", "gob", "go", "goal", "foo"},
			[]int64{1, 2, 3, 4, 5, 1},
			[]int64{1},
			false,
		},
		{
			"Valid Trie from Slides w/ N > IDS",
			"go",
			5,
			[]string{"go", "git", "gob", "go", "goal", "foo"},
			[]int64{1, 2, 3, 4, 5, 1},
			[]int64{1, 4, 5, 3},
			false,
		},
		{
			"Valid Trie from Slides w/ N > IDS",
			"go",
			3,
			[]string{"go", "git", "gob", "go", "goal", "foo"},
			[]int64{1, 2, 3, 4, 5, 1},
			[]int64{1, 4, 5},
			false,
		},
		{
			"Valid Trie from Slides w/ N near IDS",
			"go",
			4,
			[]string{"go", "git", "gob", "go", "goal", "foo"},
			[]int64{1, 2, 3, 4, 5, 1},
			[]int64{1, 4, 5, 3},
			false,
		},
		{
			"Valid Trie from Slides w/ N near IDS",
			"goal",
			4,
			[]string{"go", "git", "gob", "go", "goal", "foo", "goalie"},
			[]int64{1, 2, 3, 4, 5, 1, 6},
			[]int64{5, 6},
			false,
		},
		{
			"Only 1 Value",
			"neha",
			1,
			[]string{"neha"},
			[]int64{1},
			[]int64{1},
			false,
		},
		{
			"Empty Trie",
			"go",
			1,
			[]string{},
			[]int64{},
			[]int64{},
			true,
		},
		{
			"Prefix is Empty",
			"",
			1,
			[]string{"go", "git", "gob", "go", "goal", "foo"},
			[]int64{1, 2, 3, 4, 5, 1},
			[]int64{},
			true,
		},
		{
			"Max is Zero",
			"go",
			0,
			[]string{"go", "git", "gob", "go", "goal", "foo"},
			[]int64{1, 2, 3, 4, 5, 1},
			[]int64{},
			true,
		},
		{
			"Not Found Prefix",
			"goat",
			1,
			[]string{"go", "git", "gob", "go", "goal", "foo"},
			[]int64{1, 2, 3, 4, 5, 1},
			[]int64{},
			true,
		},
	}
	for _, c := range cases {
		trie := NewTrie()
		testset := int64set{}
		for i, k := range c.keys {
			trie.Add(k, c.vals[i])
		}
		for _, expectedID := range c.expectedVals {
			testset.add(expectedID)
		}

		results := trie.Find(c.inputKey, c.inputVal)
		if !c.expectedErr && results == nil {
			t.Errorf("case %s: wasn't expecting an error but results was nil", c.name)
		}
		if c.expectedErr && results != nil {
			t.Errorf("case %s: was expecting an error but results was not nil", c.name)
		}
		if len(c.expectedVals) != len(results) {
			t.Errorf("case %s: incorrect length when comparing results(%v) and expectedVals(%v)",
				c.name, results, c.expectedVals)
		}
	}
}

func Testint64SetAddHas(t *testing.T) {
	cases := []struct {
		name     string
		values   []int64
		expected []int64
	}{
		{
			"Single Value",
			[]int64{1},
			[]int64{1},
		},
		{
			"Duplicate Values",
			[]int64{1, 1},
			[]int64{1},
		},
		{
			"Distinct Values",
			[]int64{1, 2, 3},
			[]int64{1, 2, 3},
		},
		{
			"Distinct Values Then Duplicates",
			[]int64{1, 2, 3, 1, 2, 3},
			[]int64{1, 2, 3},
		},
		{
			"Duplicates Then Distinct Values",
			[]int64{1, 1, 2, 2, 3, 3},
			[]int64{1, 2, 3},
		},
		{
			"int64ermixed",
			[]int64{1, 2, 1, 3, 4, 3},
			[]int64{1, 2, 3, 4},
		},
	}

	for _, c := range cases {
		testset := int64set{}
		for _, v := range c.values {
			expectedRet := !testset.has(v)
			ret := testset.add(v)
			if expectedRet != ret {
				t.Errorf("case %s: incorrect return value when adding %d: expected %t but got %t",
					c.name, v, expectedRet, ret)
			}
		}
		if len(testset) != len(c.expected) {
			t.Errorf("case %s: incorrect length: expected %d but got %d",
				c.name, len(c.expected), len(testset))
		}
		for _, v := range c.expected {
			if !testset.has(v) {
				t.Errorf("case %s: expected value %d is not in the set",
					c.name, v)
			}
		}
	}
}

func Testint64SetRemove(t *testing.T) {
	cases := []struct {
		name     string
		values   []int64
		toRemove int64
		expected []int64
	}{
		{
			"One Removed from Many",
			[]int64{1, 2, 3},
			2,
			[]int64{1, 3},
		},
		{
			"Last One",
			[]int64{1},
			1,
			[]int64{},
		},
		{
			"Not Found",
			[]int64{1},
			2,
			[]int64{1},
		},
		{
			"Empty",
			[]int64{},
			2,
			[]int64{},
		},
	}

	for _, c := range cases {
		testset := int64set{}
		for _, v := range c.values {
			testset.add(v)
		}
		expectedRet := testset.has(c.toRemove)
		ret := testset.remove(c.toRemove)
		if expectedRet != ret {
			t.Errorf("case %s: incorrect return value when removing %d: expected %t but got %t",
				c.name, c.toRemove, expectedRet, ret)
		}
		if len(testset) != len(c.expected) {
			t.Errorf("case %s: incorrect length after remove: expected %d but got %d",
				c.name, len(c.expected), len(testset))
		}
		for _, v := range c.expected {
			if !testset.has(v) {
				t.Errorf("case %s: expected value %d was not in set after remove",
					c.name, v)
			}
		}
	}
}