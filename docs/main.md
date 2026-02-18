# Main Entry Point

`main.go` is the entry point for the flavour-gen CLI application.

## Overview

This file initializes the root command and executes the CLI.

## Code

```go
package main

import (
	"github.com/ZenDeveloper7/flavour-gen/cmd"
)

func main() {
	cmd.Execute()
}
```

## Flow

1. `main()` calls `cmd.Execute()`
2. `cmd` package initializes Cobra root command
3. Root command sets up color output
4. Executes the appropriate subcommand (e.g., `create`)

## Related Files

- `cmd/root.go` - Root command setup
- `cmd/create.go` - Create subcommand
