package genetree

type GeneTree struct {
	fitness           int
	nutrients         int     // fitness component from soil
	energy            int     // fitness component from sunlight
	fitnessPercentile float32 // fitness as a percentile

	root TreeNode
	age  int
}
