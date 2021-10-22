package adapter

import (
	opt "github.com/ohowland/cgc_optimize"
	"github.com/ohowland/highs"
)

func SolveLp(w opt.LinearProgram) []float64 {
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

func SolveMip(w opt.MipLinearProgram) []float64 {

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
