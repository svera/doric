// Package mocks implement mocks structs to be used in tests
package mocks

// Randomizer implements the Randomizer interface to generate random integers
type Randomizer struct {
	Values  []int
	current int
}

// Intn return numbers set in the Values property in the same order
// If all numbers inside Values are returned, the slice is ran again from the beginning
func (m *Randomizer) Intn(n int) int {
	if m.current == len(m.Values) {
		m.current = 0
	}
	val := m.Values[m.current]
	m.current++
	return val
}
