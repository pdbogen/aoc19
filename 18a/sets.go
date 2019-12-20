package main

type KeySet uint64

func (k KeySet) String() string {
	ret := ""
	for _, key := range KeyBytes {
		if k & Keys[key] > 0 {
			ret += string(key)
		}
	}
	return ret
}

var Keys = map[byte]KeySet{}
var KeySymbols = map[KeySet]byte{}
var KeyBytes = []byte("abcdefghijklmnopqrstuvwxyz")

func init() {
	for i, k := range KeyBytes {
		Keys[k] = 1 << i
		KeySymbols[1<<i] = k
	}
}
