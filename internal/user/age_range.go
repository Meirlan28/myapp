package user

const Junior = "junior"
const Middle = "middle"
const Senior = "senior"

type AgeRange struct {
	Min, Max int
}

var categories = map[string]AgeRange{
	Junior: {0, 18},
	Middle: {18, 64},
	Senior: {64, 121},
}

func (ar AgeRange) Contains(age int) bool {
	if ar.Min <= age && age < ar.Max {
		return true
	} else {
		return false
	}
}
