package mocks

type Randomizer struct {
	Values  []int
	current int
}

func (m *Randomizer) Intn(n int) int {
	if m.current == len(m.Values) {
		m.current = 0
	}
	val := m.Values[m.current]
	m.current++
	return val
}
