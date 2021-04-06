package esslp

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
			assert.Equal(t, u.Cp, cc[i])
			assert.Equal(t, u.Cn, cc[i+1])
			assert.Equal(t, u.Cc, cc[i+2])
			assert.Equal(t, u.Ce, cc[i+3])
			i += u.ColumnSize()
		}
	}
}

func TestClusterConstraints1(t *testing.T) {
	cl := NewTestCluster()
	c := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	lb := 0.0
	ub := 1.0
	err := cl.NewConstraint(c, lb, ub)
	assert.Nil(t, err)

	clc := cl.Constraints()
	assert.Equal(t, append(append([]float64{lb}, c...), ub), clc[0])
}

func TestClusterConstraints2(t *testing.T) {
	ag1 := NewTestGroup()
	c1 := []float64{1, 0, 1, 0, 1, 0, 1, 0}
	lb := 0.0
	ub := 10.0
	err := ag1.NewConstraint(c1, lb, ub)
	assert.Nil(t, err)

	ag2 := NewTestGroup()
	c2 := []float64{0, 2, 0, 2, 0, 2, 0, 2}
	err = ag2.NewConstraint(c2, lb, ub)
	assert.Nil(t, err)

	cl := NewCluster([]Group{ag1, ag2}...)
	clc := cl.Constraints()

	exp_clc0 := append(append(append([]float64{lb}, c1...), make([]float64, len(c2))...), ub)
	assert.Equal(t, exp_clc0, clc[0])

	exp_clc1 := append(append(append([]float64{lb}, make([]float64, len(c1))...), c2...), ub)
	assert.Equal(t, exp_clc1, clc[1])
}

func TestClusterConstraints3(t *testing.T) {
	pid1, _ := uuid.NewUUID()
	a1 := NewUnit(pid1, 1.0, 2.0, 3.0, 4.0)
	lb1 := 0.0
	ub1 := 2.0
	c1 := []float64{1, 0, 1, 0}
	a1.NewConstraint(c1, lb1, ub1)

	pid2, _ := uuid.NewUUID()
	a2 := NewUnit(pid2, 5.0, 6.0, 7.0, 8.0)
	lb2 := 1.0
	ub2 := 10.0
	c2 := []float64{0, 2, 0, 2}
	a2.NewConstraint(c2, lb2, ub2)

	ag1 := NewGroup(a1, a2)
	ag2 := NewGroup(a1)

	cl := NewCluster([]Group{ag1, ag2}...)
	clc := cl.Constraints()

	exp_clc0 := append(append(append(append([]float64{lb1}, c1...), make([]float64, len(c2))...), make([]float64, len(c1))...), ub1)
	assert.Equal(t, exp_clc0, clc[0])

	exp_clc1 := append(append(append(append([]float64{lb2}, make([]float64, len(c1))...), c2...), make([]float64, len(c1))...), ub2)
	assert.Equal(t, exp_clc1, clc[1])

	exp_clc2 := append(append(append(append([]float64{lb1}, make([]float64, len(c1))...), make([]float64, len(c2))...), c1...), ub1)
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
