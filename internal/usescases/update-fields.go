package usescases

import (
	"errors"
	"fmt"

	"github.com/prev-updater/internal/model"
)

type AdoRepository interface {
	GetPipelineRuns(pipelineId int) ([]model.PipelineRuns, error)
	GetPipelineRun(pipelineId, runId int) (*model.PipelineRuns, error)
	GetBuildWorkItem(buildId int) ([]model.BuildWorkItems, error)
	GetWorkitem(workItemId int) (*model.BuildWorkItems, error)
	UpdateWorkitemField(workItemId string, operation model.OperationFields) error
}

type AdoUsesCases struct {
	Repository AdoRepository
}

func NewAdoUsesCases(adoRepository AdoRepository) *AdoUsesCases {
	return &AdoUsesCases{
		Repository: adoRepository,
	}
}

func (u *AdoUsesCases) UpdateFieldsByLastRuns(pipelineId int) error {
	adoRep := u.Repository
	result, err := adoRep.GetPipelineRuns(pipelineId)
	if err != nil {
		return err
	} else if len(result) == 0 {
		return errors.New("No data")
	}

	lastRuns := result[0]

	//Get all WorkItems
	workItems, err := adoRep.GetBuildWorkItem(lastRuns.Id)
	if err != nil {
		return err
	} else if len(workItems) == 0 {
		return errors.New("No data")
	}

	for _, workItem := range workItems {
		err = u.updateFields(workItem.Id, lastRuns.Name)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func (u *AdoUsesCases) UpdateFieldsByPipelineId(pipelineId int) error {
	return nil
}

// Test
func (u *AdoUsesCases) updateFields(woritemId, name string) error {
	repo := u.Repository
	modelToUpdload := model.OperationFields{
		Op:    "add",
		Path:  "/fields/Custom.c14cc8ed-7be8-4c1a-92b3-ebe7f8923d18",
		Value: name,
	}

	return repo.UpdateWorkitemField(woritemId, modelToUpdload)
}
