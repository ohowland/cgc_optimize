package esslp

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSimpleLPFormulation(t *testing.T) {
	pid1, _ := uuid.NewUUID()
	pid2, _ := uuid.NewUUID()
	a1 := NewUnit(pid1, 1.0, 2.0, 3.0, 4.0)
	a2 := NewUnit(pid2, 5.0, 6.0, 7.0, 8.0)
	ag1 := NewGroup(a1, a2)

	ag1.NewConstraint([]float64{1, 1, 0, 0, 1, 1, 0, 0}, 10, 10)

	//fmt.Println(lp.group.CostCoefficients())
	//fmt.Println(lp.group.Bounds())
	//fmt.Println(lp.group.Constraints())
	sol := Solve(ag1)
	assert.Equal(t, []float64{10, 0, 0, 0, 0, 0, 0, 0}, sol)
}
