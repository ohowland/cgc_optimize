package cpl_adapter

import (
	"github.com/lanl/clp"
)

func Solve(w LinearProgram) []float64 {
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
