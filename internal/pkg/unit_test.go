package esslp

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func NewTestUnit() Unit {
	pid, _ := uuid.NewUUID()
	rand.Seed(time.Now().Unix())
	return NewUnit(pid, rand.Float64(), rand.Float64(), rand.ExpFloat64(), rand.Float64())
}

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
	a := NewTestUnit()
	assert.Equal(t, []float64{a.Cp, a.Cn, a.Cc, a.Ce}, a.CostCoefficients())
}

func TestGetUnitColumnSize(t *testing.T) {
	a1 := NewTestUnit()
	assert.Equal(t, a1.ColumnSize(), 4)
}

func TestApplyUnitConstraint(t *testing.T) {
	a1 := NewTestUnit()
	c := []float64{2, 0, 1, 0}
	lb := 0.0
	ub := 1.0

	err := a1.NewConstraint(c, lb, ub)
	assert.Nil(t, err)

	assert.Equal(t, a1.Constraints()[0], []float64{lb, 2, 0, 1, 0, ub})
}

func TestBadConstraint(t *testing.T) {
	a1 := NewTestUnit()
	c := []float64{2, 0, 1, 0.3, 0.1}
	lb := 0.0
	ub := 1.0

	err := a1.NewConstraint(c, lb, ub)
	assert.Error(t, err)
}

func TestReturnConstraints(t *testing.T) {
	a1 := NewTestUnit()
	c1 := []float64{0.0, 1.0, 2.0, 3.0}
	c2 := []float64{1.1, 2.2, 3.3, 4.4}
	c3 := []float64{-0.3, 0.2, -0.1, 0.0}
	lb := 0.0
	ub := 1.0

	err := a1.NewConstraint(c1, lb, ub)
	assert.Nil(t, err)
	err = a1.NewConstraint(c2, lb, ub)
	assert.Nil(t, err)
	err = a1.NewConstraint(c3, lb, ub)
	assert.Nil(t, err)

	cx := a1.Constraints()

	for i, c := range cx {
		switch i {
		case 0:
			assert.Equal(t, []float64{lb, 0.0, 1.0, 2.0, 3.0, ub}, c)
		case 1:
			assert.Equal(t, []float64{lb, 1.1, 2.2, 3.3, 4.4, ub}, c)
		case 2:
			assert.Equal(t, []float64{lb, -0.3, 0.2, -0.1, 0.0, ub}, c)
		default:
			assert.Fail(t, "unexpected constraint in unit: %v\n", c)
		}
	}
}

func TestUnitBounds(t *testing.T) {
	a1 := NewTestUnit()
	b := a1.Bounds()

	inf := math.Inf(1)
	assert.Equal(t, [][2]float64{{0, inf}, {0, inf}, {0, inf}, {0, inf}}, b)
}
