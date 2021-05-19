package gocoinor

import (
	"github.com/lanl/clp"
)

type LinearProgram interface {
	CostCoefficients() []float64
	Bounds() [][2]float64
	Constraints() [][]float64
}

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
