package user

const Underage = "underage"
const Adult = "adult"
const Senior = "senior"

type AgeRange struct {
	Min, Max int
}

var categories = map[string]AgeRange{
	Underage: {0, 17},
	Adult:    {18, 63},
	Senior:   {64, 120},
}

func (ar AgeRange) Contains(age int) bool {
	return ar.Min <= age && age <= ar.Max
}
