package user

const Underage = "underage"
const Adult = "adult"
const Senior = "senior"

type AgeRange struct {
	Min, Max int
}

var categories = map[string]AgeRange{
	Underage: {0, 18},
	Adult:    {18, 64},
	Senior:   {64, 121},
}

func (ar AgeRange) Contains(age int) bool {
	return ar.Min <= age && age < ar.Max
}
