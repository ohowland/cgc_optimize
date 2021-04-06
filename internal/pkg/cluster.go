package esslp

import (
	"errors"
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
			pre[0] = gc[0]                    // shift lower bound to index 0
			post[len(post)-1] = gc[len(gc)-1] // shift upper bound at last index
			c := append(append(pre, gc[1:len(gc)-1]...), post...)
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

func (cl *Cluster) NewConstraint(t_c []float64, t_lb float64, t_ub float64) error {
	if len(t_c) != cl.ColumnSize() {
		return errors.New("column size mismatch")
	}

	c := []float64{t_lb}
	c = append(c, t_c...) // constraint
	c = append(c, t_ub)   // upper bound

	cl.constraints = append(cl.constraints, c)
	return nil
}
