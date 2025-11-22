package main

import (
	"flag"
	"fmt"
	"lab/grading"
)

func main() {
	direct := flag.Bool("d", false, "Use Direct grading method")
	midpoint := flag.Bool("m", false, "Use Midpoint grading method")
	churchmanAckoff := flag.Bool("c", false, "Use Churchman-Ackoff grading method")
	thurstone := flag.Bool("t", false, "Use Thurstone grading method")
	n := flag.Int("n", 0, "Number of alternatives")
	experts := flag.Int("e", 0, "Number of experts (for Thurstone method)")

	flag.Parse()

	usage := `Usage: 
	
	grade [options]

Options:
	-d              Use Direct grading method
	-m              Use Midpoint grading method
	-c              Use Churchman-Ackoff grading method
	-t              Use Thurstone grading method
	-n <number>     Number of alternatives (required)
	-e <number>     Number of experts (required for Thurstone method)

Example:
	grade -d -n 5
	grade -m -n 4
	grade -c -n 6
	grade -t -n 5 -e 10
`

	if flag.NFlag() == 0 {
		fmt.Println(usage)
		return
	}

	if *n <= 1 {
		fmt.Println("Please provide a valid number of alternatives using -n flag (n > 1).")
		return
	}

	fmt.Println(*experts, *n)

	if *direct {
		grading.Direct(*n)
	} else if *midpoint {
		grading.Midpoint(*n)
	} else if *churchmanAckoff {
		grading.ChurchmanAckoff(*n)
	} else if *thurstone && *experts > 0 {
		grading.Thurstone(*n, *experts)
	} else {
		fmt.Println(usage)
	}
}
