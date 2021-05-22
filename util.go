package cgc_optimize
// lb returns the lower bounds of a constraint
func lb(c []float64) float64 {
	return c[0]
}

// ub returns the upper bounds of a constraint
func ub(c []float64) float64 {
	return c[len(c)-1]
}

// cons returns the constraint without the upper and lower bounds
func cons(c []float64) []float64 {
	return c[1 : len(c)-1]
}

// returns a bounded constraint []float64{lb, cons..., ub}
func boundConstraint(cons []float64, lb float64, ub float64) []float64 {
	return append(append([]float64{lb}, cons...), ub)
}
