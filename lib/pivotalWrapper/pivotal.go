package pivotalWrapper

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/salsita/go-pivotaltracker/v5/pivotal"
)

type Pivotal struct {
	client *pivotal.Client
}

type CycleTimeDetail struct {
	Kind           string `json:"kind"`
	TotalCycleTime int    `json:"total_cycle_time"`
	StartedCount   int    `json:"started_count"`
	FinishedCount  int    `json:"finished_count"`
	DeliveredCount int    `json:"delivered_count"`
	StoryId        int    `json:"story_id"`
}

type StoryDetails struct {
	Kind    string          `json:"kind"`
	Stories []pivotal.Story `json:"stories"`
}

func (object *Pivotal) SetupClient() {
	if object.client == nil {
		apiKey := os.Getenv("PIVOTAL_KEY")
		if apiKey == "" {
			fmt.Fprintf(os.Stderr, "PIVOTAL_KEY environment variable must be set.\n")
			os.Exit(1)
		}
		client := pivotal.NewClient(apiKey)
		object.client = client
		if client == nil {
			log.Fatalf("Failed to create client")
		}
		// [END auth]
	}

}

// ListServices return the list of services in a kubernetes cluster
func (object *Pivotal) ListProjects() ([]string, []*pivotal.Project) {
	object.SetupClient()
	projectList, _, err := object.client.Projects.List()
	if err != nil {
		panic(err.Error())
	}

	projectCount := len(projectList)
	result := make([]*pivotal.Project, 0, projectCount)
	projects := make(map[string]*pivotal.Project)
	projectSort := make([]string, 0, projectCount)

	for _, project := range projectList {
		name := project.Name
		projectSort = append(projectSort, name)
		projects[name] = project
	}

	sort.Strings(projectSort) //sort by key
	for _, sort := range projectSort {
		result = append(result, projects[sort])
	}
	return projectSort, result
}

// Get cycle time details in an iteration
func (object *Pivotal) GetCycleTimeDetails(projectId, iterationId int) ([]*CycleTimeDetail, error) {
	u := fmt.Sprintf("projects/%v/iterations/%v/analytics/cycle_time_details", projectId, iterationId)
	req, err := object.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	var cycleTimeDetail []*CycleTimeDetail
	_, err = object.client.Do(req, &cycleTimeDetail)
	if err != nil {
		return nil, err
	}

	return cycleTimeDetail, nil
}

// Get story details in an iteration
func (object *Pivotal) GetIterationStories(projectId, iterationId int) (*StoryDetails, error) {
	u := fmt.Sprintf("projects/%v/iterations/%v", projectId, iterationId)
	req, err := object.client.NewRequest("GET", u, nil)

	if err != nil {
		return nil, err
	}
	var storyDetails StoryDetails
	_, err = object.client.Do(req, &storyDetails)
	if err != nil {
		return nil, err
	}

	return &storyDetails, nil
}
