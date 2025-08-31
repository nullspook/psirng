package providers

type RngProvider interface {
	Close()
	ClearBuffer() error
	RandBytes([]byte, int32) error
	RandIntegers([]int32, int32, int32, int32) error
	RandUniform([]float64, int32, float64, float64) error
	RandNormal([]float64, int32, float64, float64) error
}
