package cgc_optimize

import (
	"errors"
	"fmt"
	"math"
	"sort"

	"github.com/google/uuid"
)

type BasicUnit struct {
	pid            uuid.UUID
	coefficients   []float64
	bounds         [][2]float64
	constraints    [][]float64
	integrality    []int
	criticalPoints []CriticalPoint
}

type CriticalPoint struct {
	kw         float64
	costPerKwh float64
}

func NewCriticalPoint(kw float64, costPerKwh float64) CriticalPoint {
	return CriticalPoint{kw, costPerKwh}
}

func (cp CriticalPoint) KW() float64 {
	return cp.kw
}

func (cp CriticalPoint) CostPerKWH() float64 {
	return cp.costPerKwh
}

// NewBasicUnit returns a configured unit struct.
func NewBasicUnit(pid uuid.UUID, c []CriticalPoint) BasicUnit {

	// order critical points ascending.
	sort.Slice(c, func(i, j int) bool {
		return (c[i].kw < c[j].kw)
	})

	coefficients := buildCoefficients(c)
	bounds := buildBounds(c)
	binaryMask := buildBinaryMask(c, len(coefficients))
	constraints := buildConstraints(c, len(coefficients))

	return BasicUnit{pid, coefficients, bounds, constraints, binaryMask, c}
}

func buildCoefficients(c []CriticalPoint) []float64 {
	// Variable cost is split into continious segments described by critical points. Constraints are formed to allow
	// only one segement to be active at a time (i.e. the line between two critical points).
	coefficients := []float64{}
	for _, cp := range c {
		coefficients = append(coefficients, cp.costPerKwh)
	}

	binaryCoeff := make([]float64, len(c)-1)
	coefficients = append(coefficients, binaryCoeff...)

	capacityCoeff := make([]float64, 2)
	coefficients = append(coefficients, capacityCoeff...)

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

	// positive and negative capacity bounds
	bounds = append(bounds, [][2]float64{{0, 1}, {0, 1}}...)

	return bounds
}

// buildBinaryMask returns a binary integer slice masking the integer decision variables
func buildBinaryMask(c []CriticalPoint, l int) []int {
	mask := make([]int, l)
	for i := len(c); i < len(c)*2-1; i++ {
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

func (u BasicUnit) PID() uuid.UUID {
	return u.pid
}

func (u BasicUnit) CostCoefficients() []float64 {
	return u.coefficients
}

func (u BasicUnit) ColumnSize() int {
	return len(u.coefficients)
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

func (u BasicUnit) Integrality() []int {
	return u.integrality
}

func (u BasicUnit) CriticalPoints() []CriticalPoint {
	return u.criticalPoints
}

func (u BasicUnit) RealPowerLoc() []int {
	locs := make([]int, len(u.criticalPoints))
	for i := range locs {
		locs[i] = i
	}

	return locs
}

func (u BasicUnit) RealPositiveCapacityLoc() []int {
	loc := []int{len(u.criticalPoints)*2 - 1}

	return loc
}

func (u BasicUnit) RealNegativeCapacityLoc() []int {
	loc := []int{len(u.criticalPoints) * 2}

	return loc
}

// Constraints

func UnitRealPowerConstraint(u *BasicUnit, setpt float64) []float64 {
	rpl := u.RealPowerLoc()
	cp := u.CriticalPoints()

	c := make([]float64, u.ColumnSize())
	for i, loc := range rpl {
		c[loc] = cp[i].KW()
	}

	return boundConstraint(c, setpt, setpt)
}

func UnitPositiveCapacityConstraint(u *BasicUnit, pCap float64) []float64 {
	rpl := u.RealPowerLoc()
	cp := u.CriticalPoints()
	pcl := u.RealPositiveCapacityLoc()[0]

	c := make([]float64, u.ColumnSize())
	for i, loc := range rpl {
		c[loc] = cp[i].KW()
	}

	c[pcl] = -pCap

	c = boundConstraint(c, math.Inf(-1), 0)
	return c
}

func UnitNegativeCapacityConstraint(u *BasicUnit, nCap float64) []float64 {
	rpl := u.RealPowerLoc()
	cp := u.CriticalPoints()
	pcl := u.RealNegativeCapacityLoc()[0]

	c := make([]float64, u.ColumnSize())
	for i, loc := range rpl {
		c[loc] = cp[i].KW()
	}

	c[pcl] = nCap

	c = boundConstraint(c, 0, math.Inf(1))
	return c
}
