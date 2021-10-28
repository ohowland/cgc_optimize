package cgc_optimize

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type EssUnit struct {
	bu     BasicUnit
	energy float64
}

// NewEssUnit returns a configured unit struct.
func NewEssUnit(pid uuid.UUID, c []CriticalPoint, pCap CriticalPoint, nCap CriticalPoint, e float64) EssUnit {

	basicUnit := NewBasicUnit(pid, c, pCap, nCap)

	basicUnit.coefficients = append(basicUnit.coefficients, 0)
	basicUnit.bounds = append(basicUnit.bounds, [2]float64{0, 1})
	basicUnit.integrality = append(basicUnit.integrality, 0)

	cons := make([][]float64, 0)
	for _, c := range basicUnit.constraints {
		// insert decision variable into existing constraints
		cons = append(cons, append(append(c[0:len(c)-2], 0), c[len(c)-1]))
	}

	basicUnit.constraints = cons

	return EssUnit{basicUnit, e}
}

func (u EssUnit) PID() uuid.UUID {
	return u.bu.pid
}

func (u EssUnit) CostCoefficients() []float64 {
	return u.bu.coefficients
}

func (u EssUnit) ColumnSize() int {
	return len(u.bu.coefficients)
}

func (u *EssUnit) NewConstraint(t_c ...[]float64) error {
	cx := make([][]float64, 0)
	for _, c := range t_c {
		if len(c) != u.ColumnSize()+2 {
			err := fmt.Sprintf("constraint contains %v columns, expected: %v", len(c), u.ColumnSize()+2)
			return errors.New(err)
		}
		cx = append(cx, c)
	}

	// if no errors: add constraints to unit
	u.bu.constraints = append(u.bu.constraints, cx...)
	return nil
}

func (u EssUnit) Constraints() [][]float64 {
	return u.bu.constraints
}

func (u EssUnit) Bounds() [][2]float64 {
	return u.bu.bounds
}

func (u EssUnit) Integrality() []int {
	return u.bu.integrality
}

func (u EssUnit) CriticalPoints() []CriticalPoint {
	return u.bu.criticalPoints
}

func (u EssUnit) RealPositiveCapacity() []float64 {
	return []float64{u.bu.realPositiveCapacity.KW()}
}

func (u EssUnit) RealNegativeCapacity() []float64 {
	return []float64{u.bu.realNegativeCapacity.KW()}
}

func (u EssUnit) RealPowerLoc() []int {
	locs := make([]int, len(u.bu.criticalPoints))
	for i := range locs {
		locs[i] = i
	}

	return locs
}

func (u EssUnit) RealPositiveCapacityLoc() []int {
	loc := []int{len(u.bu.criticalPoints)*2 - 1}

	return loc
}

func (u EssUnit) RealNegativeCapacityLoc() []int {
	loc := []int{len(u.bu.criticalPoints) * 2}

	return loc
}
func (u EssUnit) StoredEneryLoc() []int {
	loc := []int{len(u.bu.criticalPoints)*2 + 1}

	return loc
}

func (u EssUnit) StoredEnergy() []float64 {
	return []float64{u.energy}
}

// Constraints

func EssUnitRealPowerConstraint(u *EssUnit, setpt float64) []float64 {
	rpl := u.RealPowerLoc()
	cp := u.CriticalPoints()

	c := make([]float64, u.ColumnSize())
	for i, loc := range rpl {
		c[loc] = cp[i].KW()
	}

	return boundConstraint(c, setpt, setpt)
}
