package bankeralgo

import (
	"fmt"
	"sync"
)

type BankerAlgo struct {
	mu             sync.Mutex // mutex for all methods that read/write to shared resources
	numProcesses   int
	totalResources int
	avaiable       int
	max            []int
	allocated      []int
	needed         []int
}

// It is something like the constructor in Java, It will initialize the struct with the given values, make calculations (needed and available) and return the pointer to struc BankerAlgo
func NewBankerAlgo(n int, total_resources int, max []int, alloc []int) (*BankerAlgo, error) {
	// Input validation
	if len(max) != n {
		return nil, fmt.Errorf("Error: max has %d elements, expected %d", len(max), n)
	}
	if len(alloc) != n {
		return nil, fmt.Errorf("Error: alloc has %d elements, expected %d", len(alloc), n)
	}

	ba := &BankerAlgo{
		numProcesses:   n,
		totalResources: total_resources,
		max:            max,
		allocated:      alloc,
		needed:         make([]int, n),
	}
	// Calculate the needed resources and update the variable inside the BankerAlgo struct
	ba.neededRes()

	// Calculate the current available resources
	ba.avaiable = ba.totalResources - sum(ba.allocated)

	if err := ba.checkConditions(); err != nil {
		return nil, err
	}
	return ba, nil
}

// =============== GETTERS ================

func (ba *BankerAlgo) TotalResources() int {
	return ba.totalResources
}

func (ba *BankerAlgo) NumProcesses() int {
	return ba.numProcesses
}

// =============== ALGORITHM ===============

func (ba *BankerAlgo) RequestResources(process_idx, request int) error {
	ba.mu.Lock()
	defer ba.mu.Unlock() // Automatically will unlock aftrer the method ends

	// Input validation
	if process_idx < 0 || process_idx >= ba.numProcesses {
		return fmt.Errorf("Error: invalid process index %d", process_idx)
	}
	if request <= 0 {
		return fmt.Errorf("Error: request must be positive, got %d", request)
	}

	fmt.Printf("==== START ====\n")

	// 1. Check if the Process` request is in range of what process asked for at the beginning
	if request > ba.needed[process_idx] {
		return fmt.Errorf("Error: Process %d asks for more resources [%d] than it needs [%d].", process_idx, request, ba.needed[process_idx])
	}

	// 2. Check if banker has resources available at the moment
	if request > ba.avaiable {
		return fmt.Errorf("Error: Request %d is not available at the moment. Process %d must wait.", request, process_idx)
	}

	// Temporarily allocate resources, reallocate=false states that we are allocating resources (just mathematical operation)
	ba.allocate(process_idx, request, false)
	fmt.Printf("ALLOCATED P_%d\n", process_idx)
	ba.PrintSummary()

	err := ba.safetyCheck()
	if err != nil {
		// Reallocate resources back
		ba.allocate(process_idx, request, true)
		fmt.Printf("ROLLBACK P_%d\n", process_idx)
		ba.PrintSummary()
		return err
	}

	// Release resources if process has already used all its max resorces
	if ba.needed[process_idx] == 0 {
		fmt.Printf("Process %d has all resources, releasing automatically.\n", process_idx)
		ba.release(process_idx)
	}

	fmt.Printf("==== END ====\n")
	ba.PrintSummary()
	return nil
}

// We need to find out whether there exists a sequence where all processes will finish their work and can free up their resources
func (ba *BankerAlgo) safetyCheck() error {
	fmt.Printf("---RUNNING SAFETY CHECK---\n")
	available := ba.avaiable                  // local variable
	finished := make([]bool, ba.numProcesses) // all values are false by default
	processes := make([]int, 0)               // Append process ids of successfully finished processes

	finished_process := true // we need to find at least one process from all processes that can finish with current available resorces

	for finished_process {
		finished_process = false
		// Loop over the processes
		for i := 0; i < ba.numProcesses; i++ {
			if !finished[i] && ba.needed[i] <= available {
				available += ba.allocated[i] // process can finish, simulate: "free" its resources
				finished_process = true
				finished[i] = true // mark process as finished
				processes = append(processes, i)
			}
		}
	}
	// Check if all processes could finish their work
	for i, v := range finished {
		if !v {
			return fmt.Errorf("Error: unsafe state: process %d could not finish. Rollback.", i)
		}
	}
	// Print safe processes
	fmt.Printf("Safe sequence: \n")
	for _, p_idx := range processes {
		fmt.Printf("P_%d, ", p_idx)
	}
	fmt.Println()

	return nil
}

func (ba *BankerAlgo) allocate(process_idx, request int, reallocate bool) {
	if reallocate {
		request = -request
	}
	ba.avaiable -= request
	ba.allocated[process_idx] += request
	ba.needed[process_idx] -= request
}

func (ba *BankerAlgo) ReleaseResources(process_idx int) error {
	ba.mu.Lock()
	defer ba.mu.Unlock()

	if process_idx < 0 || process_idx >= ba.numProcesses {
		return fmt.Errorf("Error: invalid process index %d\n", process_idx)
	}

	if ba.allocated[process_idx] == 0 {
		return fmt.Errorf("Error: process %d has no allocated resources\n", process_idx)
	}

	ba.release(process_idx)
	return nil
}

func (ba *BankerAlgo) release(process_idx int) {
	ba.avaiable += ba.allocated[process_idx]
	ba.allocated[process_idx] = 0
	ba.needed[process_idx] = ba.max[process_idx]
}
