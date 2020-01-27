package doric

// MockRandomizer implements the MockRandomizer interface to generate random integers
type MockRandomizer struct {
	Values  []int
	current int
}

// Intn return numbers set in the Values property in the same order
// If all numbers inside Values were returned, the slice is ran again from the beginning
func (m *MockRandomizer) Intn(n int) int {
	if m.current == len(m.Values) {
		m.current = 0
	}
	val := m.Values[m.current]
	m.current++
	return val
}
