package games

func CoinGameWith(same, diff int) [][]int {
	return [][]int{
		{same, diff},
		{diff, same},
	}
}

func CoinGame() [][]int {
	return CoinGameWith(2, -3)
}

// Chicken (Hawk-Dove) / payoff to row player
// Numeric default Chicken (Hawk,Dove) with a negative HH entry:
// H vs H = -2, H vs D = 2, D vs H = 0, D vs D = 1
func ChickenDefault() [][]int {
	return [][]int{
		{-2, 2},
		{0, 1},
	}
}

func PrisonersDilemma() [][]int {
	// Row/Col: Cooperate (C), Defect (D)
	// Typical payoffs for row player: R=3 (C,C), S=0 (C,D), T=5 (D,C), P=1 (D,D)
	// Matrix:
	//   C  D
	// C 3  0
	// D 5  1
	return [][]int{
		{3, 0},
		{5, 1},
	}
}

func Game2_s_vals() [][]int {
	// s,t in {-1,0,1}, payoff = s*(t-s) + t*(t+s)
	values := []int{-1, 0, 1}
	m := make([][]int, 3)
	for i, s := range values {
		row := make([]int, 3)
		for j, t := range values {
			row[j] = s*(t-s) + t*(t+s)
		}
		m[i] = row
	}
	return m
}

// 0 tie, +1 win for row, -1 loss
// order: Rock, Paper, Scissors
func RPS() [][]int {
	return [][]int{
		{0, -1, 1},
		{1, 0, -1},
		{-1, 1, 0},
	}
}

// strategies: (show, guess) in {1..fingers} -> fingers^2 x fingers^2
// payoff to player1:
// if p1_guess == p2_show && p2_guess != p1_show -> + (p1_show + p2_show)
// if p2_guess == p1_show && p1_guess != p2_show -> - (p1_show + p2_show)
// if both guess correctly -> 0 (payments cancel)
func Morra(fingers int) [][]int {
	strats := make([]int, fingers)
	for i := 0; i < fingers; i++ {
		strats[i] = i + 1
	}
	n := fingers * fingers
	m := make([][]int, n)
	for i := 0; i < n; i++ {
		m[i] = make([]int, n)
	}

	find_idx := func(show, guess int) int {
		return (show-1)*fingers + (guess - 1)
	}

	for _, p1_show := range strats {
		for _, p1_guess := range strats {
			for _, p2_show := range strats {
				for _, p2_guess := range strats {
					i := find_idx(p1_show, p1_guess)
					j := find_idx(p2_show, p2_guess)

					payoff := p1_show + p2_show
					p1_correct := (p1_guess == p2_show)
					p2_correct := (p2_guess == p1_show)

					multiplier := 0
					if p1_correct && !p2_correct {
						multiplier = 1
					} else if p2_correct && !p1_correct {
						multiplier = -1
					}

					m[i][j] = payoff * multiplier
				}
			}
		}
	}
	return m
}

func Game6(k int) [][]int {
	// choices 1..k, payoff: if i>=j -> i-j, else -> -(i+j)
	m := make([][]int, k)
	for i := 0; i < k; i++ {
		m[i] = make([]int, k)
		for j := 0; j < k; j++ {
			ii := i + 1
			jj := j + 1
			if ii >= jj {
				m[i][j] = ii - jj
			} else {
				m[i][j] = -(ii + jj)
			}
		}
	}
	return m
}

func Blotto(attacker, defender int) [][]int {
	// allocations between two positions. rows: attacker split a in [0..attacker] (a in pos1)
	// columns: defender split d in [0..defender]
	// attacker wins (payoff +1) if a>c OR (attacker_pos2 > defender_pos2), else payoff -1
	rows := attacker + 1
	cols := defender + 1
	m := make([][]int, rows)
	for a := 0; a <= attacker; a++ {
		m[a] = make([]int, cols)
		for c := 0; c <= defender; c++ {
			att_pos1 := a
			att_pos2 := attacker - a
			def_pos1 := c
			def_pos2 := defender - c
			if att_pos1 > def_pos1 || att_pos2 > def_pos2 {
				m[a][c] = 1
			} else {
				m[a][c] = -1
			}
		}
	}
	return m
}

func SellerProblem(k, a, b, alpha, beta int) [][]int {
	// seller chooses s in [0..k]; demand d in [alpha..beta]
	// payoff = a*min(s,d) - b*max(0, s-d)
	rows := k + 1
	cols := beta - alpha + 1
	m := make([][]int, rows)
	for s := 0; s <= k; s++ {
		m[s] = make([]int, cols)
		for d := alpha; d <= beta; d++ {
			j := d - alpha
			if d >= s {
				m[s][j] = a * s
			} else {
				m[s][j] = a*d - b*(s-d) // = d*(a+b) - b*s
			}
		}
	}
	return m
}
