package adapter

import (
	"github.com/lanl/clp"
	opt "github.com/ohowland/cgc_optimize"
)

func Solve(w opt.LinearProgram) []float64 {
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
