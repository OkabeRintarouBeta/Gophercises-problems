## How to run
1. Install all packages required: `go mod tidy`
2. Compiles the Go package in the current directory, along with all its dependencies: `go install .`
3. Run the task manager in the command line:
   1. `task --help`
   2. `task add <task name>`
   3. `task do <task_id>`
   4. `task list`

## Troubleshoot
If `task` command is not found after step 2, make sure that GOBIN and GOPATH are set correctly.
