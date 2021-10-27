package cgc_optimize

import (
	"math"
	"sort"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func NewTestBasicUnit() BasicUnit {
	pid, _ := uuid.NewUUID()
	cp := []CriticalPoint{NewCriticalPoint(-10, -1), NewCriticalPoint(0, 0), NewCriticalPoint(10, 1)}
	return NewBasicUnit(pid, cp)
}

func TestNewBasicUnitSize(t *testing.T) {
	u := NewTestBasicUnit()
	assert.Equal(t, 7, u.ColumnSize())
}

func TestBasicUnitLocs(t *testing.T) {
	u := NewTestBasicUnit()
	assert.Equal(t, []int{0, 1, 2}, u.RealPowerLoc())
	assert.Equal(t, []int{5}, u.RealPositiveCapacityLoc())
	assert.Equal(t, []int{6}, u.RealNegativeCapacityLoc())
}

func TestBasicUnitIntegrality(t *testing.T) {
	u := NewTestBasicUnit()

	assert.Equal(t, []int{0, 0, 0, 1, 1, 0, 0}, u.Integrality())
}

func TestBasicUnitCriticalPoints(t *testing.T) {
	pid, _ := uuid.NewUUID()
	cp := []CriticalPoint{NewCriticalPoint(-10, -1), NewCriticalPoint(0, 0), NewCriticalPoint(10, 1)}
	u := NewBasicUnit(pid, cp)
	assert.Equal(t, cp, u.CriticalPoints())

	rev_cp := make([]CriticalPoint, len(cp))
	copy(rev_cp, cp)
	sort.Slice(rev_cp, func(i, j int) bool { return rev_cp[i].KW() > rev_cp[j].KW() })

	assert.NotEqual(t, rev_cp, u.CriticalPoints())
}

func TestBasicUnitRealPowerConstraint(t *testing.T) {
	u := NewTestBasicUnit()

	cn := UnitRealPowerConstraint(&u, 5.5)

	assert.Equal(t, []float64{5.5, -10, 0, 10, 0, 0, 0, 0, 5.5}, cn)
}

func TestBasicUnitPositiveCapacityConstraint(t *testing.T) {
	u := NewTestBasicUnit()

	cn := UnitPositiveCapacityConstraint(&u, 10)

	assert.Equal(t, []float64{math.Inf(-1), -10, 0, 10, 0, 0, -10, 0, 0}, cn)
}

func TestBasicUnitNegativeCapacityConstraint(t *testing.T) {
	u := NewTestBasicUnit()

	cn := UnitNegativeCapacityConstraint(&u, 10)

	assert.Equal(t, []float64{0, -10, 0, 10, 0, 0, 0, 10, math.Inf(1)}, cn)
}
