/*
Copyright © 2025 Zingui Fred Mike <mikezingui@yahoo.com>
*/
package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/MiltonJ23/Fako/internal/adapters/secondary/network"
	"github.com/MiltonJ23/Fako/internal/core/services"
	"github.com/spf13/cobra"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply [file]",
	Short: "Apply configuration to the network ",
	Args:  cobra.ExactArgs(1),
	Run:   runApply,
}

func init() {
	rootCmd.AddCommand(applyCmd)
}

func runApply(cmd *cobra.Command, args []string) {
	filename := args[0]         // we get the name of the file
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM) // this will help us make sure to intercept the sigterm signal and gracefully shutdown 
	defer cancel() // always make sure to clean up

	go func() {
		<-ctx.Done() // we are waiting for that interruption signal
		fmt.Println("-> Signal Received !! Shutting down gracefully......")
	}()

	data, ReadingFileError := os.ReadFile(filename)
	if ReadingFileError != nil {
		fmt.Printf("Error Reading file %v : %v", filename, ReadingFileError)
		os.Exit(1)
	}
	// now let's pass the file's data to the service layer
	graph, parsingError := services.ParseAndValidate(data)
	if parsingError != nil {
		fmt.Printf("Validation failed for file %v : %v", filename, parsingError)
		os.Exit(1)
	}
	// reaching here means that the file was parsed successfully
	fmt.Println(" -> Validation successful ! No cycles detected")
	fmt.Printf("Found %d resources \n", len(graph.Nodes))

	// now let's determine the execution order
	executionList, SortingError := graph.TopologicalSort()
	if SortingError != nil {
		fmt.Printf("An error occured during sort : %v", SortingError)
		os.Exit(1)
	}

	driver := network.NewMockDriver()

	if err := services.Enforce(ctx, executionList, driver); err != nil {
		fmt.Printf("-> Apply Failed: %v\n", err)
		os.Exit(1)
	}

}
