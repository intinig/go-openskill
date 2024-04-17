package ptr_test

import (
	"testing"

	"github.com/intinig/go-openskill/ptr"
)

func TestFloat64(t *testing.T) {
	t.Parallel()
	x := 64.0
	y := ptr.Float64(x)
	if *y != x {
		t.Errorf("Expected %f, got %f", x, *y)
	}
}

func TestInt(t *testing.T) {
	t.Parallel()
	x := 64
	y := ptr.Int(x)
	if *y != x {
		t.Errorf("Expected %d, got %d", x, *y)
	}
}
