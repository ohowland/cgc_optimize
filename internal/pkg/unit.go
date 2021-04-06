package esslp

import (
	"errors"
	"math"

	"github.com/google/uuid"
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

func (u *Unit) NewConstraint(t_c []float64, t_lb float64, t_ub float64) error {
	if len(t_c) != u.ColumnSize() {
		return errors.New("column size mismatch")
	}
	c := []float64{t_lb}
	c = append(c, t_c...)
	c = append(c, t_ub)

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

func (u Unit) RealPositivePowerLoc() []int {
	return []int{0}
}

func (u Unit) RealNegativePowerLoc() []int {
	return []int{1}
}

func (u Unit) RealCapacityLoc() []int {
	return []int{2}
}

func (u Unit) StoredCapacityLoc() []int {
	return []int{3}
}
