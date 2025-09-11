package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/prev-updater/internal/model"
	httpclient "github.com/prev-updater/pkg/http-client"
)

const (
	apiVersion string = "7.1"
)

var (
	NotFoundError       error = errors.New("Ressource not found")
	InternalServerError error = errors.New("Internal server error")
	BadRequestError     error = errors.New("Bad request error")
	IdkError            error = errors.New("IDK what's happened")
)

var mappingError map[int]error = map[int]error{
	http.StatusBadRequest:          BadRequestError,
	http.StatusInternalServerError: InternalServerError,
	http.StatusNotFound:            NotFoundError,
}

type AzureDevOpsRepository struct {
	client  *httpclient.HttpClient
	version string
}

func New(client *httpclient.HttpClient) *AzureDevOpsRepository {
	return &AzureDevOpsRepository{
		client:  client,
		version: apiVersion,
	}
}

func (r *AzureDevOpsRepository) GetPipelineRuns(pipelineId int) ([]model.PipelineRuns, error) {
	type PipelineRuns model.PaginatedValue[model.PipelineRuns]
	var paginationValue PipelineRuns
	url := r.configureRouteWithVersion("pipelines/%d/runs", pipelineId)
	httpResponse, err := r.client.Get(url, nil)

	if err != nil {
		return []model.PipelineRuns{}, err
	}

	if err := treatResult(httpResponse, http.StatusOK); err != nil {
		return []model.PipelineRuns{}, err
	}

	if err := readAndUnmarshal[PipelineRuns](httpResponse.Body, &paginationValue); err != nil {
		return []model.PipelineRuns{}, err
	}

	return paginationValue.Value, err
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

	if err := readAndUnmarshal[model.PipelineRuns](httpResponse.Body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *AzureDevOpsRepository) GetBuildWorkItem(buildId int) ([]model.BuildWorkItems, error) {
	type BuildWorkItems model.PaginatedValue[model.BuildWorkItems]
	var workItem BuildWorkItems
	url := r.configureRouteWithVersion("build/builds/%d/workitems", buildId)
	httpResponse, err := r.client.Get(url, nil)
	if err != nil {
		return []model.BuildWorkItems{}, err
	}

	if err := treatResult(httpResponse, http.StatusOK); err != nil {
		return []model.BuildWorkItems{}, err
	}

	if err := readAndUnmarshal[BuildWorkItems](httpResponse.Body, &workItem); err != nil {
		return []model.BuildWorkItems{}, err
	}
	return workItem.Value, nil
}

func (r *AzureDevOpsRepository) GetWorkitem(workItemId int) (*model.BuildWorkItems, error) {
	var buildWorkItems model.BuildWorkItems
	url := r.configureRouteWithVersion("wit/workItems/%d", workItemId)
	httpResponse, err := r.client.Get(url, nil)
	if err != nil {
		return nil, err
	}

	if err := treatResult(httpResponse, http.StatusOK); err != nil {
		return nil, err
	}

	if err := readAndUnmarshal[model.BuildWorkItems](httpResponse.Body, &buildWorkItems); err != nil {
		return nil, err
	}
	return &buildWorkItems, nil
}

func (r *AzureDevOpsRepository) UpdateWorkitemField(workItemId int, version string) error {
	url := r.configureRouteWithVersion("build/builds/workitems/%d", workItemId)
	httpResponse, err := r.client.Get(url, nil)
	if err != nil {
		return err
	}

	if err := treatResult(httpResponse, http.StatusOK); err != nil {
		return err
	}

	_, _ = io.ReadAll(httpResponse.Body)
	return nil
}

func (r *AzureDevOpsRepository) configureRouteWithVersion(route string, values ...any) string {
	url := fmt.Sprintf(route, values...)
	url = fmt.Sprintf("_apis/%s?api-version=%s", url, r.version)
	return url
}

func readAndUnmarshal[T any](body io.Reader, model *T) error {
	buffBody, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(buffBody, model); err != nil {
		return err
	}
	return nil
}

func errorCodeMapping(errorCode int) error {
	if err, ok := mappingError[errorCode]; !ok {
		return IdkError
	} else {
		return err
	}
}

func treatResult(response *http.Response, expectedReturnCode int) error {
	if returnCode := response.StatusCode; expectedReturnCode != returnCode {
		return errorCodeMapping(returnCode)
	}
	return nil
}
