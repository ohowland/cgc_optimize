package cgc_optimize

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestNewPiecewiseUnit(t *testing.T) {

	pid, _ := uuid.NewUUID()

	cp := []CriticalPoint{CriticalPoint{-5, 0.8}, CriticalPoint{0, 0}, CriticalPoint{5, 1.2}}
	pu := NewPiecewiseUnit(pid, cp)

	fmt.Println(pu)
}
