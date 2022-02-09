package main

import (
	"fmt"
	"sync"
)

func main() {
	var (
		m  = new(sync.Map)
		ds = make([]string, 0, 4)
	)
	const key = "key"
	m.Store(key, ds)
	var v, _ = m.Load(key)
	var s, ok = v.([]string)
	if !ok {
		return
	}
	s = append(s, key)
	m.Store(key, s)

	var v2, _ = m.Load(key)
	var s2, ok2 = v2.([]string)
	if !ok2 {
		return
	}
	fmt.Println(len(s2))
}
