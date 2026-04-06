package bankeralgo

import "fmt"

func (ba *BankerAlgo) checkConditions() error {
	// 1. sum of max > totalResources, this check is from the homework, otherwise it is redundant
	if sum(ba.max) <= ba.totalResources {
		return fmt.Errorf("Error: sum of max resources (%d) must be greater than total resources (%d), otherwise deadlock cannot occur", sum(ba.max), ba.totalResources)
	}
	// 2. For each process 0 <= allocated[i] <= max[i]
	for i := 0; i < ba.numProcesses; i++ {
		if ba.allocated[i] < 0 || ba.allocated[i] > ba.max[i] {
			return fmt.Errorf("Error: For each process: 0 <= allocated[i] <= max[i], Process %d: alloc %d vs max %d", i, ba.allocated[i], ba.max[i])

		}
	}
	return nil
}

// Updates needed array with max[i] - allocated[i] values
func (ba *BankerAlgo) neededRes() {
	for i := 0; i < ba.numProcesses; i++ {
		ba.needed[i] = ba.max[i] - ba.allocated[i]
	}
}

func sum(arr []int) int {
	sum := 0
	for _, n := range arr {
		sum += n
	}
	return sum
}

func (ba *BankerAlgo) PrintSummary() {
	fmt.Println("╔══════════════════════════════╗")
	fmt.Println("║      BANKER ALGORITHM        ║")
	fmt.Println("╚══════════════════════════════╝")
	fmt.Printf("  Total resources:     %d\n", ba.totalResources)
	fmt.Printf("  Available:           %d\n", ba.avaiable)
	fmt.Printf("  Allocated (sum):     %d\n", sum(ba.allocated))
	fmt.Println()
	fmt.Printf("  %-10s %-8s %-8s %-8s\n", "Process", "Max", "Alloc", "Needed")
	fmt.Println("  ----------------------------------------")
	for i := 0; i < ba.numProcesses; i++ {
		fmt.Printf("  P_%-8d %-8d %-8d %-8d\n", i, ba.max[i], ba.allocated[i], ba.needed[i])
	}
	fmt.Println("  ----------------------------------------")
	fmt.Println()
}
