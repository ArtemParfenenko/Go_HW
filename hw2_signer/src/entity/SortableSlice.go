package entity

import "sort"

type SortableSlice []string

func (sortableSlice *SortableSlice) Add(value string) {
	*sortableSlice = append(*sortableSlice, value)
}

func (sortableSlice *SortableSlice) Sort() {
	sort.Strings(*sortableSlice)
}

func (sortableSlice *SortableSlice) String() string {
	result := ""
	for _, value := range *sortableSlice {
		result = result + value + "_"
	}
	return result[:len(result)-1]
}
