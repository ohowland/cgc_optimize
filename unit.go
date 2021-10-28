package cgc_optimize

import (
	"github.com/google/uuid"
)

type Unit interface {
	PID() uuid.UUID
	CostCoefficients() []float64
	Bounds() [][2]float64
	Constraints() [][]float64
	ColumnSize() int
	Integrality() []int

	CriticalPoints() []CriticalPoint
	RealPositiveCapacity() []float64
	RealNegativeCapacity() []float64

	RealPowerLoc() []int
	RealPositiveCapacityLoc() []int
	RealNegativeCapacityLoc() []int
	//StoredEnergyLoc() []int
}

type EnergyStorageUnit interface {
	Unit
	StoredEnery() []float64
	StoredEnergyLoc() []int
}
