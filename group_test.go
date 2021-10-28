package cgc_optimize

import (
	"math"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func NewTestGroup() Group {
	a1 := NewTestBasicUnit()
	a2 := NewTestBasicUnit()
	return NewGroup(a1, a2)
}

func TestNewAssetVarsGroup(t *testing.T) {
	ag0 := NewGroup()
	assert.Equal(t, ag0.units, []Unit{}, "empty group does not return empty units slice")

	a1 := NewTestBasicUnit()
	ag1 := NewGroup(a1)
	assert.Equal(t, ag1.units, []Unit{a1}, "group does not contain unit assigned in new")

	a2 := NewTestBasicUnit()
	ag2 := NewGroup(a1, a2)
	assert.Equal(t, ag2.units, []Unit{a1, a2}, "group does not contain multiple units assigned in new")
}

func TestGetGroupDecisionVariableCoefficients(t *testing.T) {

	a1 := NewTestBasicUnit()
	a2 := NewTestBasicUnit()
	ag2 := NewGroup(a1, a2)

	cc1 := a1.CostCoefficients()
	cc2 := a2.CostCoefficients()
	cc := append(cc1, cc2...)
	assert.Equal(t, cc, ag2.CostCoefficients(), "group incorrectly formulates decision variable coefficents from internal units")

	ag0 := NewGroup()
	assert.Equal(t, []float64{}, ag0.CostCoefficients(), "group incorrectly formulates cost coefficents from empty unit slice")

}

func TestGetGroupColumnsSize(t *testing.T) {
	ag := NewTestGroup()
	assert.Equal(t, ag.ColumnSize(), 14)
}

func TestNewGroupConstraint(t *testing.T) {
	pid1, _ := uuid.NewUUID()
	pid2, _ := uuid.NewUUID()

	a1 := NewBasicUnit(
		pid1,
		[]CriticalPoint{NewCriticalPoint(-10, -1), NewCriticalPoint(0, 0), NewCriticalPoint(10, 1)},
		NewCriticalPoint(10, 1),
		NewCriticalPoint(10, 1))
	a2 := NewBasicUnit(
		pid2,
		[]CriticalPoint{NewCriticalPoint(0, 0), NewCriticalPoint(5, 1), NewCriticalPoint(15, 5)},
		NewCriticalPoint(15, 1),
		NewCriticalPoint(0, 0))
	ag1 := NewGroup(a1, a2)

	c := []float64{0, -1, 0, 1, 0, 0, 0, 0, 0, 1, 5, 0, 0, 0, 0, 0}
	err := ag1.NewConstraint(c)
	assert.Nil(t, err)

	cns := ag1.Constraints()
	assert.Equal(t, c, cns[len(cns)-1])
}

func TestGroupUnitConstraints(t *testing.T) {
	pid1, _ := uuid.NewUUID()
	pid2, _ := uuid.NewUUID()

	a1 := NewBasicUnit(
		pid1,
		[]CriticalPoint{NewCriticalPoint(-10, -1), NewCriticalPoint(0, 0), NewCriticalPoint(10, 1)},
		NewCriticalPoint(10, 1),
		NewCriticalPoint(10, 1))
	a1c := []float64{0, -10, 0, 10, 0, 0, 0, 10, 0}
	a1.NewConstraint(a1c)

	a2 := NewBasicUnit(
		pid2,
		[]CriticalPoint{NewCriticalPoint(0, 0), NewCriticalPoint(5, 1), NewCriticalPoint(15, 5)},
		NewCriticalPoint(15, 1),
		NewCriticalPoint(0, 0))
	a2c := []float64{0, 0, 5, 15, 0, 0, 0, 15, 0}
	a2.NewConstraint(a2c)

	ag1 := NewGroup(a1, a2)
	ag1c := []float64{0, -10, 0, 10, 0, 0, 0, 0, 0, 5, 15, 0, 0, 0, 0, 0}
	err := ag1.NewConstraint(ag1c)
	assert.Nil(t, err)

	cons := ag1.Constraints()
	assert.Equal(t, []float64{0, -10, 0, 10, 0, 0, 0, 10, 0, 0, 0, 0, 0, 0, 0, 0}, cons[len(cons)/2-1])
	assert.Equal(t, []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 15, 0, 0, 0, 15, 0}, cons[(len(cons)-2)])
	assert.Equal(t, ag1c, cons[(len(cons)-1)])
}

