package user

const Underage = "underage"
const Adult = "adult"
const Senior = "senior"

const MinAge = 0
const MaxAge = 120

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
