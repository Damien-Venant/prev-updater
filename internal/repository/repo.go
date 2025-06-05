package repository

import (
	"fmt"
	"io"

	"github.com/prev-updater/pkg/http-client"
)

type AzureDevOpsRepository struct {
	client *httpclient.HttpClient
}

func (r *AzureDevOpsRepository) GetPipelineRuns(pipelineId int) error {
	url := fmt.Sprintf("/pipelines/%d/runs", pipelineId)
	httpRequest, err := r.client.Get(url, nil)
	if err != nil {
		return err
	}

	body, _ := io.ReadAll(httpRequest.Body)
	fmt.Println(body)
	return nil
}

func (r *AzureDevOpsRepository) GetPipelineRun(pipelineId, runId int) error {
	url := fmt.Sprintf("/pipelines/%d/runs/%d", pipelineId, runId)
	httpRequest, err := r.client.Get(url, nil)
	if err != nil {
		return err
	}

	body, _ := io.ReadAll(httpRequest.Body)
	fmt.Println(body)

	return nil
}

func (r *AzureDevOpsRepository) GetBuildWorkItem(buildId int) error {
	url := fmt.Sprintf("/build/builds/%d/workitems", buildId)
	httpRequest, err := r.client.Get(url, nil)
	if err != nil {
		return err
	}

	body, _ := io.ReadAll(httpRequest.Body)
	fmt.Println(body)
	return nil
}

func (r *AzureDevOpsRepository) GetWorkitem(workItemId int) error {
	url := fmt.Sprintf("/wit/workitems/%d", workItemId)
	httpRequest, err := r.client.Get(url, nil)
	if err != nil {
		return err
	}

	body, _ := io.ReadAll(httpRequest.Body)
	fmt.Println(body)
	return nil
}

func (r *AzureDevOpsRepository) UpdateWorkitemField(version string) error {
	url := fmt.Sprintf("/build/builds/workitems")
	httpRequest, err := r.client.Get(url, nil)
	if err != nil {
		return err
	}

	body, _ := io.ReadAll(httpRequest.Body)
	fmt.Println(body)
	return nil
}
