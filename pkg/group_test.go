package la

import (
	"math"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func NewTestGroup() Group {
	a1 := NewTestUnit()
	a2 := NewTestUnit()
	return NewGroup(a1, a2)
}

func TestNewAssetVarsGroup(t *testing.T) {
	ag0 := NewGroup()
	assert.Equal(t, ag0.units, []Unit{}, "empty group does not return empty units slice")

	a1 := NewTestUnit()
	ag1 := NewGroup(a1)
	assert.Equal(t, ag1.units, []Unit{a1}, "group does not contain unit assigned in new")

	a2 := NewTestUnit()
	ag2 := NewGroup(a1, a2)
	assert.Equal(t, ag2.units, []Unit{a1, a2}, "group does not contain multiple units assigned in new")
}

func TestGetGroupDecisionVariableCoefficients(t *testing.T) {

	a1 := NewTestUnit()
	a2 := NewTestUnit()
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
	assert.Equal(t, ag.ColumnSize(), 8)
}

func TestNewGroupConstraint(t *testing.T) {
	pid1, _ := uuid.NewUUID()
	pid2, _ := uuid.NewUUID()
	inf := math.Inf(1)
	a1 := NewUnit(pid1, 1.0, 2.0, 3.0, 4.0, inf, inf, inf, inf)
	a2 := NewUnit(pid2, 5.0, 6.0, 7.0, 8.0, inf, inf, inf, inf)
	ag1 := NewGroup(a1, a2)

	c := []float64{0.0, -1.0, 0, 1.0, 0, 1.0, 0.2, 2.1, 1.1, 1.0}
	err := ag1.NewConstraint(c)
	assert.Nil(t, err)

	assert.Equal(t, ag1.Constraints()[0], c)
}

func TestGroupUnitConstraints(t *testing.T) {
	pid1, _ := uuid.NewUUID()
	pid2, _ := uuid.NewUUID()
	inf := math.Inf(1)
	a1 := NewUnit(pid1, 1.0, 2.0, 3.0, 4.0, inf, inf, inf, inf)
	a1c := []float64{-10, 1.0, 0.0, 1.0, 0.0, 10}
	a1.NewConstraint(a1c)

	a2 := NewUnit(pid2, 5.0, 6.0, 7.0, 8.0, inf, inf, inf, inf)
	a2c := []float64{-20, 0.0, 2.0, 0.0, 2.0, 20}
	a2.NewConstraint(a2c)

	ag1 := NewGroup(a1, a2)
	ag1c := []float64{0.0, -1.0, 0, 1.0, 0, 1.0, 0.2, 2.1, 1.1, 1.0}
	err := ag1.NewConstraint(ag1c)
	assert.Nil(t, err)

	assert.Equal(t, []float64{-10, 1.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 10}, ag1.Constraints()[0])
	assert.Equal(t, []float64{-20, 0.0, 0.0, 0.0, 0.0, 0.0, 2.0, 0.0, 2.0, 20}, ag1.Constraints()[1])
	assert.Equal(t, []float64{0, -1, 0, 1, 0, 1, .2, 2.1, 1.1, 1}, ag1.Constraints()[2])
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
	inf := math.Inf(1)
	assert.Equal(t, [][2]float64{
		{0, inf},
		{0, inf},
		{0, inf},
		{0, inf},
		{0, inf},
		{0, inf},
		{0, inf},
		{0, inf}},
		bounds)
}

func TestLocateGroupRealPositivePower(t *testing.T) {
	ag1 := NewTestGroup()
	loc := ag1.RealPositivePowerLoc()
	assert.Equal(t, []int{0, 4}, loc)
}

func TestLocateGroupRealNegativePower(t *testing.T) {
	ag1 := NewTestGroup()
	loc := ag1.RealNegativePowerLoc()
	assert.Equal(t, []int{1, 5}, loc)
}

func TestLocateRealCapacity(t *testing.T) {
	ag1 := NewTestGroup()
	loc := ag1.RealCapacityLoc()
	assert.Equal(t, []int{2, 6}, loc)
}

func TestLocateStoredCapacity(t *testing.T) {
	ag1 := NewTestGroup()
	loc := ag1.StoredEnergyLoc()
	assert.Equal(t, []int{3, 7}, loc)
}

// constraint generation
func TestNetloadConstraint(t *testing.T) {
	ag1 := NewTestGroup()
	nl := 11.1
	nlc := NetLoadConstraint(&ag1, nl)
	assert.Equal(t, []float64{nl, 1, -1, 0, 0, 1, -1, 0, 0, nl}, nlc)
}

func TestPositiveCapacityConstraint(t *testing.T) {
	ag1 := NewTestGroup()
	pc := 22.2
	pcc := GroupPositiveCapacityConstraint(&ag1, pc)
	inf := math.Inf(1)
	assert.Equal(t, []float64{pc, 0, 0, 1, 0, 0, 0, 1, 0, inf}, pcc)
}
