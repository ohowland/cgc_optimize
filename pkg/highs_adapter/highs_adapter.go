package highs_adapter

import (
	"github.com/ohowland/highs"
)

func SolveLp(w LinearProgram) []float64 {
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

func SolveMip(w MipLinearProgram) []float64 {
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
