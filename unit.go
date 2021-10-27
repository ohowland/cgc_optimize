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

	RealPowerLoc() []int
	RealPositiveCapacityLoc() []int
	RealNegativeCapacityLoc() []int
	//StoredEnergyLoc() []int
}

type EnergyStorageUnit interface {
	Unit
	StoredEnergyLoc() []int
}
