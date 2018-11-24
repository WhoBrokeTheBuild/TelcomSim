package context

// Update is a context of timing data
type Update struct {
	DeltaTime   float32
	ElapsedTime float64
	TotalTime   float64
}