func TestBadGroupConstraint(t *testing.T) {
	ag1 := NewTestGroup()
	c := []float64{-1.0, 1.1}

	err := ag1.NewConstraint(c)
	assert.Error(t, err)
}

func TestGroupBounds(t *testing.T) {
	ag1 := NewTestGroup()
	bounds := ag1.Bounds()
	assert.Equal(t, [][2]float64{
		{0, 1},
		{0, 1},
		{0, 1},
		{0, 1},
		{0, 1},
		{0, 1},
		{0, 1},
		{0, 1},
		{0, 1},
		{0, 1},
		{0, 1},
		{0, 1},
		{0, 1},
		{0, 1}},
		bounds)
}

func TestLocateGroupRealPositivePower(t *testing.T) {
	ag1 := NewTestGroup()
	loc := ag1.RealPowerLoc()
	assert.Equal(t, []int{0, 1, 2, 7, 8, 9}, loc)
}

func TestLocateRealPositiveCapacity(t *testing.T) {
	ag1 := NewTestGroup()
	loc := ag1.RealPositiveCapacityLoc()
	assert.Equal(t, []int{5, 12}, loc)
}

func TestLocateRealNegativeCapacity(t *testing.T) {
	ag1 := NewTestGroup()
	loc := ag1.RealNegativeCapacityLoc()
	assert.Equal(t, []int{6, 13}, loc)
}

func TestLocateStoredCapacity(t *testing.T) {
	ag1 := NewTestGroup()
	loc := ag1.StoredEnergyLoc()
	assert.Equal(t, []int{}, loc)
}

// constraint generation
func TestGroupNetloadConstraint(t *testing.T) {
	pid1, _ := uuid.NewUUID()
	pid2, _ := uuid.NewUUID()

	a1 := NewBasicUnit(
		pid1,
		[]CriticalPoint{NewCriticalPoint(-10, -1), NewCriticalPoint(0, 0), NewCriticalPoint(10, 1)},
		NewCriticalPoint(10, 1),
		NewCriticalPoint(10, 1))
	a1c := []float64{0, -10, 0, 10, 0, 0, 0, 10, 0}
	a1.NewConstraint(a1c)

	a2 := NewBasicUnit(
		pid2,
		[]CriticalPoint{NewCriticalPoint(0, 0), NewCriticalPoint(5, 1), NewCriticalPoint(15, 5)},
		NewCriticalPoint(15, 1),
		NewCriticalPoint(0, 0))
	a2c := []float64{0, 0, 5, 15, 0, 0, 0, 15, 0}
	a2.NewConstraint(a2c)

	ag1 := NewGroup(a1, a2)
	nl := 11.1
	nlc := NetLoadConstraint(&ag1, nl)
	assert.Equal(t, []float64{nl, -10, 0, 10, 0, 0, 0, 0, 0, 5, 15, 0, 0, 0, 0, nl}, nlc)
}

func TestGroupPositiveCapacityConstraint(t *testing.T) {
	ag1 := NewTestGroup()
	pc := 22.2
	pcc := GroupPositiveCapacityConstraint(&ag1, pc)
	inf := math.Inf(1)

	assert.Equal(t, []float64{pc, 0, 0, 0, 0, 0, 10, 0, 0, 0, 0, 0, 0, 10, 0, inf}, pcc)
}

func TestGroupNegativeCapacityConstraint(t *testing.T) {
	ag1 := NewTestGroup()
	nc := 14.5
	ncc := GroupNegativeCapacityConstraint(&ag1, nc)
	inf := math.Inf(1)

	assert.Equal(t, []float64{nc, 0, 0, 0, 0, 0, 0, 10, 0, 0, 0, 0, 0, 0, 10, inf}, ncc)
}
