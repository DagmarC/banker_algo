# Banker's Algorithm Simulator
 
A Go implementation of the Banker's Algorithm for deadlock avoidance, simulating a banker with 4 clients and a single resource type.
 
## Project Structure
 
```
banker_algo/
├── go.mod
├── main.go
└── bankeralgo/
    ├── banker.go
    └── utils.go
```
 
## Usage

| Parameter | Description |
| :--- | :--- |
| `-n` | (Optional) Number of processes (e.g., `4`). Determines the expected length of the arrays. |
| `totalResources` | Total number of resources the banker has available (positive integer). |
| `max[n]` | Maximum resource demands for each process, comma-separated (e.g., `"3,2,1,4"`). |
| `allocated[n]` | Currently allocated resources for each process, comma-separated (e.g., `"1,2,0,3"`). |


```bash
go run . [-n <processes>] <total_resources> <max[n]> <allocated[n]>
```
### Run on WINDOWS PS
To display usage and documentation
```bash
.\banker_simulator.exe -h
```
```bash
.\banker_simulator.exe -n 4 10 "3,2,1,4" "1,2,0,3"
```
### Run on MacOS

```bash
chmod +x ./banker_simulator_mac
./banker_simulator_mac -h
./banker_simulator_mac -n 4 10 "3,2,1,4" "1,2,0,3"

```
 
### Arguments
 
| Argument | Description |
|---|---|
| `total_resources` | Total resources the banker has (positive integer) |
| `max[n]` | Max resources each process can request (comma-separated) |
| `allocated[n]` | Currently allocated resources per process (comma-separated) |
 
### Flags
 
| Flag | Description | Default |
|---|---|---|
| `-n` | Number of processes | `4` |
 
## Examples
 
### macOS (zsh) / PowerShell — brackets require quotes
```bash
go run . 8 "[3,2,1,4]" "[1,2,0,3]"
go run . -n 4 8 "[3,2,1,4]" "[1,2,0,3]"
```
 
### Linux (bash) / Windows cmd — no quotes needed
```bash
go run . 8 [3,2,1,4] [1,2,0,3]
go run . -n 4 8 [3,2,1,4] [1,2,0,3]
```
 
### Without brackets — works on all platforms
```bash
go run . 8 3,2,1,4 1,2,0,3
go run . -n 4 8 3,2,1,4 1,2,0,3
```
 
## Input Validation
 
The program checks:
- `sum(max) > total_resources` — otherwise deadlock cannot occur
- `0 <= allocated[i] <= max[i]` — for each process
- `request <= need[i]` — process cannot request more than it declared
- `request <= available` — banker must have enough free resources
 
## How It Works
 
1. Program loads `total_resources`, `max[]` and `allocated[]` as arguments
2. Calculates `need[i] = max[i] - allocated[i]` and `available = total_resources - sum(allocated)`
3. In a loop, user selects a process and a request amount
4. Banker's algorithm runs:
   - Validates the request
   - Temporarily allocates resources
   - Runs safety check to find a safe sequence
   - Confirms allocation if safe, rolls back if unsafe
5. When `need[i] == 0`, process releases all resources automatically
 
## Interactive Loop
 
```
Enter process id (0-3): 1
Enter request for process 1 (1-8): 1
  SAFE: process 1 got 1 resources.
  Safe sequence: P0 P1 P2 P3
 
Continue? (y/n): y
```
 
## Requirements
 
- Go 1.21+ (if running directly with go run command, binary does not need it)
