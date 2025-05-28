package repository

import httpclient "github.com/prev-updater/pkg/http-client"

type AzureDevOpsRepository struct {
	httpclient *httpclient.HttpClient
}

func (r *AzureDevOpsRepository) GetPipelineRuns() error {
	return nil
}

func (r *AzureDevOpsRepository) GetPipelineRun() error {
	return nil
}

func (r *AzureDevOpsRepository) GetBuildWorkItem() error {
	return nil
}

func (r *AzureDevOpsRepository) GetWorkitem() error {
	return nil
}

func (r *AzureDevOpsRepository) UpdateWorkitemField() error {
	return nil
}
