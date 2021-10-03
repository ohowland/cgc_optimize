package cgc_optimize

import (
	"errors"
	"fmt"
	"math"

	"github.com/google/uuid"
)

type Unit interface {
	PID() uuid.UUID
	CostCoefficients() []float64
	Bounds() [][2]float64
	Constraints() [][]float64
	ColumnSize() int

	RealPower() float64
}

type BasicUnit struct {
	pid          uuid.UUID
	coefficients []float64
	bounds       [][2]float64
	constraints  [][]float64
}

// NewBasicUnit returns a configured unit struct.
//
// Cp: Cost coefficient for real positive power
// Cn: Cost coefficient for real negative power
// Cc: Cost coefficient for real capacity
// Ce: Cost coefficient for stored energy
//
// XpUb: Upper bound for real positive power decision variable
// XnUb: Upper bound for real negative power decision variable (positive value)
// XcUb: Upper bound for real capacity decision variable
// XeUb: Upper bound for stored energy
func NewBasicUnit(pid uuid.UUID, Cp float64, Cn float64, Cc float64, Ce float64, XpUb float64, XnUb float64, XcUb float64, XeUb float64) BasicUnit {
	coefficients := []float64{Cp, Cn, Cc, Ce}
	bounds := [][2]float64{{0, XpUb}, {0, XnUb}, {0, XcUb}, {0, XeUb}}

	return BasicUnit{pid, coefficients, bounds, [][]float64{}}
}

func (u BasicUnit) PID() uuid.UUID {
	return u.pid
}

func (u BasicUnit) CostCoefficients() []float64 {
	return u.coefficients
}

func (u BasicUnit) ColumnSize() int {
	return 4
}

func (u *BasicUnit) NewConstraint(t_c ...[]float64) error {
	cx := make([][]float64, 0)
	for _, c := range t_c {
		if len(c) != u.ColumnSize()+2 {
			err := fmt.Sprintf("constraint contains %v columns, expected: %v", len(c), u.ColumnSize()+2)
			return errors.New(err)
		}
		cx = append(cx, c)
	}

	// if no errors: add constraints to unit
	u.constraints = append(u.constraints, cx...)
	return nil
}

func (u BasicUnit) Constraints() [][]float64 {
	return u.constraints
}

func (u BasicUnit) Bounds() [][2]float64 {
	return u.bounds
}

func (u BasicUnit) RealPositivePowerLoc() []int {
	return []int{0}
}

func (u BasicUnit) RealNegativePowerLoc() []int {
	return []int{1}
}

func (u BasicUnit) RealCapacityLoc() []int {
	return []int{2}
}

func (u BasicUnit) StoredEnergyLoc() []int {
	return []int{3}
}

// Constraints

func BasicUnitCapacityConstraints(u *BasicUnit) [][]float64 {
	cx := make([][]float64, 0)
	cx = append(cx, BasicUnitPositiveCapacityConstraint(u))
	cx = append(cx, BasicUnitNegativeCapacityConstraint(u))
	return cx
}

func BasicUnitPositiveCapacityConstraint(u *BasicUnit) []float64 {
	xp := u.RealPositivePowerLoc()[0]
	xc := u.RealCapacityLoc()[0]

	cp := make([]float64, u.ColumnSize())
	cp[xp] = -1
	cp[xc] = 1

	cp = boundConstraint(cp, 0, math.Inf(1))
	return cp
}

func BasicUnitNegativeCapacityConstraint(u *BasicUnit) []float64 {
	xn := u.RealNegativePowerLoc()[0]
	xc := u.RealCapacityLoc()[0]

	cn := make([]float64, u.ColumnSize())
	cn[xn] = -1
	cn[xc] = 1

	cn = boundConstraint(cn, 0, math.Inf(1))
	return cn
}
