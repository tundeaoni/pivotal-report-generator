// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
)

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "generate-report",
	Short: "generate pivotal tracker iteration report",
	Long:  `generate pivotal tracker iteration report`,
	Run: func(cmd *cobra.Command, args []string) {
		// select project
		fmt.Println("Retrieving list of projects...")
		projectNames, projectDetails := pivotal.ListProjects()
		for i, name := range projectNames {
			fmt.Printf("%d -> %s\n", i+1, name)
		}

		reader := bufio.NewReader(os.Stdin)
		var selectedIndex int
		for {
			fmt.Println("Select desired project:")

			response, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}

			selectedIndex, err = strconv.Atoi(strings.Trim(response, "\n"))
			if selectedIndex > 0 && selectedIndex <= len(projectDetails) {
				break
			} else {
				fmt.Println("- Oh Nooo !!! Please select valid index")
			}
		}
		// selectedProject := projectNames[selectedIndex-1]
		fmt.Println("... ....")

		selectedProjectDetail := projectDetails[selectedIndex-1]
		// select iteration
		currectIteration := selectedProjectDetail.CurrentIterationNumber
		for {
			fmt.Printf("Enter iteration number to generate report.  [1 - %d] :\n", currectIteration)
			response, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}

			selectedIndex, err = strconv.Atoi(strings.Trim(response, "\n"))
			if selectedIndex > 0 && selectedIndex <= currectIteration {
				break
			} else {
				fmt.Printf("- Oh Nooo !!! Please select iteration between 1 and %d \n", currectIteration)
			}
		}
		iteration := selectedIndex
		collatedStories := make(map[string][]map[string]int)
		// get delivered story points https://www.pivotaltracker.com/help/api/rest/v5#Iterations
		allStories, _ := pivotal.GetIterationStories(selectedProjectDetail.Id, iteration)
		for _, v := range allStories.Stories {
			item := map[string]int{
				"random": 12,
			}
			for _, ownerID := range v.OwnerIds {
				id := strconv.Itoa(ownerID)
				collatedStories[id] = append(collatedStories[id], item)
			}
		}
		// get cycle time
		spew.Dump(allStories.Stories)
		// allCycleTime, _ := pivotal.GetCycleTimeDetails(2120784, 25)
		// spew.Dump(len(allCycleTime))
		os.Exit(1)

	},
}

func init() {
	rootCmd.AddCommand(reportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
