package random

func (s stdGenerator) UInt64() uint64 {
	return s.seed.Uint64()
}
