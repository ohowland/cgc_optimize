package esslp

import (
	"errors"
	"math"
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

func (g *Group) NewConstraint(t_c []float64, t_lb float64, t_ub float64) error {
	if len(t_c) != g.ColumnSize() {
		return errors.New("column size mismatch")
	}

	c := []float64{t_lb}
	c = append(c, t_c...) // constraint
	c = append(c, t_ub)   // upper bound

	g.constraints = append(g.constraints, c)
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
			pre[0] = uc[0]                    // shift lower bound to index 0
			post[len(post)-1] = uc[len(uc)-1] // shift upper bound at last index

			c := append(append(pre, uc[1:len(uc)-1]...), post...)
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

func (g Group) StoredCapacityLoc() []int {
	loc := make([]int, 0)
	i := 0
	for _, u := range g.units {
		for _, p := range u.StoredCapacityLoc() {
			loc = append(loc, p+i)
		}
		i += u.ColumnSize()
	}
	return loc
}

// Constraint Generation

func NetLoadConstraint(t_nl float64, g *Group) []float64 {
	c := make([]float64, g.ColumnSize())

	rpp := g.RealPositivePowerLoc()
	rnp := g.RealNegativePowerLoc()
	for _, i := range rpp {
		c[i] = 1.0
	}
	for _, i := range rnp {
		c[i] = -1.0
	}

	lbub := []float64{t_nl}

	return append(append(lbub, c...), lbub...) // [low_bound, constraints, upper_bound]
}

func PositiveCapacityConstraint(t_cap float64, g *Group) []float64 {
	c := make([]float64, g.ColumnSize())

	pc := g.RealCapacityLoc()
	for _, i := range pc {
		c[i] = 1.0
	}

	lb := []float64{t_cap}

	return append(append(lb, c...), math.Inf(1))
}
