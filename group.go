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

func (g Group) CriticalPoints() []CriticalPoint {
	cp := make([]CriticalPoint, 0)
	for _, u := range g.units {
		cp = append(cp, u.CriticalPoints()...)
	}
	return cp
}

func (g Group) RealPositiveCapacity() []float64 {
	rpc := make([]float64, 0)
	for _, u := range g.units {
		rpc = append(rpc, u.RealPositiveCapacity()...)
	}
	return rpc
}

func (g Group) RealNegativeCapacity() []float64 {
	rpc := make([]float64, 0)
	for _, u := range g.units {
		rpc = append(rpc, u.RealNegativeCapacity()...)
	}
	return rpc
}

func (g Group) StoredEnergyCapacityPid(t_pid uuid.UUID) []float64 {
	eCap := make([]float64, 0)
	for _, u := range g.units {
		switch v := u.(type) {
		case EnergyStorageUnit:
			if u.PID() == t_pid {
				eCap = append(eCap, v.StoredEnergyCapacity()...)
			}
		default:
		}
	}
	return eCap
}

func (g Group) CriticalPointsPid(t_pid uuid.UUID) []CriticalPoint {
	cp := make([]CriticalPoint, 0)
	for _, u := range g.units {
		if u.PID() == t_pid {
			cp = append(cp, u.CriticalPoints()...)
		}
	}
	return cp
}

func (g Group) RealPowerLoc() []int {
	loc := make([]int, 0)
	i := 0
	for _, u := range g.units {
		for _, p := range u.RealPowerLoc() {
			loc = append(loc, p+i)
		}
		i += u.ColumnSize()
	}
	return loc
}

func (g Group) RealPositiveCapacityLoc() []int {
	loc := make([]int, 0)
	i := 0
	for _, u := range g.units {
		for _, p := range u.RealPositiveCapacityLoc() {
			loc = append(loc, p+i)
		}
		i += u.ColumnSize()
	}
	return loc
}

func (g Group) RealNegativeCapacityLoc() []int {
	loc := make([]int, 0)
	i := 0
	for _, u := range g.units {
		for _, p := range u.RealNegativeCapacityLoc() {
			loc = append(loc, p+i)
		}
		i += u.ColumnSize()
	}
	return loc
}

func (g Group) RealPowerPidLoc(t_pid uuid.UUID) []int {
	loc := make([]int, 0)
	i := 0
	for _, u := range g.units {
		if u.PID() == t_pid {
			for _, p := range u.RealPowerLoc() {
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
		switch v := u.(type) {
		case EnergyStorageUnit:
			if u.PID() == t_pid {
				for _, p := range v.StoredEnergyLoc() {
					loc = append(loc, p+i)
				}
			}
		default:
		}
		i += u.ColumnSize()
	}
	return loc
}

// Constraint Generation

// NetLoadPiecewiseConstraint return a constraint of the form: Sum_i(x1+x2+...xn) == t_nl
func NetLoadConstraint(g *Group, t_nl float64) []float64 {
	c := make([]float64, g.ColumnSize())
	rpl := g.RealPowerLoc() // RealNegativePowerLoc would return same value.
	cps := g.CriticalPoints()
	for i, loc := range rpl {
		c[loc] = cps[i].KW()
	}

	return boundConstraint(c, t_nl, t_nl)
}

// GroupPositiveCapacityConstriant returns a constraint
func GroupPositiveCapacityConstraint(g *Group, t_cap float64) []float64 {
	c := make([]float64, g.ColumnSize())
	rpcl := g.RealPositiveCapacityLoc()
	rpc := g.RealPositiveCapacity()
	for i, loc := range rpcl {
		c[loc] = rpc[i]
	}

	return boundConstraint(c, t_cap, math.Inf(1))
}

// GroupNegativeCapacityConstriant returns a constraint
func GroupNegativeCapacityConstraint(g *Group, t_cap float64) []float64 {
	c := make([]float64, g.ColumnSize())
	rncl := g.RealNegativeCapacityLoc()
	rnc := g.RealNegativeCapacity()
	for i, loc := range rncl {
		c[loc] = rnc[i]
	}

	return boundConstraint(c, t_cap, math.Inf(1))
}
