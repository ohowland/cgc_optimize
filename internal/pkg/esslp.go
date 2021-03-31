package esslp

import (
	"errors"

	"github.com/google/uuid"
)

type Unit struct {
	pid         uuid.UUID
	Cp          float64 // cost of energy production
	Cn          float64 // cost of energy consumption
	Cc          float64 // cost of energy capacity
	Ce          float64 // cost of energy storage
	constraints [][]float64
}

func NewUnit(pid uuid.UUID, Cp float64, Cn float64, Cc float64, Ce float64) Unit {
	return Unit{pid, Cp, Cn, Cc, Ce, [][]float64{}}
}

func (u Unit) Flatten() []float64 {
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

type Group struct {
	Units []Unit
}

func NewGroup(units ...Unit) Group {
	ux := make([]Unit, 0)
	ux = append(ux, units...)

	return Group{ux}
}

func (g Group) CostCoefficients() []float64 {
	cx := make([]float64, 0)
	for _, u := range g.Units {
		cx = append(cx, u.Flatten()...)
	}

	return cx
}

func (g Group) ColumnSize() int {
	var s int
	for _, u := range g.Units {
		s += u.ColumnSize()
	}

	return s
}
