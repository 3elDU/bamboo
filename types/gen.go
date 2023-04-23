package types

type WorldGenerator interface {
	// Request chunk generation at those coordinates
	Generate(chunk Chunk)
	// Generates a dummy chunk ( just fills it with water ).
	// Immedinately writes changes to the chunk, skipping queue
	GenerateDummy(chunk Chunk)
	// Returns generated chunks, if any
	Receive() []Chunk
}
