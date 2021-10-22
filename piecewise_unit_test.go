package cgc_optimize

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func NewTestPiecewiseUnit() PiecewiseUnit {
	pid, _ := uuid.NewUUID()
	cp := []CriticalPoint{NewCriticalPoint(-10, 1), NewCriticalPoint(0, 0), NewCriticalPoint(10, 1)}
	return NewPiecewiseUnit(pid, cp)
}

func TestNewPiecewiseUnit(t *testing.T) {
	pu := NewTestPiecewiseUnit()

	assert.Equal(t, 5, pu.ColumnSize())

}
