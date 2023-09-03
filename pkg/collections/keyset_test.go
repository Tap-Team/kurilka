package collections_test

import (
	"slices"
	"testing"

	"github.com/Tap-Team/kurilka/pkg/collections"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_KeySet_New(t *testing.T) {
	cases := []struct {
		elems    []collections.Keyable
		len      int
		setelems []collections.Keyable
	}{
		{
			elems:    collections.NewKeyableIntList([]int{1, 1, 1, 1, 2, 3}),
			len:      3,
			setelems: collections.NewKeyableIntList([]int{1, 2, 3}),
		},
		{
			elems:    collections.NewKeyableStringList([]string{"Hello", "Hello.", "Hello World", "Hello", "Hello."}),
			len:      3,
			setelems: collections.NewKeyableStringList([]string{"Hello", "Hello.", "Hello World"}),
		},
	}

	for _, cs := range cases {
		set := collections.NewKeySet(cs.elems)
		require.Equal(t, cs.len, set.Len(), "wrong len")
		ok := slices.Equal(cs.setelems, set.Elements())
		require.True(t, ok, "slices not equal")
	}
}

func Test_KeySet_Add(t *testing.T) {
	cases := []struct {
		elems    []collections.Keyable
		setelems []collections.Keyable
	}{
		{
			elems:    collections.NewKeyableIntList([]int{1, 1, 1, 1, 2, 3}),
			setelems: collections.NewKeyableIntList([]int{1, 2, 3}),
		},
		{
			elems:    collections.NewKeyableStringList([]string{"Hello", "Hello.", "Hello World", "Hello", "Hello."}),
			setelems: collections.NewKeyableStringList([]string{"Hello", "Hello.", "Hello World"}),
		},
	}

	for _, cs := range cases {
		set := collections.NewKeySet([]collections.Keyable{})
		for _, element := range cs.elems {
			exists := set.Exist(element)
			ok := set.Add(element)
			assert.Equal(t, !exists, ok, "wrong result from add element")

			exists = set.Exist(element)
			assert.True(t, exists, "add not work")
		}
		ok := slices.Equal(cs.setelems, set.Elements())
		assert.True(t, ok, "slices not equal")
	}
}

func Test_KeySet_Remove(t *testing.T) {
	cases := []struct {
		elems          []collections.Keyable
		removeElements []collections.Keyable
		setelems       []collections.Keyable
	}{
		{
			elems:          collections.NewKeyableIntList([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
			removeElements: collections.NewKeyableIntList([]int{4, 7, 10, 154}),
			setelems:       collections.NewKeyableIntList([]int{1, 2, 3, 5, 6, 8, 9}),
		},
		{
			elems:          collections.NewKeyableStringList([]string{"ABC", "AMIDMAN", "HELLO World", "HELLO WORLD", "HIMAN", ""}),
			removeElements: collections.NewKeyableStringList([]string{"ABC", "AMIDMAN", "HELLO WORLD", "ABRACADABRA", "asdfjlaskd", ""}),
			setelems:       collections.NewKeyableStringList([]string{"HELLO World", "HIMAN"}),
		},
	}

	for _, cs := range cases {
		set := collections.NewKeySet(cs.elems)
		for _, element := range cs.removeElements {
			exists := set.Exist(element)
			ok := set.Remove(element)
			assert.Equal(t, exists, ok, "wrong result from remove")

			exists = set.Exist(element)
			assert.False(t, exists, "delete not work, %v", element)
		}
		ok := slices.Equal(cs.setelems, set.Elements())
		assert.True(t, ok, "slices not equal")
	}
}

func Test_KeySet_RemoveIndex(t *testing.T) {

}
