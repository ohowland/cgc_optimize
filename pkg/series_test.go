package gocoinor

import (
	"math"
	"math/rand"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func NewTestSeries(n int) Series {
	cl1 := NewTestCluster()
	clx := make([]Sequencer, n)
	for i := 0; i < n; i++ {
		clx[i] = cl1
	}

	return NewSeries(clx...)
}

func TestSeriesCostCoefficients(t *testing.T) {
	pid1, _ := uuid.NewUUID()
	pid2, _ := uuid.NewUUID()
	inf := math.Inf(1)
	a1 := NewUnit(pid1, 1, 2, 3, 4, inf, inf, inf, inf)
	a2 := NewUnit(pid2, 5, 6, 7, 8, inf, inf, inf, inf)
	g := NewGroup(a1, a2)
	cl := NewCluster(g)
	s := NewSeries(cl, cl)

	cc := s.CostCoefficients()
	assert.Equal(t, []float64{1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8}, cc)
}
func TestSeriesConstraints1(t *testing.T) {
	s := NewTestSeries(24)

	c := make([]float64, s.ColumnSize()+2)
	for i := 1; i < s.ColumnSize()-1; i++ {
		c[i] = rand.Float64()
	}
	err := s.NewConstraint(c)
	assert.Nil(t, err)

	sc := s.Constraints()
	assert.Equal(t, c, sc[0])

}
func TestSeriesBounds(t *testing.T) {
	s := NewTestSeries(2)
	b := s.Bounds()

	inf := math.Inf(1)
	expBnds := make([][2]float64, s.ColumnSize())
	for i := 0; i < s.ColumnSize(); i++ {
		expBnds[i] = [2]float64{0, inf}
	}
	assert.Equal(t, expBnds, b)
}

func TestSeriesBatteryEnergyConstraint(t *testing.T) {
	pid1, _ := uuid.NewUUID()
	pid2, _ := uuid.NewUUID()
	inf := math.Inf(1)
	a1 := NewUnit(pid1, 1, 2, 3, 4, inf, inf, inf, inf)
	a2 := NewUnit(pid2, 5, 6, 7, 8, inf, inf, inf, inf)
	g := NewGroup(a1, a2)
	cl := NewCluster(g)
	s := NewSeries(cl, cl, cl)
	bec := BatteryEnergyConstraint(&s, pid1, 1)
	for _, c := range bec {
		//fmt.Println(bec)
		err := s.NewConstraint(c)
		assert.Nil(t, err)
	}

	sec := s.Constraints()
	assert.Equal(t, []float64{0, -1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, -1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, sec[0])
	assert.Equal(t, []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, -1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, -1, 0, 0, 0, 0, 0}, sec[1])
}
