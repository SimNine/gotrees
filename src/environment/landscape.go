package environment

type Landscape struct {
	groundBaseline     int
	groundDegree       int
	groundFrequency    []float32
	groundAmplitude    []float32
	groundDisplacement []float32

	tileType     [][]int
	groundLevels []int
}
