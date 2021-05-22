package cgc_optimize

type LinearProgram interface {
	CostCoefficients() []float64
	Bounds() [][2]float64
	Constraints() [][]float64
}

type MipLinearProgram interface {
	LinearProgram
	Integrality() []int
}