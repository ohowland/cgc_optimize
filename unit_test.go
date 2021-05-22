package cgc_optimize

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
	inf := math.Inf(1)
	return NewUnit(pid, rand.Float64(), rand.Float64(), rand.ExpFloat64(), rand.Float64(), inf, inf, inf, inf)
}

func TestNewAssetVar(t *testing.T) {
	pid, _ := uuid.NewUUID()
	Cp := 0.2
	Cn := -0.1
	Cc := 1.0
	Ce := 0.0
	inf := math.Inf(1)

	a := NewUnit(pid, Cp, Cn, Cc, Ce, inf, inf, inf, inf)

	assert.Equal(t, a.coefficients[0], Cp)
	assert.Equal(t, a.coefficients[1], Cn)
	assert.Equal(t, a.coefficients[2], Cc)
	assert.Equal(t, a.coefficients[3], Ce)
}

func TestGetCostCoefficients(t *testing.T) {
	a := NewTestUnit()
	assert.Equal(t,
		[]float64{a.coefficients[0],
			a.coefficients[1],
			a.coefficients[2],
			a.coefficients[3]},
		a.CostCoefficients())
}

func TestGetUnitColumnSize(t *testing.T) {
	a1 := NewTestUnit()
	assert.Equal(t, a1.ColumnSize(), 4)
}

func TestApplyUnitConstraint(t *testing.T) {
	a1 := NewTestUnit()
	c := []float64{0, 2, 0, 1, 0, 1}

	err := a1.NewConstraint(c)
	assert.Nil(t, err)

	assert.Equal(t, a1.Constraints()[0], []float64{0, 2, 0, 1, 0, 1})
}

func TestBadConstraint(t *testing.T) {
	a1 := NewTestUnit()
	c := []float64{0, 2, 0, 1, 0.3, 0.1, 1}

	err := a1.NewConstraint(c)
	assert.Error(t, err)
}

func TestReturnConstraints(t *testing.T) {
	a1 := NewTestUnit()
	c1 := []float64{0.0, 0.0, 1.0, 2.0, 3.0, 1.0}
	c2 := []float64{0.0, 1.1, 2.2, 3.3, 4.4, 1.0}
	c3 := []float64{0.0, -0.3, 0.2, -0.1, 0.0, 1.0}

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
			assert.Equal(t, []float64{0.0, 0.0, 1.0, 2.0, 3.0, 1.0}, c)
		case 1:
			assert.Equal(t, []float64{0.0, 1.1, 2.2, 3.3, 4.4, 1.0}, c)
		case 2:
			assert.Equal(t, []float64{0.0, -0.3, 0.2, -0.1, 0.0, 1.0}, c)
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

// Constraint tests

func TestUnitCapacityConstraint(t *testing.T) {
	pid, _ := uuid.NewUUID()
	a := NewUnit(pid, 1, 1, 0.1, 0, 10, 10, 10, 0)

	ucc := UnitCapacityConstraints(&a)
	for _, c := range ucc {
		err := a.NewConstraint(c)
		assert.Nil(t, err)
	}

	ac := a.Constraints()
	assert.Equal(t, []float64{0, -1, 0, 1, 0, math.Inf(1)}, ac[0], "positive capacity constraint malformed")
	assert.Equal(t, []float64{0, 0, -1, 1, 0, math.Inf(1)}, ac[1], "negative capacity constraint malformed")
	assert.Len(t, ac, 2)
}
