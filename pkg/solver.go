package gocoinor

import (
	"github.com/lanl/clp"
	"github.com/ohowland/highs"
)

type MipLinearProgram interface {
	LinearProgram
	Integrality() []int
}

type LinearProgram interface {
	CostCoefficients() []float64
	Bounds() [][2]float64
	Constraints() [][]float64
}

func HighsLpSolve(w LinearProgram) []float64 {
	s, err := highs.New(
		w.CostCoefficients(),
		w.Bounds(),
		w.Constraints(),
		[]int{})

	if err != nil {
		panic(err)
	}

	s.SetObjectiveSense(highs.Minimize)
	s.RunSolver()
	return s.PrimalColumnSolution()
}

func HighsMipSolve(w MipLinearProgram) []float64 {
	s, err := highs.New(
		w.CostCoefficients(),
		w.Bounds(),
		w.Constraints(),
		w.Integrality())

	if err != nil {
		panic(err)
	}

	s.SetObjectiveSense(highs.Minimize)
	s.RunSolver()
	return s.PrimalColumnSolution()
}

func ClpSolve(w LinearProgram) []float64 {
	s := clp.NewSimplex()
	s.EasyLoadDenseProblem(
		w.CostCoefficients(),
		w.Bounds(),
		w.Constraints(),
	)

	s.SetOptimizationDirection(clp.Minimize)
	s.Primal(clp.NoValuesPass, clp.NoStartFinishOptions)
	return s.PrimalColumnSolution()
}
