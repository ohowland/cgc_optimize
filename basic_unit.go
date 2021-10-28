package cgc_optimize

import (
	"errors"
	"fmt"
	"math"
	"sort"

	"github.com/google/uuid"
)

type BasicUnit struct {
	pid          uuid.UUID
	coefficients []float64
	bounds       [][2]float64
	constraints  [][]float64
	integrality  []int

	criticalPoints       []CriticalPoint
	realPositiveCapacity CriticalPoint
	realNegativeCapacity CriticalPoint
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
func NewBasicUnit(pid uuid.UUID, c []CriticalPoint, pCap CriticalPoint, nCap CriticalPoint) BasicUnit {

	if pCap.KW() < 0 || nCap.KW() < 0 {
		panic("basic unit positive and negative capacity must be greater than or equal to 0")
	}
	// order critical points ascending.
	sort.Slice(c, func(i, j int) bool {
		return (c[i].kw < c[j].kw)
	})

	coefficients := buildCoefficients(c, pCap, nCap)
	bounds := buildBounds(c)
	binaryMask := buildBinaryMask(c, len(coefficients))
	constraints := make([][]float64, 0)

	return BasicUnit{pid, coefficients, bounds, constraints, binaryMask, c, pCap, nCap}
}

func buildCoefficients(c []CriticalPoint, pCap CriticalPoint, nCap CriticalPoint) []float64 {
	// Variable cost is split into continious segments described by critical points. Constraints are formed to allow
	// only one segement to be active at a time (i.e. the line between two critical points).
	coefficients := []float64{}
	for _, cp := range c {
		coefficients = append(coefficients, cp.costPerKwh)
	}

	binaryCoeff := make([]float64, len(c)-1)
	coefficients = append(coefficients, binaryCoeff...)

	coefficients = append(coefficients, pCap.CostPerKWH())
	coefficients = append(coefficients, nCap.CostPerKWH())

	return coefficients
}

func buildBounds(c []CriticalPoint) [][2]float64 {

	bounds := [][2]float64{}
	for i := 0; i < len(c); i++ {
		bounds = append(bounds, [2]float64{0, 1})
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

func (u BasicUnit) RealPositiveCapacity() []float64 {
	return []float64{u.realPositiveCapacity.KW()}
}

func (u BasicUnit) RealNegativeCapacity() []float64 {
	return []float64{u.realNegativeCapacity.KW()}
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

func UnitPositiveCapacityConstraint(u *BasicUnit) []float64 {
	cons := make([]float64, u.ColumnSize())
	for i, cp := range u.CriticalPoints() {
		cons[i] = cp.KW()
	}

	cons[u.RealPositiveCapacityLoc()[0]] = -u.RealPositiveCapacity()[0]

	cons = boundConstraint(cons, math.Inf(-1), 0)
	return cons
}

func UnitNegativeCapacityConstraint(u *BasicUnit) []float64 {
	cons := make([]float64, u.ColumnSize())
	for i, cp := range u.CriticalPoints() {
		cons[i] = cp.KW()
	}

	cons[u.RealNegativeCapacityLoc()[0]] = u.RealNegativeCapacity()[0]

	cons = boundConstraint(cons, 0, math.Inf(1))
	return cons
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

func UnitRealPowerConstraint(u *BasicUnit, setpt float64) []float64 {
	rpl := u.RealPowerLoc()
	cp := u.CriticalPoints()

	c := make([]float64, u.ColumnSize())
	for i, loc := range rpl {
		c[loc] = cp[i].KW()
	}

	return boundConstraint(c, setpt, setpt)
}

// buildConstraintsSegement creates a lower diagonal matrix
func segmentPartitionConstraints(u *BasicUnit) [][]float64 {
	cons := [][]float64{}
	cplen := len(u.criticalPoints)
	for i := range u.CriticalPoints() {
		con := make([]float64, u.ColumnSize())
		con[i] = 1

		if i < cplen-1 {
			con[i+cplen] = -1
		}

		if i > 0 {
			con[i+cplen-1] = -1
		}

		cons = append(cons, boundConstraint(con, math.Inf(-1), 0))
	}

	return cons
}

func segmentInterpolationConstraint(u *BasicUnit) []float64 {
	// segment interpolation constraint
	// [0, 1, 1, 1, 0, 0, 0, 0 1]
	cons := make([]float64, u.ColumnSize())
	for i := range u.CriticalPoints() {
		cons[i] = 1
	}

	return boundConstraint(cons, 1, 1)
}

func segmentActivationConstraint(u *BasicUnit) []float64 {

	// segment activation constraint
	// [0, 0, 0, 0, 1, 1, 0, 0 1]
	cons := make([]float64, u.ColumnSize())
	cplen := len(u.CriticalPoints())
	for i := cplen; i < cplen*2-1; i++ {
		cons[i] = 1
	}

	return boundConstraint(cons, 1, 1)
}

func UnitSegmentConstraints(u *BasicUnit) [][]float64 {

	cons := segmentPartitionConstraints(u)
	cons = append(cons, segmentInterpolationConstraint(u))
	cons = append(cons, segmentActivationConstraint(u))

	return cons
}
