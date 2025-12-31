/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"fmt"
	"os"

	"github.com/MiltonJ23/Fako/internal/core/services"
	"github.com/spf13/cobra"
)

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan [file]",
	Short: "Preview the execution graph",
	Long:  `Reads the YAML intent and displays the dependency graph and execution order`,
	Args:  cobra.ExactArgs(1), // This will help us validate the input by ensuring the arguments are not more than 1
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("plan called")
	},
}

func init() {
	rootCmd.AddCommand(planCmd)

}
func runPlan(cmd *cobra.Command, args []string) {

	filename := args[0]
	fmt.Printf("Reading intent from file %s", filename)
	// now let's read the file data
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
	// we reached here means that all the check passed and now let's print the execution order

	fmt.Println("\n -------------------------- Execution Order -----------------------------\n ")
	for i, res := range executionList {
		fmt.Printf("%d. [%s] Create %s\n", i+1, res.Kind, res.ID)
	}
}
