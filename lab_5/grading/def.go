package grading

type asp struct {
	alt   int
	score float32
}

func AscendingSort() func(a, b asp) int {
	return func(a, b asp) int {
		if a.score < b.score {
			return -1
		} else if a.score > b.score {
			return 1
		}
		return 0
	}
}

func DescendingSort() func(a, b asp) int {
	return func(a, b asp) int {
		if a.score < b.score {
			return 1
		} else if a.score > b.score {
			return -1
		}
		return 0
	}
}
