package la

import (
	"math"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func NewTestCluster() Cluster {
	ag1 := NewTestGroup()
	ag2 := NewTestGroup()

	return NewCluster([]Group{ag1, ag2}...)
}

func TestClusterCostCoefficients(t *testing.T) {
	cl := NewTestCluster()

	cc := cl.CostCoefficients()

	var i int
	for _, g := range cl.groups {
		for _, u := range g.units {
			assert.Equal(t, u.coefficients[0], cc[i])
			assert.Equal(t, u.coefficients[1], cc[i+1])
			assert.Equal(t, u.coefficients[2], cc[i+2])
			assert.Equal(t, u.coefficients[3], cc[i+3])
			i += u.ColumnSize()
		}
	}
}

func TestClusterConstraints1(t *testing.T) {
	cl := NewTestCluster()
	c := []float64{-1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	err := cl.NewConstraint(c)
	assert.Nil(t, err)

	clc := cl.Constraints()
	assert.Equal(t, c, clc[0])
}

func TestClusterConstraints2(t *testing.T) {
	ag1 := NewTestGroup()
	c1 := []float64{0, 1, 0, 1, 0, 1, 0, 1, 0, 10}
	err := ag1.NewConstraint(c1)
	assert.Nil(t, err)

	ag2 := NewTestGroup()
	c2 := []float64{0, 0, 2, 0, 2, 0, 2, 0, 2, 10}
	err = ag2.NewConstraint(c2)
	assert.Nil(t, err)

	cl := NewCluster([]Group{ag1, ag2}...)
	clc := cl.Constraints()

	exp_clc0 := append(append(append([]float64{lb(c1)}, cons(c1)...), make([]float64, len(c2)-2)...), ub(c1))
	assert.Equal(t, exp_clc0, clc[0])

	exp_clc1 := append(append(append([]float64{lb(c2)}, make([]float64, len(c1)-2)...), cons(c2)...), ub(c2))
	assert.Equal(t, exp_clc1, clc[1])
}

func TestClusterConstraints3(t *testing.T) {
	pid1, _ := uuid.NewUUID()
	inf := math.Inf(1)
	a1 := NewUnit(pid1, 1.0, 2.0, 3.0, 4.0, inf, inf, inf, inf)
	c1 := []float64{-1, 1, 0, 1, 0, 10}
	a1.NewConstraint(c1)

	pid2, _ := uuid.NewUUID()
	a2 := NewUnit(pid2, 5.0, 6.0, 7.0, 8.0, inf, inf, inf, inf)
	c2 := []float64{-2, 0, 2, 0, 2, 20}
	a2.NewConstraint(c2)

	ag1 := NewGroup(a1, a2)
	ag2 := NewGroup(a1)

	cl := NewCluster([]Group{ag1, ag2}...)
	clc := cl.Constraints()

	exp_clc0 := append(append(append(append([]float64{lb(c1)}, cons(c1)...), make([]float64, len(c2)-2)...), make([]float64, len(c1)-2)...), ub(c1))
	assert.Equal(t, exp_clc0, clc[0])

	exp_clc1 := append(append(append(append([]float64{lb(c2)}, make([]float64, len(c1)-2)...), cons(c2)...), make([]float64, len(c1)-2)...), ub(c2))
	assert.Equal(t, exp_clc1, clc[1])

	exp_clc2 := append(append(append(append([]float64{lb(c1)}, make([]float64, len(c1)-2)...), make([]float64, len(c2)-2)...), cons(c1)...), ub(c1))
	assert.Equal(t, exp_clc2, clc[2])
}
func TestClusterBounds(t *testing.T) {
	cl := NewTestCluster()
	b := cl.Bounds()

	inf := math.Inf(1)
	assert.Equal(t, [][2]float64{
		{0, inf},
		{0, inf},
		{0, inf},
		{0, inf},
		{0, inf},
		{0, inf},
		{0, inf},
		{0, inf},
		{0, inf},
		{0, inf},
		{0, inf},
		{0, inf},
		{0, inf},
		{0, inf},
		{0, inf},
		{0, inf},
	}, b)
}

// Cluster Constraints

func TestLinkbusConstraint(t *testing.T) {
	pid1, _ := uuid.NewUUID()
	inf := math.Inf(1)
	a1 := NewUnit(pid1, 1.0, 2.0, 3.0, 4.0, inf, inf, inf, inf)

	ag1 := NewGroup(a1)
	ag2 := NewGroup(a1)

	cl := NewCluster([]Group{ag1, ag2}...)

	lbc := LinkedBusConstraints(&cl, pid1)

	assert.Equal(t, []float64{0, 1, 0, 0, 0, -1, 0, 0, 0, 0}, lbc[0])
	assert.Equal(t, []float64{0, 0, 1, 0, 0, 0, -1, 0, 0, 0}, lbc[1])
}
