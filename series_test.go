package cgc_optimize

import (
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
	a1 := NewBasicUnit(
		pid1,
		[]CriticalPoint{NewCriticalPoint(-10, -1), NewCriticalPoint(0, 0), NewCriticalPoint(10, 1)},
		NewCriticalPoint(10, 1.1),
		NewCriticalPoint(10, 2.2))
	a2 := NewBasicUnit(
		pid2,
		[]CriticalPoint{NewCriticalPoint(0, 0), NewCriticalPoint(5, 2), NewCriticalPoint(15, 6)},
		NewCriticalPoint(20, 3.3),
		NewCriticalPoint(0, 0))
	g := NewGroup(a1, a2)
	cl := NewCluster(g)
	s := NewSeries(cl, cl)

	cc := s.CostCoefficients()
	assert.Equal(t, []float64{-1, 0, 1, 0, 0, 1.1, 2.2, 0, 2, 6, 0, 0, 3.3, 0, -1, 0, 1, 0, 0, 1.1, 2.2, 0, 2, 6, 0, 0, 3.3, 0}, cc)
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
	assert.Equal(t, c, sc[len(sc)-1])

}
func TestSeriesBounds(t *testing.T) {
	s := NewTestSeries(2)
	b := s.Bounds()

	expBnds := make([][2]float64, s.ColumnSize())
	for i := 0; i < s.ColumnSize(); i++ {
		expBnds[i] = [2]float64{0, 1}
	}
	assert.Equal(t, expBnds, b)
}

func TestSeriesBatteryEnergyConstraint(t *testing.T) {
	pid1, _ := uuid.NewUUID()
	a1 := NewEssUnit(
		pid1,
		[]CriticalPoint{
			NewCriticalPoint(-10, -1),
			NewCriticalPoint(0, 0),
			NewCriticalPoint(10, 1)},
		NewCriticalPoint(10, 0.1),
		NewCriticalPoint(10, 0.1),
		100)

	g1 := NewGroup(a1)
	s := NewSeries(g1, g1, g1)

	bec := BatteryEnergyConstraint(&s, pid1, 1)
	for _, c := range bec {
		err := s.NewConstraint(c)
		assert.Nil(t, err)
	}

	sec := s.Constraints()
	assert.Equal(t, []float64{0, 10, 0, -10, 0, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, -100, 0, 0, 0, 0, 0, 0, 0, 0, 0}, sec[0])
	assert.Equal(t, []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 0, -10, 0, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, -100, 0}, sec[1])
}

func TestSeriesBatteryEnergyConstraintMultiUnit(t *testing.T) {
	pid1, _ := uuid.NewUUID()
	a1 := NewBasicUnit(
		pid1,
		[]CriticalPoint{
			NewCriticalPoint(0, 0),
			NewCriticalPoint(5, 1),
			NewCriticalPoint(20, 5)},
		NewCriticalPoint(20, 0.1),
		NewCriticalPoint(0, 0))

	pid2, _ := uuid.NewUUID()
	a2 := NewEssUnit(
		pid2,
		[]CriticalPoint{
			NewCriticalPoint(-10, -1),
			NewCriticalPoint(0, 0),
			NewCriticalPoint(10, 1)},
		NewCriticalPoint(10, 0.1),
		NewCriticalPoint(10, 0.1),
		100)

	g1 := NewGroup(a1, a2)
	s := NewSeries(g1, g1, g1)

	bec := BatteryEnergyConstraint(&s, a2.PID(), 1)
	for _, c := range bec {
		err := s.NewConstraint(c)
		assert.Nil(t, err)
	}

	sec := s.Constraints()
	assert.Equal(t, []float64{0, 0, 0, 0, 0, 0, 0, 0, 10, 0, -10, 0, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, sec[0])
	assert.Equal(t, []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 0, -10, 0, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -100, 0}, sec[1])
}
