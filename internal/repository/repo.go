package repository

import (
	"encoding/json"
	"fmt"
	"io"

	httpclient "github.com/prev-updater/pkg/http-client"
)

type AzureDevOpsRepository struct {
	client  *httpclient.HttpClient
	version string
}

type jsonMap map[string]interface{}

func New(client *httpclient.HttpClient) *AzureDevOpsRepository {
	return &AzureDevOpsRepository{
		client:  client,
		version: "7.1",
	}
}

func (r *AzureDevOpsRepository) GetPipelineRuns(pipelineId int) error {
	url := r.configureRootWithVersion("/pipelines/%d/runs", pipelineId)
	httpResponse, err := r.client.Get(url, nil)
	if err != nil {
		return err
	}

	body, _ := io.ReadAll(httpResponse.Body)
	fmt.Println(string(body))
	return nil
}

func (r *AzureDevOpsRepository) GetPipelineRun(pipelineId, runId int) error {
	var result jsonMap
	url := r.configureRootWithVersion("/pipelines/%d/runs/%d", pipelineId, runId)
	httpResponse, err := r.client.Get(url, nil)
	if err != nil {
		return err
	}

	body, _ := io.ReadAll(httpResponse.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}

func (r *AzureDevOpsRepository) GetBuildWorkItem(buildId int) error {
	url := r.configureRootWithVersion("/build/builds/%d/workitems", buildId)
	httpResponse, err := r.client.Get(url, nil)
	if err != nil {
		return err
	}

	body, _ := io.ReadAll(httpResponse.Body)
	fmt.Println(body)
	return nil
}

func (r *AzureDevOpsRepository) GetWorkitem(workItemId int) error {
	url := r.configureRootWithVersion("/wit/workitems/%d", workItemId)
	httpResponse, err := r.client.Get(url, nil)
	if err != nil {
		return err
	}

	body, _ := io.ReadAll(httpResponse.Body)
	fmt.Println(body)
	return nil
}

func (r *AzureDevOpsRepository) UpdateWorkitemField(workItemId int, version string) error {
	url := r.configureRootWithVersion("/build/builds/workitems/%d", workItemId)
	httpResponse, err := r.client.Get(url, nil)
	if err != nil {
		return err
	}

	body, _ := io.ReadAll(httpResponse.Body)
	fmt.Println(body)
	return nil
}

func (r *AzureDevOpsRepository) configureRootWithVersion(route string, values ...any) string {
	url := fmt.Sprintf(route, values...)
	url = fmt.Sprintf("%s?api-version=%s", url, r.version)
	return url
}
