package cgc_optimize

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type Series struct {
	clusters    []Sequencer
	constraints [][]float64
}

type Sequencer interface {
	CostCoefficients() []float64
	Constraints() [][]float64
	Bounds() [][2]float64
	ColumnSize() int
	PowerLoc
	StorageLoc
}

type PowerLoc interface {
	RealPowerPidLoc(uuid.UUID) []int
}

type StorageLoc interface {
	StoredEnergyPidLoc(uuid.UUID) []int
}

func NewSeries(sequence ...Sequencer) Series {
	return Series{sequence, [][]float64{}}
}

func (se Series) CostCoefficients() []float64 {
	cc := []float64{}
	for _, cl := range se.clusters {
		cc = append(cc, cl.CostCoefficients()...)
	}

	return cc
}

func (se Series) Bounds() [][2]float64 {
	b := make([][2]float64, 0)

	for _, cl := range se.clusters {
		b = append(b, cl.Bounds()...)
	}
	return b
}

func (se Series) Constraints() [][]float64 {
	s := se.ColumnSize()
	sec := make([][]float64, 0) // Series Constraint

	i := 0
	for _, cl := range se.clusters {
		pre := make([]float64, i+1)
		post := make([]float64, s-i-cl.ColumnSize()+1)
		for _, clc := range cl.Constraints() {
			pre[0] = lb(clc)            // shift lower bound to index 0
			post[len(post)-1] = ub(clc) // shift upper bound at last index
			c := append(append(pre, cons(clc)...), post...)
			sec = append(sec, c)
		}
		i += cl.ColumnSize()
	}

	sec = append(sec, se.constraints...)
	return sec
}

func (se *Series) NewConstraint(t_c ...[]float64) error {
	cx := make([][]float64, 0)
	for _, c := range t_c {
		if len(c) != se.ColumnSize()+2 {
			err := fmt.Sprintf("constraint contains %v columns, expected: %v", len(c), se.ColumnSize()+2)
			return errors.New(err)
		}
		cx = append(cx, c)

	}

	se.constraints = append(se.constraints, cx...)
	return nil
}

func (se *Series) ColumnSize() int {
	var s int
	for _, cl := range se.clusters {
		s += cl.ColumnSize()
	}

	return s
}

func (se Series) RealPowerPidLoc(t_pid uuid.UUID) []int {
	loc := make([]int, 0)
	i := 0
	for _, cl := range se.clusters {
		for _, p := range cl.RealPowerPidLoc(t_pid) {
			loc = append(loc, p+i)
		}
		i += cl.ColumnSize()
	}

	return loc
}

func (se Series) StoredEnergyPidLoc(t_pid uuid.UUID) []int {
	loc := make([]int, 0)
	i := 0
	for _, cl := range se.clusters {
		for _, p := range cl.StoredEnergyPidLoc(t_pid) {
			loc = append(loc, p+i)
		}
		i += cl.ColumnSize()
	}

	return loc
}

// BatteryInitialEnergyConstraint returns a constraint of the form: e_t0 = t_e
func BatteryInitialEnergyConstraint(t_se *Series, t_pid uuid.UUID, t_e float64) []float64 {
	eLoc := t_se.StoredEnergyPidLoc(t_pid)
	c := make([]float64, t_se.ColumnSize())
	c[eLoc[0]] = 1
	c = boundConstraint(c, t_e, t_e)

	return c
}

// BatteryEnergyConstraint returns a constraint of the form: e_ti - (p_ti-n_ti)*t = e_t(i+1)
func BatteryEnergyConstraint(t_se *Series, t_pid uuid.UUID, t_tstep float64) [][]float64 {
	pLoc := t_se.RealPowerPidLoc(t_pid)
	eLoc := t_se.StoredEnergyPidLoc(t_pid)

	cx := make([][]float64, 0)
	for i := 0; i < len(eLoc)-1; i++ {
		c := make([]float64, t_se.ColumnSize())
		c[pLoc[i]] = -t_tstep
		c[eLoc[i]] = 1
		c[eLoc[i+1]] = -1
		c = boundConstraint(c, 0, 0)
		cx = append(cx, c)
	}

	return cx
}
