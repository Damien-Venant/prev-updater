package usescases

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Damien-Venant/prev-updater/internal/model"
	"github.com/Damien-Venant/prev-updater/pkg/queryslice"
	"github.com/rs/zerolog"
)

const (
	AdoIntegrationBuildFieldName string = "Microsoft.VSTS.Build.IntegrationBuild"
	AdoIntegrationPath           string = "/fields/" + AdoIntegrationBuildFieldName
	AdoTitleFieldName            string = "System.Title"
	AdoTagsFieldName             string = "System.Tags"
)

type AdoRepository interface {
	GetPipelineRuns(pipelineId int) ([]model.PipelineRuns, error)
	GetPipelineRun(pipelineId, runId int) (*model.PipelineRuns, error)
	GetBuildWorkItem(fromBuildId, toBuildId int) ([]model.BuildWorkItems, error)
	GetWorkItem(workItemId string) (*model.WorkItem, error)
	GetRepositoryById(uuid string) (*model.Repository, error)
	UpdateWorkitemField(workItemId string, operation model.OperationFields) error
}

type N8nRepository interface {
	PostWebhook(data model.N8nResult) error
}

type (
	Version      [4]int
	AdoUsesCases struct {
		N8nRepo    N8nRepository
		Repository AdoRepository
		Logger     *zerolog.Logger
	}

	UpdateFieldsParams struct {
		PipelineId   int
		RepositoryId string
		FieldName    string
		BranchName   string
	}
)

func NewAdoUsesCases(adoRepository AdoRepository, N8nRepo N8nRepository, logger *zerolog.Logger) *AdoUsesCases {
	return &AdoUsesCases{
		N8nRepo:    N8nRepo,
		Repository: adoRepository,
		Logger:     logger,
	}
}

func (u *AdoUsesCases) UpdateFieldsByPipelineId(pipelineId int) error {
	return nil
}

func (u *AdoUsesCases) UpdateFieldsByLastRuns(param UpdateFieldsParams) error {
	adoRep := u.Repository
	result, err := adoRep.GetPipelineRuns(param.PipelineId)
	if err != nil {
		return err
	} else if len(result) == 0 {
		return nil
	}

	builds, err := u.getRunsToUpdate(result, param.RepositoryId, param.PipelineId, param.BranchName)
	if err != nil {
		return err
	}
	lastBuild := builds[0]

	workItems, err := u.getAllWorkItems(builds)
	if err != nil {
		return err
	}

	versionName := lastBuild.Name
	tabFieldName := strings.Split(param.FieldName, "/")
	fieldName := tabFieldName[len(tabFieldName)-1]
	workItemsToUpdatePrev := u.getAllWorkItemsToUpdatePrev(workItems, builds[0].Name, fieldName)

	if len(workItemsToUpdatePrev) > 0 {
		var errMap error = nil
		for _, workItem := range workItemsToUpdatePrev {
			err = u.updateFields(strconv.FormatInt(int64(workItem.Id), 10), versionName, param.FieldName)
			if err != nil {
				errMap = errors.Join(errMap, err)
			}
		}
		if errMap != nil {
			return errMap
		}
	}

	if len(workItems) > 0 {
		u.updateAdoIntegrationBuild(workItems, versionName)
		err := u.sendDataToN8N(workItems, versionName, param.BranchName)
		return err
	}
	return nil
}

