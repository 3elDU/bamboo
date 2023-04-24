package types

type WorldGenerator interface {
	// Request chunk generation at those coordinates
	Generate(chunk Chunk)
	// Generates a dummy chunk
	// Immedinately writes changes to the chunk, skipping queue
	GenerateDummy(chunk Chunk)
	// Returns generated chunks, if any
	Receive() []Chunk

	// Runs a generator main loop.
	// Must be called with `go Run()`
	Run()
}
