package cgc_optimize

import (
	"errors"
	"fmt"
	"math"

	"github.com/google/uuid"
)

type PiecewiseUnit struct {
	pid            uuid.UUID
	coefficients   []float64
	bounds         [][2]float64
	constraints    [][]float64
	binaries       []int
	criticalPoints []CriticalPoint
}

type CriticalPoint struct {
	val  float64
	cost float64
}

// NewPiecewiseUnit returns a configured unit struct.
//
// XpUb: Upper bound for real positive power decision variable
// XnUb: Upper bound for real negative power decision variable (positive value)
// XcUb: Upper bound for real capacity decision variable
func NewPiecewiseUnit(pid uuid.UUID, C []CriticalPoint) PiecewiseUnit {

	// Variable cost is split into continious segments described by critical points. Constraints are formed to allow
	// only one segement to be active at a time (i.e. the line between two critical points).
	coefficients := []float64{}
	for _, cp := range C {
		coefficients = append(coefficients, cp.cost)
	}

	binaryCoeff := make([]float64, len(C)-1)
	coefficients = append(coefficients, binaryCoeff...)

	// set bounds on segment decision variables
	bounds := [][2]float64{}
	for i := 0; i < len(C); i++ {
		if C[i].val >= 0 {
			bounds = append(bounds, [2]float64{0, C[i].val})
		} else {
			bounds = append(bounds, [2]float64{C[i].val, 0})
		}
	}

	// set bounds on binary decision vairables
	for i := len(C); i < len(coefficients); i++ {
		bounds = append(bounds, [2]float64{0, 1})
	}

	// mask binary decision variables
	binaryIndex := make([]int, len(coefficients))
	for i := len(C); i < len(coefficients); i++ {
		binaryIndex[i] = 1
	}

	// create segment constraints, this is a diagonal matrix
	constraints := [][]float64{}
	for i := range C {
		constraint := make([]float64, len(coefficients))
		constraint[i] = 1

		if i < len(C)-1 {
			constraint[i+len(C)] = -1
		}
		if i > 0 {
			constraint[i+len(C)-1] = -1
		}

		constraints = append(constraints, constraint)
	}

	return PiecewiseUnit{pid, coefficients, bounds, constraints, binaryIndex, C}
}

func (u PiecewiseUnit) PID() uuid.UUID {
	return u.pid
}

func (u PiecewiseUnit) CostCoefficients() []float64 {
	return u.coefficients
}

func (u PiecewiseUnit) ColumnSize() int {
	return len(u.coefficients)
}

func (u *PiecewiseUnit) NewConstraint(t_c ...[]float64) error {
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

func (u PiecewiseUnit) Constraints() [][]float64 {
	return u.constraints
}

func (u PiecewiseUnit) Bounds() [][2]float64 {
	return u.bounds
}

func (u PiecewiseUnit) RealPositivePowerLoc() []int {
	return []int{0}
}

func (u PiecewiseUnit) RealNegativePowerLoc() []int {
	return []int{1}
}

func (u PiecewiseUnit) RealCapacityLoc() []int {
	return []int{2}
}

// Constraints

func PiecewiseUnitCapacityConstraints(u *PiecewiseUnit) [][]float64 {
	cx := make([][]float64, 0)
	cx = append(cx, PiecewiseUnitPositiveCapacityConstraint(u))
	cx = append(cx, PiecewiseUnitNegativeCapacityConstraint(u))
	return cx
}

func PiecewiseUnitPositiveCapacityConstraint(u *PiecewiseUnit) []float64 {
	xp := u.RealPositivePowerLoc()[0]
	xc := u.RealCapacityLoc()[0]

	cp := make([]float64, u.ColumnSize())
	cp[xp] = -1
	cp[xc] = 1

	cp = boundConstraint(cp, 0, math.Inf(1))
	return cp
}

func PiecewiseUnitNegativeCapacityConstraint(u *PiecewiseUnit) []float64 {
	xn := u.RealNegativePowerLoc()[0]
	xc := u.RealCapacityLoc()[0]

	cn := make([]float64, u.ColumnSize())
	cn[xn] = -1
	cn[xc] = 1

	cn = boundConstraint(cn, 0, math.Inf(1))
	return cn
}
