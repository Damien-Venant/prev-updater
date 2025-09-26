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
	defaultRefName, err := adoRep.GetRepositoryById(repositoryId)
	if err != nil {
		u.Logger.Error().
			Err(err).
			Stack().
			Dict("metadata", zerolog.Dict().Int("pipeline-id", pipelineId)).
			Send()
		return err
	}

	var beforeRuns model.PipelineRuns
	lastRuns := result[0]
	if len(result) > 0 {
		beforeRuns = result[1]
	} else {
		beforeRuns = lastRuns
	}

	actualRefName := lastRuns.Resources.Repositories.Self.RefName
	if actualRefName != defaultRefName.DefaultBranch {
		filterRuns := queryslice.Filter(result[1:len(result)-1], func(pre model.PipelineRuns) bool {
			return pre.Resources.Repositories.Self.RefName == actualRefName
		})
		if len(filterRuns) > 0 {
			beforeRuns = filterRuns[0]
		}
	}
	//Get all WorkItems
	workItems, err := adoRep.GetBuildWorkItem(beforeRuns.Id, lastRuns.Id)
	if err != nil {
		u.Logger.
			Error().
			Err(err).
			Stack().
			Dict("metadata", zerolog.Dict().Int("pipeline-id", pipelineId)).
			Send()
		return err
	} else if len(workItems) == 0 {
		u.Logger.
			Warn().
			Dict("metadata", zerolog.Dict().Int("pipeline-id", pipelineId)).
			Msgf("GetBuildWorkitem return no data, pipeline id : %d", lastRuns.Id)
		return nil
	}

	for _, workItem := range workItems {
		err = u.updateFields(workItem.Id, lastRuns.Name, fieldName)
		if err != nil {
			u.Logger.Err(err).Dict("pipeline-id",
				zerolog.Dict().Int("pipeline-id", pipelineId)).
				Send()
		}
	}
	return nil
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
