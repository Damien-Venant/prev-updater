package usescases

import (
	"github.com/Damien-Venant/prev-updater/internal/model"
	"github.com/Damien-Venant/prev-updater/pkg/queryslice"
	"github.com/rs/zerolog"
)

type AdoRepository interface {
	GetPipelineRuns(pipelineId int) ([]model.PipelineRuns, error)
	GetPipelineRun(pipelineId, runId int) (*model.PipelineRuns, error)
	GetBuildWorkItem(fromBuildId, toBuildId int) ([]model.BuildWorkItems, error)
	GetWorkitem(workItemId int) (*model.BuildWorkItems, error)
	GetRepositoryById(uuid string) (*model.Repository, error)
	UpdateWorkitemField(workItemId string, operation model.OperationFields) error
}

type AdoUsesCases struct {
	Repository AdoRepository
	Logger     *zerolog.Logger
}

func NewAdoUsesCases(adoRepository AdoRepository, logger *zerolog.Logger) *AdoUsesCases {
	return &AdoUsesCases{
		Repository: adoRepository,
		Logger:     logger,
	}
}

func (u *AdoUsesCases) UpdateFieldsByLastRuns(pipelineId int, repositoryId, fieldName string) error {
	adoRep := u.Repository
	u.Logger.Info().
		Dict("metadata", zerolog.Dict().Int("pipeline-id", pipelineId)).
		Msg("Start update fields by last runs")
	result, err := adoRep.GetPipelineRuns(pipelineId)
	if err != nil {
		u.Logger.Error().
			Err(err).
			Stack().
			Dict("metadata", zerolog.Dict().Int("pipeline-id", pipelineId)).
			Send()
		return err
	} else if len(result) == 0 {
		u.Logger.
			Warn().Msg("GetPipelinesRuns return no data")
		return nil
	}

	//for _, workItem := range workItems {
	//	err = u.updateFields(workItem.Id, lastRuns.Name, fieldName)
	//	if err != nil {
	//		u.Logger.Err(err).Dict("pipeline-id",
	//			zerolog.Dict().Int("pipeline-id", pipelineId)).
	//			Send()
	//	}
	//}
	return nil
}

// getRunsToUpdate is used to return last build and N-1 last build
// It's return a array where the first index is last build and second index is N-1 last build
// If the build have not previous build the N-1 last build is last build on defaultBranch
func (u *AdoUsesCases) getRunsToUpdate(builds []model.PipelineRuns, repositoryId string, pipelineId int) ([]model.PipelineRuns, error) {
	adoRep := u.Repository
	defaultRefName, err := adoRep.GetRepositoryById(repositoryId)
	if err != nil {
		return nil, err
	}
	builds = queryslice.Filter(builds, func(pre model.PipelineRuns) bool {
		return pre.State == "completed"
	})
	buildsOnSameRef := queryslice.Filter(builds, func(pre model.PipelineRuns) bool {
		return pre.Resources.Repositories.Self.RefName == builds[0].Resources.Repositories.Self.RefName
	})

	if len(buildsOnSameRef) > 1 {
		return []model.PipelineRuns{builds[0], builds[1]}, nil
	}

	lastBuildOnDefaultRefName := queryslice.Filter(builds, func(pre model.PipelineRuns) bool {
		return pre.Resources.Repositories.Self.RefName == defaultRefName.DefaultBranch
	})[0]

	return []model.PipelineRuns{builds[0], lastBuildOnDefaultRefName}, nil
}

func (u *AdoUsesCases) UpdateFieldsByPipelineId(pipelineId int) error {
	return nil
}

func (u *AdoUsesCases) updateFields(woritemId, name, fieldName string) error {
	repo := u.Repository
	modelToUpdload := model.OperationFields{
		Op:    "add",
		Path:  fieldName,
		Value: name,
	}

	return repo.UpdateWorkitemField(woritemId, modelToUpdload)
}
