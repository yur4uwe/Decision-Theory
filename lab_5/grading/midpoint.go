package grading

import (
	"fmt"
	"slices"
	"strconv"
)

type Range struct {
	start float32
	end   float32
}

func PrintRankings(rankings []asp) {
	fmt.Println("\nRankings:")

	slices.SortFunc(rankings, AscendingSort())

	for _, kv := range rankings {
		if kv.alt != -1 {
			fmt.Printf("Alternative %d: %.2f\n", kv.alt, kv.score)
		}
	}
}

func isUnusedAlt(alts []asp, alt int) bool {
	for _, a := range alts {
		if a.alt == alt {
			return false
		}
	}
	return true
}

func Midpoint(n int) {
	alts := alternatives(n)

	var worst, best int

	fmt.Printf("Enter the number of the best alternative (1-%d): ", n)
	_, err := fmt.Scan(&best)
	if err != nil || best <= 0 || best > n {
		fmt.Println("Invalid input for best alternative.")
		return
	}

	alts = removeAlternative(alts, best)

	fmt.Printf("Enter the number of the worst alternative (1-%d): ", n)
	_, err = fmt.Scan(&worst)
	if err != nil || worst <= 0 || worst > n {
		fmt.Println("Invalid input for worst alternative.")
		return
	}

	alts = removeAlternative(alts, worst)

	rankings := make([]asp, 0, n)

	rankings = append(rankings, asp{alt: worst, score: 0.0})
	rankings = append(rankings, asp{alt: best, score: 1.0})

	queue := []Range{{start: 0.0, end: 1.0}}

	var midpoint string
	var midpointInt int

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		start := current.start
		end := current.end

		if len(alts) == 1 && len(queue) == 1 {
			midpointInt = alts[0]
			goto putAlternative
		}

		if len(alts) == 0 {
			break
		}

		PrintRankings(rankings)

		fmt.Println("\nFrom the remaining alternatives ", alts)

	chooseMidpoint:

		fmt.Printf("Enter the number of the midpoint alternative in [%.2f-%.2f] (or 'n' if none): ", start, end)
		_, err = fmt.Scan(&midpoint)
		if err != nil {
			fmt.Println("Invalid input.")
			goto chooseMidpoint
		}

		if midpoint == "n" {
			continue
		}

		midpointInt, err = strconv.Atoi(midpoint)
		if err != nil || midpointInt <= 0 || midpointInt > n || !isUnusedAlt(rankings, midpointInt) || !slices.Contains(alts, midpointInt) {
			fmt.Println("Invalid input for midpoint alternative.")
			goto chooseMidpoint
		}

	putAlternative:

		midpointValue := (start + end) / 2

		rankings = append(rankings, asp{alt: midpointInt, score: midpointValue})

		alts = removeAlternative(alts, midpointInt)

		queue = append(queue, Range{start: start, end: midpointValue})
		queue = append(queue, Range{start: midpointValue, end: end})
	}

	PrintRankings(rankings)
}
