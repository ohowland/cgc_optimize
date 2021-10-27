package cgc_optimize

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func NewTestBasicUnit() BasicUnit {
	pid, _ := uuid.NewUUID()
	cp := []CriticalPoint{NewCriticalPoint(-10, 1), NewCriticalPoint(0, 0), NewCriticalPoint(10, 1)}
	return NewBasicUnit(pid, cp)
}

func TestNewBasicUnit(t *testing.T) {
	pu := NewTestBasicUnit()

	assert.Equal(t, 5, pu.ColumnSize())

}
