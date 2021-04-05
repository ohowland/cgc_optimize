package esslp

import (
	"errors"
	"math"

	"github.com/google/uuid"
	"github.com/lanl/clp"
)

type Unit struct {
	pid         uuid.UUID
	Cp          float64 // cost of energy production
	Cn          float64 // cost of energy consumption
	Cc          float64 // cost of energy capacity
	Ce          float64 // cost of energy storage
	constraints [][]float64
	bounds      [][]float64
}

func NewUnit(pid uuid.UUID, Cp float64, Cn float64, Cc float64, Ce float64) Unit {
	return Unit{pid, Cp, Cn, Cc, Ce, [][]float64{}, [][]float64{}}
}

func (u Unit) CostCoefficients() []float64 {
	return []float64{u.Cp, u.Cn, u.Cc, u.Ce}
}

func (u Unit) ColumnSize() int {
	return 4
}

func (u *Unit) NewConstraint(c []float64) error {
	if len(c) != u.ColumnSize() {
		return errors.New("column size mismatch")
	}
	u.constraints = append(u.constraints, c)
	return nil
}

func (u Unit) Constraints() [][]float64 {
	return u.constraints
}

func (u Unit) Bounds() [][2]float64 {
	inf := math.Inf(1)
	return [][2]float64{
		{0, inf},
		{0, inf},
		{0, inf},
		{0, inf},
	}
}

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

func (g *Group) NewConstraint(c []float64) error {
	if len(c) != g.ColumnSize() {
		return errors.New("column size mismatch")
	}
	g.constraints = append(g.constraints, c)
	return nil
}

func (g Group) Constraints() [][]float64 {
	return g.constraints
}

func (g Group) Bounds() [][2]float64 {
	b := make([][2]float64, 0)

	for _, u := range g.units {
		b = append(b, u.Bounds()...)
	}
	return b
}

// Linear Program

type LinearProgram struct {
	group *Group
}

func NewLP(g *Group) LinearProgram {
	return LinearProgram{g}
}

func (lp *LinearProgram) Solve() []float64 {
	s := clp.NewSimplex()
	s.EasyLoadDenseProblem(
		lp.group.CostCoefficients(),
		lp.group.Bounds(),
		lp.group.Constraints(),
	)

	s.SetOptimizationDirection(clp.Minimize)
	s.Primal(clp.NoValuesPass, clp.NoStartFinishOptions)
	return s.PrimalColumnSolution()
}
