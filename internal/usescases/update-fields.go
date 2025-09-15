package usescases

import (
	"github.com/Damien-Venant/prev-updater/internal/model"
	"github.com/rs/zerolog"
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
	Logger     *zerolog.Logger
}

func NewAdoUsesCases(adoRepository AdoRepository, logger *zerolog.Logger) *AdoUsesCases {
	return &AdoUsesCases{
		Repository: adoRepository,
		Logger:     logger,
	}
}

func (u *AdoUsesCases) UpdateFieldsByLastRuns(pipelineId int) error {
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

	lastRuns := result[0]
	u.Logger.
		Info().
		Dict("metadata", zerolog.Dict().Int("pipeline-id", pipelineId)).
		Msgf("Last run Id : %d", lastRuns.Id)
	//Get all WorkItems
	workItems, err := adoRep.GetBuildWorkItem(lastRuns.Id)
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
		err = u.updateFields(workItem.Id, lastRuns.Name)
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

func (u *AdoUsesCases) updateFields(woritemId, name string) error {
	repo := u.Repository
	modelToUpdload := model.OperationFields{
		Op:    "add",
		Path:  "/fields/Custom.c14cc8ed-7be8-4c1a-92b3-ebe7f8923d18",
		Value: name,
	}

	return repo.UpdateWorkitemField(woritemId, modelToUpdload)
}
