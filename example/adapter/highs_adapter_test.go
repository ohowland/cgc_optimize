package adapter

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/google/uuid"
	opt "github.com/ohowland/cgc_optimize"
	"github.com/stretchr/testify/assert"
)

func TestHighsEssLpNetLoadConstraint(t *testing.T) {
	pid1, _ := uuid.NewUUID()
	pid2, _ := uuid.NewUUID()
	a1 := opt.NewBasicUnit(pid1, 1.0, 2.0, 0.01, 0, 5, 5, 5, 5)
	a2 := opt.NewBasicUnit(pid2, 5.0, 6.0, 0.01, 0, 10, 10, 10, 10)
	ag1 := opt.NewGroup(a1, a2)

	nlc := opt.NetLoadConstraint(&ag1, 10)
	ag1.NewConstraint(nlc)

	sol := SolveLp(ag1)
	assert.InDeltaSlice(t, []float64{5, 0, 0, 0, 5, 0, 0, 0}, sol, 0.1)
}

func TestHighsPiecewiseSingleAsset(t *testing.T) {
	pid1, _ := uuid.NewUUID()
	a1 := opt.NewPiecewiseUnit(pid1, []opt.CriticalPoint{opt.NewCriticalPoint(-10, 1), opt.NewCriticalPoint(0, 0), opt.NewCriticalPoint(10, 1)})
	ag1 := opt.NewGroup(a1)

	nlc := opt.NetLoadPiecewiseConstraint(&ag1, 5)
	ag1.NewConstraint(nlc)

	fmt.Println("cost coeff:", ag1.CostCoefficients())
	fmt.Println("crit pts:", ag1.CriticalPoints())
	fmt.Println("binary mask:", ag1.Integrality())
	fmt.Println("bounds:", ag1.Bounds())
	fmt.Println("constraints:", ag1.Constraints())

	sol := SolveMip(ag1)

	fmt.Println("solution:", sol)
	assert.InDeltaSlice(t, []float64{0, 0.5, 0.5, 0, 1, 0}, sol, 0.1)
}

func TestHighsPiecewiseTwoAssets(t *testing.T) {
	pid1, _ := uuid.NewUUID()
	a1 := opt.NewPiecewiseUnit(pid1, []opt.CriticalPoint{opt.NewCriticalPoint(-10, 3), opt.NewCriticalPoint(0, 0), opt.NewCriticalPoint(10, 3)})

	pid2, _ := uuid.NewUUID()
	a2 := opt.NewPiecewiseUnit(pid2, []opt.CriticalPoint{opt.NewCriticalPoint(-5, 1), opt.NewCriticalPoint(0, 0), opt.NewCriticalPoint(5, 1)})

	ag1 := opt.NewGroup(a1, a2)
	nlc := opt.NetLoadPiecewiseConstraint(&ag1, 10)
	ag1.NewConstraint(nlc)

	sol := SolveMip(ag1)
	fmt.Println("solution:", sol)
	assert.InDeltaSlice(t, []float64{0, 0.5, 0.5, 0, 1, 0, 0, 0, 1, 0, 1, 0}, sol, 0.1)
}

func TestHighsEssLpAssetCapacityConstraint(t *testing.T) {
	pid1, _ := uuid.NewUUID()
	pid2, _ := uuid.NewUUID()
	a1 := opt.NewBasicUnit(pid1, 1.0, 2.0, 0.01, 0, 5, 5, 5, 5)
	a1.NewConstraint(opt.BasicUnitCapacityConstraints(&a1)...)
	a2 := opt.NewBasicUnit(pid2, 5.0, 6.0, 0.01, 0, 10, 10, 10, 10)
	a2.NewConstraint(opt.BasicUnitCapacityConstraints(&a2)...)

	ag1 := opt.NewGroup(a1, a2)

	nlc := opt.NetLoadConstraint(&ag1, 10)
	ag1.NewConstraint(nlc)

	sol := SolveLp(ag1)
	assert.InDeltaSlice(t, []float64{5, 0, 5, 0, 5, 0, 5, 0}, sol, 0.1)
}

