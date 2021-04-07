package la

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type Cluster struct {
	groups      []Group
	constraints [][]float64
}

func NewCluster(groups ...Group) Cluster {
	return Cluster{groups, [][]float64{}}
}

func (cl Cluster) CostCoefficients() []float64 {
	cc := []float64{}
	for _, g := range cl.groups {
		cc = append(cc, g.CostCoefficients()...)
	}

	return cc
}

func (cl Cluster) Constraints() [][]float64 {
	s := cl.ColumnSize()
	clc := make([][]float64, 0) // Cluster Constraint

	i := 0
	for _, g := range cl.groups {
		pre := make([]float64, i+1)
		post := make([]float64, s-i-g.ColumnSize()+1)
		for _, gc := range g.Constraints() {
			pre[0] = lb(gc)            // shift lower bound to index 0
			post[len(post)-1] = ub(gc) // shift upper bound at last index
			c := append(append(pre, cons(gc)...), post...)
			clc = append(clc, c)
		}
		i += g.ColumnSize()
	}

	clc = append(clc, cl.constraints...)
	return clc
}

func (cl Cluster) Bounds() [][2]float64 {
	b := make([][2]float64, 0)

	for _, g := range cl.groups {
		b = append(b, g.Bounds()...)
	}
	return b
}

func (cl Cluster) ColumnSize() int {
	var s int
	for _, g := range cl.groups {
		s += g.ColumnSize()
	}

	return s
}

func (cl *Cluster) NewConstraint(t_c ...[]float64) error {
	cx := make([][]float64, 0)
	for _, c := range t_c {
		if len(c) != cl.ColumnSize()+2 {
			err := fmt.Sprintf("constraint contains %v columns, expected: %v", len(c), cl.ColumnSize()+2)
			return errors.New(err)
		}
		cx = append(cx, c)
	}

	cl.constraints = append(cl.constraints, cx...)
	return nil
}

func (cl Cluster) RealPositivePowerPidLoc(t_pid uuid.UUID) []int {
	loc := make([]int, 0)
	i := 0
	for _, g := range cl.groups {
		for _, p := range g.RealPositivePowerPidLoc(t_pid) {
			loc = append(loc, p+i)
		}
		i += g.ColumnSize()
	}

	return loc
}

func (cl Cluster) RealNegativePowerPidLoc(t_pid uuid.UUID) []int {
	loc := make([]int, 0)
	i := 0
	for _, g := range cl.groups {
		for _, p := range g.RealNegativePowerPidLoc(t_pid) {
			loc = append(loc, p+i)
		}
		i += g.ColumnSize()
	}

	return loc
}

func (cl Cluster) StoredEnergyPidLoc(t_pid uuid.UUID) []int {
	loc := make([]int, 0)
	i := 0
	for _, g := range cl.groups {
		for _, p := range g.StoredEnergyPidLoc(t_pid) {
			loc = append(loc, p+i)
		}
		i += g.ColumnSize()
	}

	return loc
}

// Cluster Specific Constraints

func LinkedBusConstraints(t_cl *Cluster, t_pid uuid.UUID) [][]float64 {
	pLoc := t_cl.RealPositivePowerPidLoc(t_pid) // location of Positive Real Power decision variables
	nLoc := t_cl.RealNegativePowerPidLoc(t_pid) // location of Negative Real Power decision variables

	if len(pLoc) != 2 || len(nLoc) != 2 {
		return [][]float64{}
	}

	pc := make([]float64, t_cl.ColumnSize())
	pc[pLoc[0]] = 1
	pc[pLoc[1]] = -1
	pc = boundConstraint(pc, 0, 0)

	nc := make([]float64, t_cl.ColumnSize())
	nc[nLoc[0]] = 1
	nc[nLoc[1]] = -1
	nc = boundConstraint(nc, 0, 0)

	c := make([][]float64, 0)
	c = append(append(c, pc), nc)

	return c
}
