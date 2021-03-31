package main

import (
	"fmt"
	"math"

	"github.com/lanl/clp"
)

func main() {
	pinf := math.Inf(1)
	ninf := math.Inf(-1)
	simp := clp.NewSimplex()
	simp.EasyLoadDenseProblem(
		[]float64{1, 1, 1},
		[][2]float64{
			{1, 6},
			{1, 6},
			{1, 6},
		},
		[][]float64{
			{1, 1, -1, 0, pinf},
			{1, 0, 1, -1, pinf},
			{ninf, 1, -2, 1, -1},
		})
	simp.SetOptimizationDirection(clp.Maximize)

	simp.Primal(clp.NoValuesPass, clp.NoStartFinishOptions)
	soln := simp.PrimalColumnSolution()

	fmt.Printf("Die 1 = %.0f\n", soln[0])
	fmt.Printf("Die 2 = %.0f\n", soln[1])
	fmt.Printf("Die 3 = %.0f\n", soln[2])
}