func TestHighsEssLpGroupCapacityConstraint(t *testing.T) {
	pid1, _ := uuid.NewUUID()
	pid2, _ := uuid.NewUUID()
	a1 := opt.NewBasicUnit(pid1, 1.0, 2.0, 0.01, 0, 5, 5, 5, 5)
	err := a1.NewConstraint(opt.BasicUnitCapacityConstraints(&a1)...)
	assert.Nil(t, err)
	a2 := opt.NewBasicUnit(pid2, 5.0, 6.0, 0.01, 0, 10, 10, 10, 10)
	err = a2.NewConstraint(opt.BasicUnitCapacityConstraints(&a2)...)
	assert.Nil(t, err)

	ag1 := opt.NewGroup(a1, a2)
	err = ag1.NewConstraint(opt.NetLoadConstraint(&ag1, 7), opt.GroupPositiveCapacityConstraint(&ag1, 10))
	assert.Nil(t, err)

	sol := SolveLp(ag1)
	//fmt.Println(sol)
	assert.InDeltaSlice(t, []float64{5, 0, 5, 0, 2, 0, 5, 0}, sol, 0.1)
}
func TestHighsEssLpClusterLinkedBusConstraint(t *testing.T) {
	pid1, _ := uuid.NewUUID()
	pid2, _ := uuid.NewUUID()
	a1 := opt.NewBasicUnit(pid1, 0.1, 0.1, 0.01, 0, 5, 5, 5, 0)
	a2 := opt.NewBasicUnit(pid2, 2.0, 2.0, 0.01, 0, 5, 5, 5, 0)
	err := a1.NewConstraint(opt.BasicUnitCapacityConstraints(&a1)...)
	assert.Nil(t, err)
	err = a2.NewConstraint(opt.BasicUnitCapacityConstraints(&a2)...)
	assert.Nil(t, err)

	ag1 := opt.NewGroup(a1, a2)
	ag2 := opt.NewGroup(a1)
	nload := (5 * rand.Float64()) + 5
	err = ag1.NewConstraint(opt.NetLoadConstraint(&ag1, nload))
	assert.Nil(t, err)

	cl1 := opt.NewCluster(ag1, ag2)
	err = cl1.NewConstraint(opt.LinkedBusConstraints(&cl1, pid1)...)
	assert.Nil(t, err)

	sol := SolveLp(ag1)

	assert.InDeltaSlice(t, []float64{5, 0, 5, 0, nload - 5, 0, nload - 5, 0}, sol, 0.1)
	//fmt.Println(nload, nload-5, sol)
}
func TestHighsEssLpSeriesDischargeBatteryConstraint(t *testing.T) {
	pid1, _ := uuid.NewUUID()
	a1 := opt.NewBasicUnit(pid1, 0.1, 0.1, 0.01, 0, 10, 10, 10, 20)
	err := a1.NewConstraint(opt.BasicUnitCapacityConstraints(&a1)...)
	assert.Nil(t, err)
	ag1 := opt.NewGroup(a1)

	nload := 10.0
	err = ag1.NewConstraint(opt.NetLoadConstraint(&ag1, nload))
	assert.Nil(t, err)

	s1 := opt.NewSeries(ag1, ag1, ag1, ag1)
	err = s1.NewConstraint(opt.BatteryInitialEnergyConstraint(&s1, pid1, 20))
	assert.Nil(t, err)
	err = s1.NewConstraint(opt.BatteryEnergyConstraint(&s1, pid1, 0.5)...)
	assert.Nil(t, err)

	sol := SolveLp(s1)
	assert.InDeltaSlice(t, []float64{10, 0, 10, 20, 10, 0, 10, 15, 10, 0, 10, 10, 10, 0, 10, 5}, sol, 0.1, "battery positive power not decreasing stored energy")
}
func TestHighsEssLpSeriesChargeBatteryConstraint(t *testing.T) {
	pid1, _ := uuid.NewUUID()
	a1 := opt.NewBasicUnit(pid1, 0.1, 0.1, 0.01, 0, 10, 10, 10, 20)
	err := a1.NewConstraint(opt.BasicUnitCapacityConstraints(&a1)...)
	assert.Nil(t, err)
	ag1 := opt.NewGroup(a1)

	nload := -10.0
	err = ag1.NewConstraint(opt.NetLoadConstraint(&ag1, nload))
	assert.Nil(t, err)

	s1 := opt.NewSeries(ag1, ag1, ag1, ag1)
	err = s1.NewConstraint(opt.BatteryInitialEnergyConstraint(&s1, pid1, 5))
	assert.Nil(t, err)
	err = s1.NewConstraint(opt.BatteryEnergyConstraint(&s1, pid1, 0.5)...)
	assert.Nil(t, err)

	sol := SolveLp(s1)
	assert.InDeltaSlice(t, []float64{0, 10, 10, 5, 0, 10, 10, 10, 0, 10, 10, 15, 0, 10, 10, 20}, sol, 0.1, "battery negative power not increasing stored energy")
}
