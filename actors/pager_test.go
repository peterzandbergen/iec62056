package actors

import "testing"

func TestPager(t *testing.T) {
	length := 10
	i := bound(100, 0, length)
	if i != length {
		t.Errorf("expected %d, got %d", length, i)
	}
}
