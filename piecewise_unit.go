package cgc_optimize

import (
	"errors"
	"fmt"
	"math"
	"sort"

	"github.com/google/uuid"
)

type PiecewiseUnit struct {
	pid            uuid.UUID
	coefficients   []float64
	bounds         [][2]float64
	constraints    [][]float64
	integrality    []int
	criticalPoints []CriticalPoint
}

type CriticalPoint struct {
	val  float64
	cost float64
}

func NewCriticalPoint(val float64, cost float64) CriticalPoint {
	return CriticalPoint{val, cost}
}

func (cp CriticalPoint) Value() float64 {
	return cp.val
}

func (cp CriticalPoint) Cost() float64 {
	return cp.cost
}

// NewPiecewiseUnit returns a configured unit struct.
func NewPiecewiseUnit(pid uuid.UUID, c []CriticalPoint) PiecewiseUnit {

	// order critical points ascending.
	sort.Slice(c, func(i, j int) bool {
		return (c[i].val < c[j].val)
	})

	coefficients := buildCoefficients(c)
	bounds := buildBounds(c)
	binaryMask := buildBinaryMask(c, len(coefficients))
	constraints := buildConstraints(c, len(coefficients))

	return PiecewiseUnit{pid, coefficients, bounds, constraints, binaryMask, c}
}

func buildCoefficients(c []CriticalPoint) []float64 {
	// Variable cost is split into continious segments described by critical points. Constraints are formed to allow
	// only one segement to be active at a time (i.e. the line between two critical points).
	coefficients := []float64{}
	for _, cp := range c {
		coefficients = append(coefficients, cp.cost)
	}

	binaryCoeff := make([]float64, len(c)-1)
	coefficients = append(coefficients, binaryCoeff...)

	return coefficients
}

func buildBounds(c []CriticalPoint) [][2]float64 {

	bounds := [][2]float64{}
	for i := 0; i < len(c); i++ {
		bounds = append(bounds, [2]float64{0, math.Inf(1)})
	}

	// set bounds on binary decision vairables
	for i := 0; i < len(c)-1; i++ {
		bounds = append(bounds, [2]float64{0, 1})
	}

	return bounds
}

// buildBinaryMask returns a binary integer slice masking the integer decision variables
func buildBinaryMask(c []CriticalPoint, l int) []int {
	mask := make([]int, l)
	for i := len(c); i < l; i++ {
		mask[i] = 1
	}

	return mask
}

// buildConstraintsSegement creates a lower diagonal matrix
func buildConstraintsSegment(c []CriticalPoint, l int) [][]float64 {
	constraints := [][]float64{}
	for i := range c {
		constraint := make([]float64, l)
		constraint[i] = 1

		if i < len(c)-1 {
			constraint[i+len(c)] = -1
		}

		if i > 0 {
			constraint[i+len(c)-1] = -1
		}

		constraints = append(constraints, boundConstraint(constraint, math.Inf(-1), 0))
	}

	return constraints
}

func buildConstraintsInterpolation(c []CriticalPoint, l int) []float64 {
	// segment interpolation constraint
	// [0, 1, 1, 1, 0, 0, 1]
	constraint := make([]float64, l)
	for i := range c {
		constraint[i] = 1
	}

	return boundConstraint(constraint, 1, 1)
}

func buildConstraintsActivation(c []CriticalPoint, l int) []float64 {

	// segment activation constraint
	// [0, 0, 0, 0, 1, 1, 1]
	constraint := make([]float64, l)
	for i := len(c); i < l; i++ {
		constraint[i] = 1
	}

	return boundConstraint(constraint, 1, 1)
}

func buildConstraints(c []CriticalPoint, l int) [][]float64 {

	constraints := buildConstraintsSegment(c, l)
	constraints = append(constraints, buildConstraintsInterpolation(c, l))
	constraints = append(constraints, buildConstraintsActivation(c, l))

	return constraints
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

func (u PiecewiseUnit) Integrality() []int {
	return u.integrality
}

func (u PiecewiseUnit) CriticalPoints() []CriticalPoint {
	return u.criticalPoints
}

func (u PiecewiseUnit) RealPositivePowerLoc() []int {
	locs := make([]int, len(u.criticalPoints))
	for i := range locs {
		locs[i] = i
	}

	return locs
}

func (u PiecewiseUnit) RealNegativePowerLoc() []int {
	loc := make([]int, len(u.criticalPoints))
	for i := range loc {
		loc[i] = i
	}

	return loc
}

func (u PiecewiseUnit) RealCapacityLoc() []int {
	loc := make([]int, len(u.criticalPoints))
	for i := range loc {
		loc[i] = i
	}

	return loc
}

func (u PiecewiseUnit) StoredEnergyLoc() []int {
	loc := make([]int, len(u.criticalPoints))
	for i := range loc {
		loc[i] = i
	}

	return loc
}

// Constraints

/*
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
*/
