package cgc_optimize

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type Series struct {
	elem        []Sequencer
	constraints [][]float64
}

type Sequencer interface {
	CostCoefficients() []float64
	Constraints() [][]float64
	Bounds() [][2]float64
	ColumnSize() int

	RealPowerPidLoc(uuid.UUID) []int
	CriticalPointsPid(uuid.UUID) []CriticalPoint

	StoredEnergyPidLoc(uuid.UUID) []int
	StoredEnergyCapacityPid(uuid.UUID) []float64
}

func NewSeries(sequence ...Sequencer) Series {
	return Series{sequence, [][]float64{}}
}

func (se Series) CostCoefficients() []float64 {
	cc := []float64{}
	for _, cl := range se.elem {
		cc = append(cc, cl.CostCoefficients()...)
	}

	return cc
}

func (se Series) Bounds() [][2]float64 {
	b := make([][2]float64, 0)

	for _, cl := range se.elem {
		b = append(b, cl.Bounds()...)
	}
	return b
}

func (se Series) Constraints() [][]float64 {
	s := se.ColumnSize()
	sec := make([][]float64, 0) // Series Constraint

	i := 0
	for _, cl := range se.elem {
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
	for _, cl := range se.elem {
		s += cl.ColumnSize()
	}

	return s
}

func (se Series) RealPowerPidLoc(t_pid uuid.UUID) []int {
	loc := make([]int, 0)
	i := 0
	for _, cl := range se.elem {
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
	for _, cl := range se.elem {
		for _, p := range cl.StoredEnergyPidLoc(t_pid) {
			loc = append(loc, p+i)
		}
		i += cl.ColumnSize()
	}

	return loc
}

func (se Series) StoredEnergyCapacityPid(t_pid uuid.UUID) []float64 {
	eCap := make([]float64, 0)
	for _, cl := range se.elem {
		for _, e := range cl.StoredEnergyCapacityPid(t_pid) {
			eCap = append(eCap, e)
		}
	}

	return eCap
}

func (se Series) CriticalPointsPid(t_pid uuid.UUID) []CriticalPoint {
	cps := make([]CriticalPoint, 0)
	for _, e := range se.elem {
		cps = append(cps, e.CriticalPointsPid(t_pid)...)
	}

	return cps
}

// BatteryInitialEnergyConstraint returns a constraint of the form: e_t0 = t_e
func BatteryInitialEnergyConstraint(t_se *Series, t_pid uuid.UUID, t_e float64) []float64 {
	eLoc := t_se.StoredEnergyPidLoc(t_pid)
	eCap := t_se.StoredEnergyCapacityPid(t_pid)[0]

	c := make([]float64, t_se.ColumnSize())
	c[eLoc[0]] = t_e
	c = boundConstraint(c, t_e/eCap, t_e/eCap)

	return c
}

// BatteryEnergyConstraint returns a constraint of the form: e_ti - (sum(p_ti))*t = e_t(i+1)
func BatteryEnergyConstraint(t_se *Series, t_pid uuid.UUID, t_tstep float64) [][]float64 {
	// Get ESS critical points and energy capacity
	// Assumes the critical points and stored energy capacity do not change in series.
	cps := t_se.elem[0].CriticalPointsPid(t_pid)
	e := t_se.elem[0].StoredEnergyCapacityPid(t_pid)

	cx := make([][]float64, 0)

	// eCap*e_t0 - (sum(p_t0)) * -tstep = eCap*e_t1
	for i := 0; i < len(t_se.elem)-1; i++ {
		pLoc_t0 := t_se.elem[i].RealPowerPidLoc(t_pid)
		eLoc_t0 := t_se.elem[i].StoredEnergyPidLoc(t_pid)
		eLoc_t1 := t_se.elem[i+1].StoredEnergyPidLoc(t_pid)

		c := make([]float64, t_se.ColumnSize())
		for i, loc := range pLoc_t0 {
			c[loc] = cps[i].KW() * -t_tstep
		}

		c[eLoc_t0[0]] = e[0]
		c[eLoc_t1[0]] = -e[0]

		c = boundConstraint(c, 0, 0)
		cx = append(cx, c)
	}

	return cx
}
