package usescases

import (
	"strings"

	"github.com/Damien-Venant/prev-updater/internal/model"
	"github.com/Damien-Venant/prev-updater/pkg/queryslice"
	"github.com/rs/zerolog"
)

type AdoRepository interface {
	GetPipelineRuns(pipelineId int) ([]model.PipelineRuns, error)
	GetPipelineRun(pipelineId, runId int) (*model.PipelineRuns, error)
	GetBuildWorkItem(fromBuildId, toBuildId int) ([]model.BuildWorkItems, error)
	GetWorkItem(workItemId string) (*model.WorkItem, error)
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

	builds, err := u.getRunsToUpdate(result, repositoryId, pipelineId)
	if err != nil {
		u.Logger.Error().
			Err(err).
			Stack().
			Dict("metadata", zerolog.Dict().Int("pipeline-id", pipelineId)).
			Send()
		return err
	}

	workItems, err := u.getAllWorkItems(builds)
	if err != nil {
		u.Logger.Error().
			Err(err).
			Stack().
			Dict("metatdata", zerolog.Dict().Int("pipeline-id", pipelineId)).
			Send()
		return err
	}

	_ = u.getAllWorkItemsToUpdatePrev(workItems, builds[0].Name)

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

func (u *AdoUsesCases) getAllWorkItems(builds []model.PipelineRuns) ([]model.WorkItem, error) {
	adoRep := u.Repository
	buildWorkItems, err := adoRep.GetBuildWorkItem(builds[1].Id, builds[0].Id)
	if err != nil {
		return []model.WorkItem{}, err
	}

	workItems := queryslice.TransformParallel(buildWorkItems, func(val model.BuildWorkItems, _ int) model.WorkItem {
		workItem, err := adoRep.GetWorkItem(val.Id)
		if err != nil {
			return model.WorkItem{}
		}
		return *workItem
	})

	return workItems, nil
}

func (u *AdoUsesCases) getAllWorkItemsToUpdatePrev(workItems []model.WorkItem, version string) []model.WorkItem {
	workItems = queryslice.Filter(workItems, func(pre model.WorkItem) bool {
		workItemVersion, _ := pre.Fields["Microsoft.VSTS.Build.IntegrationBuild"].(string)
		if workItemVersion == "" {
			return true
		}
		return strings.Compare(workItemVersion, version) == 1
	})

	return workItems
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
