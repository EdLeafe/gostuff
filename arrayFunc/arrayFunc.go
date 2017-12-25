package arrayFunc

import (
	"sort"
)

func IntIn(val int, source []int) bool {
    slc := make([]int, len(source))
    copy(slc, source)
    sort.Ints(slc)
	i := sort.Search(len(slc), func(i int) bool { return slc[i] >= val })
	return i < len(slc) && slc[i] == val
}

func StringIn(val string, source []string) bool {
    slc := make([]string, len(source))
    copy(slc, source)
    sort.Strings(slc)
	i := sort.Search(len(slc), func(i int) bool { return slc[i] >= val })
	return i < len(slc) && slc[i] == val
}
