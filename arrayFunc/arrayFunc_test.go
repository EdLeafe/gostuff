package arrayFunc

import "testing"

func TestIntIn(t *testing.T) {
	var response bool
	sample := []int{11, 22, 33, 44, 55}
	response = IntIn(33, sample)
	if !response {
		t.Errorf("The IntIn function should have found 33 in the sample slice")
	}
	response = IntIn(34, sample)
	if response {
		t.Errorf("The IntIn function should have found 33 in the sample slice")
	}
}

func TestIntInUnsorted(t *testing.T) {
	var response bool
	sample := []int{44, 22, 55, 11, 33}
	response = IntIn(33, sample)
	if !response {
		t.Errorf("The IntIn function should have found 33 in the sample slice")
	}
	response = IntIn(34, sample)
	if response {
		t.Errorf("The IntIn function should not have found 34 in the sample slice")
	}
}

func TestStringIn(t *testing.T) {
	var response bool
	sample := []string{"aa", "bb", "cc"}
	response = StringIn("bb", sample)
	if !response {
		t.Errorf("The StringIn function should have found 'bb' in the sample slice")
	}
	response = StringIn("xx", sample)
	if response {
		t.Errorf("The StringIn function should have found 'xx' in the sample slice")
	}
}
