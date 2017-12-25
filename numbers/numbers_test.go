package numbers

import (
    "testing"
)

func TestRandRange(t *testing.T) {
    SeedRandom()
    low := 42
    high := 55
    num := RandRange(low, high)
    if num < low || num > high {
        t.Errorf("RandRange returned %d when it should have been between %d and %d",
                num, low, high)
    }
}

func TestRandRangeReverse(t *testing.T) {
    SeedRandom()
    low := 55
    high := 42
    num := RandRange(low, high)
    // RandRange should reverse them for a valid range
    if num < high  || num > low {
        t.Errorf("RandRange returned %d when it should have been between %d and %d",
                num, low, high)
    }
}

func TestRandRangeSame(t *testing.T) {
    SeedRandom()
    low := 42
    num := RandRange(low, low)
    if num != low {
        t.Errorf("RandRange returned %d when it should have been %d", num, low)
    }
}

func TestRandSet(t *testing.T) {
    SeedRandom()
    low := 42
    high := 55
    count := 4
    set := RandSet(low, high, count)
    if len(set) != count {
        t.Errorf("RandSet returned %d items; should have been %d", len(set), count)
    }
    first := set[0]
    for _, next := range set[1:] {
        if next < first {
            t.Errorf("RandSet is not sorted: %v", set)
        }
        first = next
    }
}
