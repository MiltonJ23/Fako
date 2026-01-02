/*
Copyright © 2026 Zingui Fred Mike mikezingui@yahoo.com
*/
package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/MiltonJ23/Fako/internal/adapters/secondary/network"
	"github.com/MiltonJ23/Fako/internal/adapters/secondary/persistence"
	"github.com/MiltonJ23/Fako/internal/core/services"
	"github.com/spf13/cobra"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy [file]",
	Short: "Destroy resources defined in the file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("destroy called")
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
}

func runDestroy(cmd *cobra.Command, args []string) {
	filename := args[0]
	ctx := context.Background()

	data, ReadingFileError := os.ReadFile(filename)
	if ReadingFileError != nil {
		fmt.Printf("-> an error occured when reading the file : %s", ReadingFileError)
		os.Exit(1)
	}

	graph, ParsingError := services.ParseAndValidate(data)
	if ParsingError != nil {
		fmt.Printf("-> an error occured when parsing and validating the file: %s", ParsingError)
		os.Exit(1)
	}
	destructionExecutionList, ReversedTopologicalSortError := graph.ReverseTopologicalSort()
	if ReversedTopologicalSortError != nil {
		fmt.Printf("-> an error happened when trying to get the reverse execution order list: %s", ReversedTopologicalSortError)
		os.Exit(1)
	}

	driver, DriverSelectionError := network.GetDriver(driverType)
	if DriverSelectionError != nil {
		fmt.Printf("-> an error occured when selecting the proper driver : %s", DriverSelectionError)
		os.Exit(1)
	}
	repo := persistence.NewJSONStateRepository("fako.state.json")
	DestroyingError := services.Destroy(ctx, destructionExecutionList, driver, repo)
	if DestroyingError != nil {
		fmt.Printf("-> Destroy Failed .....: %s", DestroyingError)
		os.Exit(1)
	}
}
