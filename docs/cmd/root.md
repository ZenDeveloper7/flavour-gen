# Root Command

`cmd/root.go` defines the root command for the flavour-gen CLI.

## Overview

Sets up the Cobra command structure and enables color output.

## Code

```go
var rootCmd = &cobra.Command{
	Use:   "flavor-gen",
	Short: "Flavor Generator CLI - Build Android app flavors",
	Long:  `A tool to generate complete Android app flavor packages...`,
}

func Execute() {
	color.NoColor = false
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
```

## Configuration

- **Use:** `flavor-gen`
- **Short:** Brief description
- **Long:** Detailed description

## Color Output

Color is enabled by default. Set `color.NoColor = true` to disable.

## Related Files

- `cmd/create.go` - Create subcommand
- `main.go` - Entry point
