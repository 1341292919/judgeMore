package utils

const (
	Grade20 = "20"
	Grade21 = "21"
	Grade22 = "22"
	Grade23 = "23"
	Grade24 = "24"
	Grade25 = "25"
)

func IsGradeValid(grade string) bool {
	if grade != Grade20 &&
		grade != Grade21 &&
		grade != Grade22 &&
		grade != Grade23 &&
		grade != Grade24 &&
		grade != Grade25 {
		return false
	}
	return true
}
