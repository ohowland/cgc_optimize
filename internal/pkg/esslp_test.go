package esslp

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewAssetVar(t *testing.T) {
	pid, _ := uuid.NewUUID()
	Cp := 0.2
	Cn := -0.1
	Cc := 1.0
	Ce := 0.0

	a := NewUnit(pid, Cp, Cn, Cc, Ce)

	assert.Equal(t, a.Cp, Cp)
	assert.Equal(t, a.Cn, Cn)
	assert.Equal(t, a.Cc, Cc)
	assert.Equal(t, a.Ce, Ce)
}

func TestGetCostCoefficients(t *testing.T) {
	pid, _ := uuid.NewUUID()
	Cp := 1.0
	Cn := 2.0
	Cc := 3.0
	Ce := 4.0

	a := NewUnit(pid, Cp, Cn, Cc, Ce)

	assert.Equal(t, a.CostCoefficients(), []float64{1.0, 2.0, 3.0, 4.0})
}

func TestGetUnitColumnSize(t *testing.T) {
	pid1, _ := uuid.NewUUID()
	a1 := NewUnit(pid1, 1.0, 2.0, 3.0, 4.0)
	assert.Equal(t, a1.ColumnSize(), 4)
}

func TestApplyUnitConstraint(t *testing.T) {
	pid1, _ := uuid.NewUUID()
	a1 := NewUnit(pid1, 1.0, 2.0, 3.0, 4.0)
	c := []float64{2.0, 0, 1.0, 0}

	err := a1.NewConstraint(c)
	assert.Nil(t, err)

	assert.Equal(t, a1.Constraints()[0], c)

}

func TestBadConstraint(t *testing.T) {
	pid1, _ := uuid.NewUUID()
	a1 := NewUnit(pid1, 1.0, 2.0, 3.0, 4.0)
	c := []float64{2.0, 0, 1.0, 0.3, 0.1}

	err := a1.NewConstraint(c)
	assert.Error(t, err)
}

func TestReturnConstraints(t *testing.T) {
	pid1, _ := uuid.NewUUID()
	a1 := NewUnit(pid1, 1.0, 2.0, 3.0, 4.0)
	c1 := []float64{0.0, 1.0, 2.0, 3.0}
	c2 := []float64{1.1, 2.2, 3.3, 4.4}
	c3 := []float64{-0.3, 0.2, -0.1, 0.0}

	err := a1.NewConstraint(c1)
	assert.Nil(t, err)
	err = a1.NewConstraint(c2)
	assert.Nil(t, err)
	err = a1.NewConstraint(c3)
	assert.Nil(t, err)

	cx := a1.Constraints()

	for i, c := range cx {
		switch i {
		case 0:
			assert.Equal(t, c, c1)
		case 1:
			assert.Equal(t, c, c2)
		case 2:
			assert.Equal(t, c, c3)
		default:
			assert.Fail(t, "unexpected constraint in unit: %v\n", c)
		}
	}
}

func TestUnitBounds(t *testing.T) {
	pid1, _ := uuid.NewUUID()
	a1 := NewUnit(pid1, 1.0, 2.0, 3.0, 4.0)

	b := a1.Bounds()
	assert.Fail(t, "bounds not specified for test %v", b)
}

// GROUP TESTS

func TestNewAssetVarsGroup(t *testing.T) {

	ag0 := NewGroup()
	assert.Equal(t, ag0.units, []Unit{}, "empty group does not return empty units slice")

	pid1, _ := uuid.NewUUID()
	a1 := NewUnit(pid1, 1.0, 2.0, 3.0, 4.0)
	ag1 := NewGroup(a1)
	assert.Equal(t, ag1.units, []Unit{a1}, "group does not contain unit assigned in new")

	pid2, _ := uuid.NewUUID()
	a2 := NewUnit(pid2, 5.0, 6.0, 7.0, 8.0)
	ag2 := NewGroup(a1, a2)
	assert.Equal(t, ag2.units, []Unit{a1, a2}, "group does not contain multiple units assigned in new")
}

func TestGetGroupDecisionVariableCoefficients(t *testing.T) {

	pid1, _ := uuid.NewUUID()
	pid2, _ := uuid.NewUUID()
	a1 := NewUnit(pid1, 1.0, 2.0, 3.0, 4.0)
	a2 := NewUnit(pid2, 5.0, 6.0, 7.0, 8.0)
	ag2 := NewGroup(a1, a2)

	assert.Equal(t, ag2.CostCoefficients(), []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0}, "group incorrectly formulates decision variable coefficents from internal units")

	ag0 := NewGroup()
	assert.Equal(t, ag0.CostCoefficients(), []float64{}, "group incorrectly formulates cost coefficents from empty unit slice")

}

func TestGetGroupColumnsSize(t *testing.T) {
	pid1, _ := uuid.NewUUID()
	pid2, _ := uuid.NewUUID()
	a1 := NewUnit(pid1, 1.0, 2.0, 3.0, 4.0)
	a2 := NewUnit(pid2, 5.0, 6.0, 7.0, 8.0)
	ag2 := NewGroup(a1, a2)

	assert.Equal(t, ag2.ColumnSize(), 8)
}

func TestNewGroupConstraint(t *testing.T) {
	pid1, _ := uuid.NewUUID()
	pid2, _ := uuid.NewUUID()
	a1 := NewUnit(pid1, 1.0, 2.0, 3.0, 4.0)
	a2 := NewUnit(pid2, 5.0, 6.0, 7.0, 8.0)
	ag1 := NewGroup(a1, a2)

	c := []float64{-1.0, 0, 1.0, 0, 1.0, 0.2, 2.1, 1.1}

	err := ag1.NewConstraint(c)
	assert.Nil(t, err)
	assert.Equal(t, ag1.Constraints()[0], c)
}

func TestBadGroupConstraint(t *testing.T) {
	pid1, _ := uuid.NewUUID()
	pid2, _ := uuid.NewUUID()
	a1 := NewUnit(pid1, 1.0, 2.0, 3.0, 4.0)
	a2 := NewUnit(pid2, 5.0, 6.0, 7.0, 8.0)
	ag1 := NewGroup(a1, a2)

	c := []float64{-1.0, 1.1}

	err := ag1.NewConstraint(c)
	assert.Error(t, err)
}

func TestGroupBounds(t *testing.T) {
	pid1, _ := uuid.NewUUID()
	pid2, _ := uuid.NewUUID()
	a1 := NewUnit(pid1, 1.0, 2.0, 3.0, 4.0)
	a2 := NewUnit(pid2, 5.0, 6.0, 7.0, 8.0)
	ag1 := NewGroup(a1, a2)

	bounds := ag1.Bounds()
	assert.Fail(t, "bounds %v", bounds)
}

// LP

func TestSimpleLPFormulation(t *testing.T) {
	pid1, _ := uuid.NewUUID()
	pid2, _ := uuid.NewUUID()
	a1 := NewUnit(pid1, 1.0, 2.0, 3.0, 4.0)
	a2 := NewUnit(pid2, 5.0, 6.0, 7.0, 8.0)
	ag1 := NewGroup(a1, a2)

	lp := NewLP(&ag1)

	sol := lp.Solve()
	assert.Fail(t, "sol[0] = %v", sol[0])
}