// getRunsToUpdate is used to return last build and N-1 last build
// It's return a array where the first index is last build and second index is N-1 last build
// If the build have not previous build the N-1 last build is last build on defaultBranch
func (u *AdoUsesCases) getRunsToUpdate(builds []model.PipelineRuns, repositoryId string, pipelineId int, branchName string) ([]model.PipelineRuns, error) {
	var lastBuild model.PipelineRuns
	adoRep := u.Repository
	defaultRefName, err := adoRep.GetRepositoryById(repositoryId)
	if err != nil {
		return nil, err
	}
	builds = queryslice.Filter(builds, func(pre model.PipelineRuns) bool {
		return pre.State == "completed"
	})

	index := 0
	if branchName == "" {
		lastBuild = builds[0]
	} else {
		index = queryslice.FindIndex(builds, func(pre model.PipelineRuns) bool {
			return strings.Contains(pre.Resources.Repositories.Self.RefName, branchName)
		})
		if index < 0 {
			return []model.PipelineRuns{}, ErrBranchNameNotExist
		}
		lastBuild = builds[index]
	}

	buildsOnSameRef := queryslice.Filter(builds, func(pre model.PipelineRuns) bool {
		return pre.Resources.Repositories.Self.RefName == lastBuild.Resources.Repositories.Self.RefName
	})

	if len(buildsOnSameRef) > 1 {
		return []model.PipelineRuns{buildsOnSameRef[0], buildsOnSameRef[1]}, nil
	}

	lastBuildOnDefaultRefName := queryslice.Filter(builds[index:], func(pre model.PipelineRuns) bool {
		return pre.Resources.Repositories.Self.RefName == defaultRefName.DefaultBranch
	})[0]

	return []model.PipelineRuns{builds[index], lastBuildOnDefaultRefName}, nil
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

func (u *AdoUsesCases) getAllWorkItemsToUpdatePrev(workItems []model.WorkItem, version, fieldName string) []model.WorkItem {
	workItems = queryslice.Filter(workItems, func(pre model.WorkItem) bool {
		workItemVers, _ := pre.Fields[fieldName].(string)
		if workItemVers == "" {
			return true
		}
		actualVersion := newVersion(version)
		workItemVersion := newVersion(workItemVers)
		return actualVersion.isSmallerThan(workItemVersion) == 1
	})

	return workItems
}

func (u *AdoUsesCases) updateAdoIntegrationBuild(workItems []model.WorkItem, version string) error {
	var errMap error = nil
	for _, workItem := range workItems {
		var concatenateVersion string
		val, _ := workItem.Fields[AdoIntegrationBuildFieldName].(string)
		if val != "" && !strings.Contains(val, version) {
			concatenateVersion = fmt.Sprintf("%s | %s", val, version)
		} else {
			concatenateVersion = version
		}
		if err := u.updateFields(strconv.FormatInt(int64(workItem.Id), 10), concatenateVersion, AdoIntegrationPath); err != nil {
			errMap = errors.Join(errMap)
		}
	}
	return errMap
}

func (u *AdoUsesCases) updateFields(woritemId, name, path string) error {
	repo := u.Repository
	modelToUpdload := model.OperationFields{
		Op:    "add",
		Path:  path,
		Value: name,
	}

	return repo.UpdateWorkitemField(woritemId, modelToUpdload)
}

func (u *AdoUsesCases) sendDataToN8N(workitems []model.WorkItem, version string, sourceBranch string) error {
	data := WorkItemToN8NResult(workitems)
	data.Version = version
	data.SourceBranch = sourceBranch

	return u.N8nRepo.PostWebhook(data)
}

func (actual Version) isSmallerThan(targetVersion Version) int {
	for index := 0; index < len(targetVersion); index++ {
		if actual[index] > targetVersion[index] {
			return -1
		} else if actual[index] < targetVersion[index] {
			return 1
		}
	}
	return 0
}

func (actual Version) isHigherThan(targetVersion Version) int {
	return -actual.isSmallerThan(targetVersion)
}

func newVersion(version string) Version {
	result := Version{}
	res := strings.Split(version, ".")
	for index, val := range res {
		resultConv, _ := strconv.ParseInt(val, 10, 32)
		result[index] = int(resultConv)
	}
	return result
}

func WorkItemToN8NResult(workitems []model.WorkItem) model.N8nResult {
	wItems := make([]model.N8NWorkItems, len(workitems))

	for index, val := range workitems {
		title := val.Fields[AdoTitleFieldName].(string)
		tags := val.Fields[AdoTagsFieldName].(string)
		integrationBuild := val.Fields[AdoIntegrationBuildFieldName].(string)
		integrationBuild = strings.TrimSpace(integrationBuild)

		wItems[index] = model.N8NWorkItems{
			Id:               val.Id,
			Title:            title,
			Tags:             strings.Split(tags, ";"),
			IntegrationBuild: strings.Split(integrationBuild, "|"),
		}
	}

	return model.N8nResult{
		WorkItems: wItems,
	}
}
