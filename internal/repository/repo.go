package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/prev-updater/internal/model"
	httpclient "github.com/prev-updater/pkg/http-client"
)

const (
	apiVersion string = "7.1"
)

var (
	NotFoundError       error = errors.New("")
	InternalServerError error = errors.New("")
)

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

func (r *AzureDevOpsRepository) GetPipelineRuns(pipelineId int) (error, []model.PipelineRuns) {
	var paginationValue model.PaginatedValue[model.PipelineRuns]
	url := r.configureRouteWithVersion("pipelines/%d/runs", pipelineId)
	httpResponse, err := r.client.Get(url, nil)
	if err != nil {
		return err, []model.PipelineRuns{}
	}

	if err := readAndUnmarshal[model.PaginatedValue[model.PipelineRuns]](httpResponse.Body, &paginationValue); err != nil {
		return err, []model.PipelineRuns{}
	}

	return nil, paginationValue.Value
}

func (r *AzureDevOpsRepository) GetPipelineRun(pipelineId, runId int) (error, *model.PipelineRuns) {
	var result model.PipelineRuns
	url := r.configureRouteWithVersion("pipelines/%d/runs/%d", pipelineId, runId)
	httpResponse, err := r.client.Get(url, nil)
	if err != nil {
		return err, nil
	}

	if err := readAndUnmarshal[model.PipelineRuns](httpResponse.Body, &result); err != nil {
		return err, nil
	}
	return nil, &result
}

func (r *AzureDevOpsRepository) GetBuildWorkItem(buildId int) (error, *model.BuildWorkItems) {
	var workItem model.BuildWorkItems
	url := r.configureRouteWithVersion("build/builds/%d/workitems", buildId)
	httpResponse, err := r.client.Get(url, nil)
	if err != nil {
		return err, nil
	}

	if err := readAndUnmarshal[model.BuildWorkItems](httpResponse.Body, &workItem); err != nil {
		return err, nil
	}
	return nil, &workItem
}

func (r *AzureDevOpsRepository) GetWorkitem(workItemId int) (error, *model.BuildWorkItems) {
	var buildWorkItems model.BuildWorkItems
	url := r.configureRouteWithVersion("wit/workitems/%d", workItemId)
	httpResponse, err := r.client.Get(url, nil)
	if err != nil {
		return err, nil
	}

	if err := readAndUnmarshal[model.BuildWorkItems](httpResponse.Body, &buildWorkItems); err != nil {
		return err, nil
	}
	return nil, &buildWorkItems
}

func (r *AzureDevOpsRepository) UpdateWorkitemField(workItemId int, version string) error {
	url := r.configureRouteWithVersion("build/builds/workitems/%d", workItemId)
	httpResponse, err := r.client.Get(url, nil)
	if err != nil {
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

// TODO: write a test for this function
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
	return nil
}
