package collections

import (
	"hash/fnv"
)

type Keyable interface {
	Key() int
}

type KeyableInt int

func (k KeyableInt) Key() int {
	return int(k)
}

type KeyableString string

func (k KeyableString) Key() int {
	hash := fnv.New64()
	hash.Write([]byte(k))
	return int(hash.Sum64())
}

func NewKeyableIntList(l []int) []Keyable {
	keyableList := make([]Keyable, 0, len(l))
	for _, i := range l {
		keyableList = append(keyableList, KeyableInt(i))
	}
	return keyableList
}

func NewKeyableStringList(l []string) []Keyable {
	keyableList := make([]Keyable, 0, len(l))
	for _, i := range l {
		keyableList = append(keyableList, KeyableString(i))
	}
	return keyableList
}
