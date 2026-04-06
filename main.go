package main

import (
	"banker_algo/bankeralgo"
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Banker's Algorithm Simulator\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  go run . [-n <processes>] <total_resources> <max[n]> <allocated[n]>\n\n")
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  total_resources   Total resources the banker has (positive integer)\n")
		fmt.Fprintf(os.Stderr, "  max[n]            Max resources each process can request (comma-separated)\n")
		fmt.Fprintf(os.Stderr, "  allocated[n]      Currently allocated resources per process (comma-separated)\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		fmt.Fprintf(os.Stderr, "  -n int            Number of processes (default 4)\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  go run . 8 3,2,1,4 1,2,0,3\n")
		fmt.Fprintf(os.Stderr, "  go run . -n 4 8 3,2,1,4 1,2,0,3\n")
		fmt.Fprintf(os.Stderr, "  macOS (zsh) / PowerShell — brackets require quotes:\n")
		fmt.Fprintf(os.Stderr, "  go run . 8 \"[3,2,1,4]\" \"[1,2,0,3]\"\n")
		fmt.Fprintf(os.Stderr, "  go run . -n 4 8 \"[3,2,1,4]\" \"[1,2,0,3]\"\n\n")
		fmt.Fprintf(os.Stderr, "  Linux (bash) / Windows cmd — brackets without quotes:\n")
		fmt.Fprintf(os.Stderr, "  go run . 8 [3,2,1,4] [1,2,0,3]\n")
		fmt.Fprintf(os.Stderr, "  go run . -n 4 8 [3,2,1,4] [1,2,0,3]\n")
		fmt.Fprintf(os.Stderr, "\nNotes:\n")
		fmt.Fprintf(os.Stderr, "  macOS (zsh) / PowerShell:  if using brackets, quotes required: \"[3,2,1,4]\"\n")
		fmt.Fprintf(os.Stderr, "  without brackets no quotes needed on any platform\n")
	}

	n := flag.Int("n", 4, "Number of processes, It determines the length of the arrays, Default value is 4")
	flag.Parse()

	if *n <= 0 {
		fmt.Println("Error: -n has to be positive whole integer")
		os.Exit(1)
	}

	args := flag.Args()
	banker := parseArgs(args, *n)

	// BONUS: Parallel run without user input, uncomment to see
	// runInParallel(banker)

	scanner := bufio.NewScanner(os.Stdin)

	for {
		banker.PrintSummary()
		// load request
		prompt := fmt.Sprintf("Enter request number in range of 1-%d", banker.TotalResources())
		request := readInt(scanner, prompt, 1, banker.TotalResources())

		// Load process id
		prompt = fmt.Sprintf("Enter process id in range of 0-%d", banker.NumProcesses()-1)
		process_idx := readInt(scanner, prompt, 0, banker.NumProcesses()-1)

		// RUN BANKER ALGORITHM
		err := banker.RequestResources(process_idx, request)
		if err != nil {
			fmt.Printf("Result: Process %d is UNSAFE.\n", process_idx)
			fmt.Println(err)
		} else {
			fmt.Printf("Result: Process %d is SAFE. Resources are allocated.\n", process_idx)
		}

		fmt.Print("Continue? (y/n): ")
		scanner.Scan()
		answer := strings.TrimSpace(scanner.Text())

		if answer == "n" {
			fmt.Print("Ending...")
			break
		}
		fmt.Println("==============RUN AGAIN===============")
	}
}

// =============== PARALLEL RUN ================

// Function runInParallel will simulate when processes will ask for resources in parallel and not one by another, the Scanner part from the user is not implemented
func runInParallel(banker *bankeralgo.BankerAlgo) {
	var wg sync.WaitGroup
	for i := 0; i < banker.NumProcesses(); i++ {
		wg.Add(1)
		go func(p_idx int) {
			defer wg.Done()

			// Randomly ask for 2 resources
			request := 2

			err := banker.RequestResources(p_idx, request)
			if err == nil {
				// If process can request the resource, simulte working
				fmt.Printf("---> Process %d is working...\n", p_idx)
				time.Sleep(time.Microsecond * 600)

				// Release resources after the work is completed
				if err := banker.ReleaseResources(p_idx); err != nil {
					fmt.Printf("Error: %v\n", err)
				}

			} else {
				fmt.Printf("Error: Process %d: %v\n", p_idx, err)
			}
		}(i)
	}
	// This will ensure that the program is waiting untill all processes (gorutines-sth like threads in go) will finish their work, without this it would just continue executing next lines of code
	wg.Wait()
	fmt.Println("DONE: All processes have finished their work.")
}

// ================ PARSE INPUT===============

func parseArgs(args []string, n int) *bankeralgo.BankerAlgo {

	if len(args) != 3 {
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  go run . 8 3,2,1,4 1,2,0,3\n")
		fmt.Fprintf(os.Stderr, "  go run . -n 4 8 3,2,1,4 1,2,0,3\n")
		fmt.Fprintf(os.Stderr, "  macOS (zsh) / PowerShell — brackets require quotes:\n")
		fmt.Fprintf(os.Stderr, "  go run . 8 \"[3,2,1,4]\" \"[1,2,0,3]\"\n")
		fmt.Fprintf(os.Stderr, "  go run . -n 4 8 \"[3,2,1,4]\" \"[1,2,0,3]\"\n\n")
		fmt.Fprintf(os.Stderr, "  Linux (bash) / Windows cmd — brackets without quotes:\n")
		fmt.Fprintf(os.Stderr, "  go run . 8 [3,2,1,4] [1,2,0,3]\n")
		fmt.Fprintf(os.Stderr, "  go run . -n 4 8 [3,2,1,4] [1,2,0,3]\n")
		os.Exit(1)
	}

	// Max number of resources the banker posesses
	total_resources, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: total_resources must be a positive integer, got: %s\n", args[0])
		os.Exit(1)
	}

	max, err := parseArr(args[1], n)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	allocated, err := parseArr(args[2], n)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	ba, err := bankeralgo.NewBankerAlgo(n, total_resources, max, allocated)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	return ba
}

func readInt(scanner *bufio.Scanner, prompt string, min, max int) int {
	for {
		fmt.Printf("%s: ", prompt)
		scanner.Scan()
		var n int
		_, err := fmt.Sscan(scanner.Text(), &n)
		if err == nil && n >= min && n <= max {
			return n
		}
		fmt.Printf("Enter integer between %d and %d.\n", min, max)
	}
}

func parseArr(s string, n int) ([]int, error) {
	s = strings.Trim(s, "[]")
	parts := strings.Split(s, ",")
	// Check if the number of total elements in the array is correct (equal to the number of processes 'n')
	if len(parts) != n {
		return nil, fmt.Errorf("Array should contain %d number of elements but it has %d.", n, len(parts))
	}
	arr := make([]int, n)
	for i, p := range parts {
		// 1 by 1 convert string input to positive int inside the array
		val, err := strconv.Atoi(strings.TrimSpace(p))

		if err != nil {
			return nil, fmt.Errorf("%s is not an integer", p)
		}
		if val < 0 {
			return nil, fmt.Errorf("%d cant be negative", val)
		}

		arr[i] = val
	}
	return arr, nil
}
