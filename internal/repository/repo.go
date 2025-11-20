package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/Damien-Venant/prev-updater/internal/model"
	httpClient "github.com/Damien-Venant/prev-updater/pkg/http-client"
)

const (
	apiVersion string = "7.1"
)

type AzureDevOpsRepository struct {
	client  httpClient.HttpClientInterface
	version string
}

func NewAdoRepository(client httpClient.HttpClientInterface) *AzureDevOpsRepository {
	return &AzureDevOpsRepository{
		client:  client,
		version: apiVersion,
	}
}

func (r *AzureDevOpsRepository) GetPipelineRuns(pipelineId int) ([]model.PipelineRuns, error) {
	type PipelineRuns model.PaginatedValue[map[string]interface{}]
	url := r.configureRouteWithVersion("pipelines/%d/runs", pipelineId)
	httpResponse, err := r.client.Get(url, nil)
	if err != nil {
		return []model.PipelineRuns{}, err
	}

	if err := treatResult(httpResponse, http.StatusOK); err != nil {
		return []model.PipelineRuns{}, err
	}

	var paginationValue PipelineRuns
	if err := readAndUnmarshal(httpResponse.Body, &paginationValue); err != nil {
		return []model.PipelineRuns{}, err
	}

	pipelineRuns := make([]model.PipelineRuns, paginationValue.Count)
	outChan := make(chan model.PipelineRuns)
	errChan := make(chan error)

	for i := 0; i < paginationValue.Count; i++ {
		runId, ok := paginationValue.Value[i]["id"].(float64)
		if !ok {
			return []model.PipelineRuns{}, err
		}
		go r.batchGetPipelineRunRequest(pipelineId, int(runId), outChan, errChan)
	}

	for i := 0; i < paginationValue.Count; i++ {
		select {
		case res := <-outChan:
			pipelineRuns[i] = res
		case err := <-errChan:
			return []model.PipelineRuns{}, err
		case <-time.After(time.Second * 5):
			return []model.PipelineRuns{}, errors.New("Timeout")
		}
	}

	slices.SortStableFunc(pipelineRuns, func(i, j model.PipelineRuns) int {
		return j.Id - i.Id
	})
	return pipelineRuns, err
}

func (r *AzureDevOpsRepository) batchGetPipelineRunRequest(pipelineId int, runId int, outChan chan model.PipelineRuns, errChan chan error) {
	if res, err := r.GetPipelineRun(pipelineId, int(runId)); err != nil {
		errChan <- err
	} else {
		outChan <- *res
	}
}

func (r *AzureDevOpsRepository) GetPipelineRun(pipelineId, runId int) (*model.PipelineRuns, error) {
	var result model.PipelineRuns
	url := r.configureRouteWithVersion("pipelines/%d/runs/%d", pipelineId, runId)
	httpResponse, err := r.client.Get(url, nil)
	if err != nil {
		return nil, err
	}

	if err := treatResult(httpResponse, http.StatusOK); err != nil {
		return nil, err
	}

	if err := readAndUnmarshal(httpResponse.Body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *AzureDevOpsRepository) GetBuildWorkItem(fromBuildId, toBuildId int) ([]model.BuildWorkItems, error) {
	type BuildWorkItems model.PaginatedValue[model.BuildWorkItems]
	var workItem BuildWorkItems
	url := r.configureRouteWithVersion("build/workitems?fromBuildId=%d&toBuildId=%d", fromBuildId, toBuildId)
	httpResponse, err := r.client.Get(url, nil)
	if err != nil {
		return []model.BuildWorkItems{}, err
	}

	if err := treatResult(httpResponse, http.StatusOK); err != nil {
		return []model.BuildWorkItems{}, err
	}

	if err := readAndUnmarshal(httpResponse.Body, &workItem); err != nil {
		return []model.BuildWorkItems{}, err
	}
	return workItem.Value, nil
}

func (r *AzureDevOpsRepository) GetWorkItem(workItemId string) (*model.WorkItem, error) {
	var buildWorkItems model.WorkItem
	url := r.configureRouteWithVersion("wit/workItems/%s", workItemId)
	httpResponse, err := r.client.Get(url, nil)
	if err != nil {
		return nil, err
	}

	if err := treatResult(httpResponse, http.StatusOK); err != nil {
		return nil, err
	}

	if err := readAndUnmarshal(httpResponse.Body, &buildWorkItems); err != nil {
		return nil, err
	}
	return &buildWorkItems, nil
}

func (r *AzureDevOpsRepository) UpdateWorkitemField(workItemId string, operation model.OperationFields) error {
	url := r.configureRouteWithVersion("wit/workItems/%s", workItemId)
	model, err := json.Marshal([]model.OperationFields{operation})
	if err != nil {
		return err
	}
	httpResponse, err := r.client.Patch(url, model, nil)
	if err != nil {
		return err
	}
	if err := treatResult(httpResponse, http.StatusOK); err != nil {
		return err
	}

	return nil
}

func (r *AzureDevOpsRepository) GetRepositoryById(uuid string) (*model.Repository, error) {
	var result model.Repository
	url := r.configureRouteWithVersion("git/repositories/%s", uuid)
	httpResponse, err := r.client.Get(url, nil)
	if err != nil {
		return nil, err
	}
	if err := treatResult(httpResponse, http.StatusOK); err != nil {
		return nil, err
	}
	if err := readAndUnmarshal(httpResponse.Body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *AzureDevOpsRepository) configureRouteWithVersion(route string, values ...any) string {
	url := fmt.Sprintf(route, values...)
	if strings.Contains(url, "?") {
		url = fmt.Sprintf("_apis/%s&api-version=%s", url, r.version)

	} else {
		url = fmt.Sprintf("_apis/%s?api-version=%s", url, r.version)
	}
	return url
}
