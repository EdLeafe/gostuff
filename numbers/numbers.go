package numbers

import (
	"math/rand"
	"sort"
	"time"
)

func SeedRandom() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func RandRange(low, high int) int {
	if low > high {
        low, high = high, low
    } else if low == high {
        return low
    }
    diff := high - low
	return rand.Intn(int(diff+1)) + low
}

func RandSet(low, high, count int) []int {
	set := make(map[int]struct{})
	dummy := struct{}{}
	for len(set) < count {
		num := RandRange(low, high)
		set[num] = dummy
	}

	keys := make([]int, count)
	i := 0
	for k := range set {
		keys[i] = k
		i++
	}
	sort.Ints(keys)
	return keys
}

func Distribution(low, high, total int) []int {
    SeedRandom()
    count := make([]int, total)

    for i:=0; i < total; i++ {
        rr := RandRange(low, high)
        count[rr-1] += 1
    }
    return count
}
