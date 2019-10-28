package entity

import (
	"sort"
	"sync"
)

type MultiHash map[string]string

var multiHashMutex = &sync.Mutex{}

func (multiHash *MultiHash) AppendPart(partId, part string) {
	multiHashMutex.Lock()
	(*multiHash)[partId] = part
	//fmt.Println(partId, " ", part)
	multiHashMutex.Unlock()
}

func (multiHash *MultiHash) GetMultiHash() string {
	result := ""
	keys := make([]string, len(*multiHash))
	for key := range *multiHash {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		//fmt.Print(len(*multiHash))
		result += (*multiHash)[key]
	}
	//fmt.Println()
	return result
}
