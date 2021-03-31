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

func TestFlattenAssetVar(t *testing.T) {
	pid, _ := uuid.NewUUID()
	Cp := 1.0
	Cn := 2.0
	Cc := 3.0
	Ce := 4.0

	a := NewUnit(pid, Cp, Cn, Cc, Ce)

	assert.Equal(t, a.Flatten(), []float64{1.0, 2.0, 3.0, 4.0})
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

func TestNewAssetVarsGroup(t *testing.T) {

	ag0 := NewGroup()
	assert.Equal(t, ag0.Units, []Unit{}, "empty group does not return empty units slice")

	pid1, _ := uuid.NewUUID()
	a1 := NewUnit(pid1, 1.0, 2.0, 3.0, 4.0)
	ag1 := NewGroup(a1)
	assert.Equal(t, ag1.Units, []Unit{a1}, "group does not contain unit assigned in new")

	pid2, _ := uuid.NewUUID()
	a2 := NewUnit(pid2, 5.0, 6.0, 7.0, 8.0)
	ag2 := NewGroup(a1, a2)
	assert.Equal(t, ag2.Units, []Unit{a1, a2}, "group does not contain multiple units assigned in new")
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
