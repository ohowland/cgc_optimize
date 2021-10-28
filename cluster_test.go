package cgc_optimize

import (
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
			assert.Equal(t, u.CostCoefficients()[0], cc[i])
			assert.Equal(t, u.CostCoefficients()[1], cc[i+1])
			assert.Equal(t, u.CostCoefficients()[2], cc[i+2])
			assert.Equal(t, u.CostCoefficients()[3], cc[i+3])
			i += u.ColumnSize()
		}
	}
}

func TestClusterConstraints(t *testing.T) {
	cl := NewTestCluster()
	c := func() []float64 {
		r := make([]float64, 0)
		for i := 0; i < cl.ColumnSize()+2; i++ {
			r = append(r, float64(i))
		}
		return r
	}
	err := cl.NewConstraint(c())
	assert.Nil(t, err)

	clc := cl.Constraints()
	assert.Equal(t, c(), clc[len(clc)-1])
}

func TestClusterConstraintsFromGroups(t *testing.T) {
	ag1 := NewTestGroup()
	c1 := []float64{10, -10, 0, 10, 0, 0, 0, 0, -10, 0, 10, 0, 0, 0, 0, 10}
	err := ag1.NewConstraint(c1)
	assert.Nil(t, err)

	ag2 := NewTestGroup()
	c2 := []float64{20, -20, 0, 20, 0, 0, 0, 0, -20, 0, 20, 0, 0, 0, 0, 20}
	err = ag2.NewConstraint(c2)
	assert.Nil(t, err)

	cl := NewCluster([]Group{ag1, ag2}...)
	clc := cl.Constraints()

	exp_clc0 := append(append(append([]float64{lb(c1)}, cons(c1)...), make([]float64, len(c2)-2)...), ub(c1))
	assert.Equal(t, exp_clc0, clc[len(clc)/2-1])

	exp_clc1 := append(append(append([]float64{lb(c2)}, make([]float64, len(c1)-2)...), cons(c2)...), ub(c2))
	assert.Equal(t, exp_clc1, clc[len(clc)-1])
}

func TestClusterConstraintsFromAssets(t *testing.T) {
	pid1, _ := uuid.NewUUID()

	a1 := NewBasicUnit(
		pid1,
		[]CriticalPoint{NewCriticalPoint(-10, -1), NewCriticalPoint(0, 0), NewCriticalPoint(10, 1)},
		NewCriticalPoint(10, 1),
		NewCriticalPoint(10, 1))
	c1 := []float64{5, -10, 0, 10, 0, 0, 0, 0, 5}
	err := a1.NewConstraint(c1)
	assert.Nil(t, err)

	pid2, _ := uuid.NewUUID()
	a2 := NewBasicUnit(
		pid2,
		[]CriticalPoint{NewCriticalPoint(-10, -1), NewCriticalPoint(0, 0), NewCriticalPoint(10, 1)},
		NewCriticalPoint(10, 1),
		NewCriticalPoint(10, 1))
	c2 := []float64{10, -20, 0, 20, 0, 0, 0, 0, 10}
	err = a2.NewConstraint(c2)
	assert.Nil(t, err)

	ag1 := NewGroup(a1, a2)
	ag2 := NewGroup(a1)

	cl := NewCluster([]Group{ag1, ag2}...)
	clc := cl.Constraints()

	exp_clc0 := append(append(append(append([]float64{lb(c1)}, cons(c1)...), make([]float64, len(c2)-2)...), make([]float64, len(c1)-2)...), ub(c1))
	assert.Equal(t, exp_clc0, clc[len(clc)/3-1])

	exp_clc1 := append(append(append(append([]float64{lb(c2)}, make([]float64, len(c1)-2)...), cons(c2)...), make([]float64, len(c1)-2)...), ub(c2))
	assert.Equal(t, exp_clc1, clc[len(clc)*2/3-1])

	exp_clc2 := append(append(append(append([]float64{lb(c1)}, make([]float64, len(c1)-2)...), make([]float64, len(c2)-2)...), cons(c1)...), ub(c1))
	assert.Equal(t, exp_clc2, clc[len(clc)-1])
}
func TestClusterBounds(t *testing.T) {
	cl := NewTestCluster()
	b := cl.Bounds()

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
		{0, 1},
		{0, 1},
	}, b)
}

// Cluster Constraints

/*
func TestLinkbusConstraint(t *testing.T) {
	pid1, _ := uuid.NewUUID()
	a1 := NewBasicUnit(pid1, []CriticalPoint{NewCriticalPoint(-10, -1), NewCriticalPoint(0, 0), NewCriticalPoint(10, 1)})

	ag1 := NewGroup(a1)
	ag2 := NewGroup(a1)

	cl := NewCluster([]Group{ag1, ag2}...)

	lbc := LinkedBusConstraints(&cl, pid1)

	assert.Equal(t, []float64{0, 1, 0, 0, 0, -1, 0, 0, 0, 0}, lbc[0])
	assert.Equal(t, []float64{0, 0, 1, 0, 0, 0, -1, 0, 0, 0}, lbc[1])
}
*/
