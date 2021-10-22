package cgc_optimize

import (
	"errors"
	"fmt"
	"math"

	"github.com/google/uuid"
)

type Group struct {
	units       []Unit
	constraints [][]float64
}

func NewGroup(units ...Unit) Group {
	ux := make([]Unit, 0)
	ux = append(ux, units...)
	cx := make([][]float64, 0)

	return Group{ux, cx}
}

func (g Group) CostCoefficients() []float64 {
	cx := make([]float64, 0)
	for _, u := range g.units {
		cx = append(cx, u.CostCoefficients()...)
	}

	return cx
}

func (g Group) ColumnSize() int {
	var s int
	for _, u := range g.units {
		s += u.ColumnSize()
	}

	return s
}

func (g *Group) NewConstraint(t_c ...[]float64) error {
	cx := make([][]float64, 0)
	for _, c := range t_c {
		if len(c) != g.ColumnSize()+2 {
			err := fmt.Sprintf("constraint contains %v columns, expected: %v", len(c), g.ColumnSize()+2)
			return errors.New(err)
		}
		cx = append(cx, c)
	}

	// if no errors, append constraints to group
	g.constraints = append(g.constraints, cx...)
	return nil
}

func (g Group) Constraints() [][]float64 {
	s := g.ColumnSize()
	gc := make([][]float64, 0)

	i := 0
	for _, u := range g.units {
		pre := make([]float64, i+1)
		post := make([]float64, s-i-u.ColumnSize()+1)
		for _, uc := range u.Constraints() {
			pre[0] = lb(uc)            // shift lower bound to index 0
			post[len(post)-1] = ub(uc) // shift upper bound at last index

			c := append(append(pre, cons(uc)...), post...)
			gc = append(gc, c)
		}
		i += u.ColumnSize()
	}

	gc = append(gc, g.constraints...)
	return gc
}

func (g Group) Bounds() [][2]float64 {
	b := make([][2]float64, 0)

	for _, u := range g.units {
		b = append(b, u.Bounds()...)
	}
	return b
}

func (g Group) Integrality() []int {
	i := make([]int, 0)
	for _, u := range g.units {
		i = append(i, u.Integrality()...)
	}
	return i
}

func (g Group) RealPositivePowerLoc() []int {
	loc := make([]int, 0)
	i := 0
	for _, u := range g.units {
		for _, p := range u.RealPositivePowerLoc() {
			loc = append(loc, p+i)
		}
		i += u.ColumnSize()
	}
	return loc
}

func (g Group) CriticalPoints() []CriticalPoint {
	cp := make([]CriticalPoint, 0)
	for _, u := range g.units {
		cp = append(cp, u.CriticalPoints()...)
	}

	return cp
}

func (g Group) RealNegativePowerLoc() []int {
	loc := make([]int, 0)
	i := 0
	for _, u := range g.units {
		for _, p := range u.RealNegativePowerLoc() {
			loc = append(loc, p+i)
		}
		i += u.ColumnSize()
	}
	return loc
}

func (g Group) RealCapacityLoc() []int {
	loc := make([]int, 0)
	i := 0
	for _, u := range g.units {
		for _, p := range u.RealCapacityLoc() {
			loc = append(loc, p+i)
		}
		i += u.ColumnSize()
	}
	return loc
}

func (g Group) StoredEnergyLoc() []int {
	loc := make([]int, 0)
	i := 0
	for _, u := range g.units {
		for _, p := range u.StoredEnergyLoc() {
			loc = append(loc, p+i)
		}
		i += u.ColumnSize()
	}
	return loc
}

func (g Group) RealPositivePowerPidLoc(t_pid uuid.UUID) []int {
	loc := make([]int, 0)
	i := 0
	for _, u := range g.units {
		if u.PID() == t_pid {
			for _, p := range u.RealPositivePowerLoc() {
				loc = append(loc, p+i)
			}
		}
		i += u.ColumnSize()
	}
	return loc
}

func (g Group) RealNegativePowerPidLoc(t_pid uuid.UUID) []int {
	loc := make([]int, 0)
	i := 0
	for _, u := range g.units {
		if u.PID() == t_pid {
			for _, p := range u.RealNegativePowerLoc() {
				loc = append(loc, p+i)
			}
		}
		i += u.ColumnSize()
	}
	return loc
}

func (g Group) StoredEnergyPidLoc(t_pid uuid.UUID) []int {
	loc := make([]int, 0)
	i := 0
	for _, u := range g.units {
		if u.PID() == t_pid {
			for _, p := range u.StoredEnergyLoc() {
				loc = append(loc, p+i)
			}
		}
		i += u.ColumnSize()
	}
	return loc
}

// Constraint Generation

// NetLoadConstraint returns a constraint of the form: Sum_i(Xp_i - Xn_i) == t_nl
func NetLoadConstraint(g *Group, t_nl float64) []float64 {
	c := make([]float64, g.ColumnSize())

	rpp := g.RealPositivePowerLoc()
	rnp := g.RealNegativePowerLoc()
	for _, i := range rpp {
		c[i] = 1.0
	}
	for _, i := range rnp {
		c[i] = -1.0
	}

	return boundConstraint(c, t_nl, t_nl)
}

// NetLoadPiecewiseConstraint return a constraint of the form: Sum_i(x1+x2+...xn) == t_nl
func NetLoadPiecewiseConstraint(g *Group, t_nl float64) []float64 {
	c := make([]float64, g.ColumnSize())
	locs := g.RealPositivePowerLoc()
	cps := g.CriticalPoints()
	for i, loc := range locs {
		c[loc] = cps[i].Value()
	}

	return c
}

// GroupCapacityConstriant returns a constraint of the form: Sum_i(Xc_i) >= t_cap
func GroupPositiveCapacityConstraint(g *Group, t_cap float64) []float64 {
	c := make([]float64, g.ColumnSize())

	pc := g.RealCapacityLoc()
	for _, i := range pc {
		c[i] = 1.0
	}

	return boundConstraint(c, t_cap, math.Inf(1))
}
