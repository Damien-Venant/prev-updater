package repository

import httpclient "github.com/prev-updater/pkg/http-client"

type AzureDevOpsRepository struct {
	httpclient *httpclient.HttpClient
}

func (r *AzureDevOpsRepository) GetPipelineRuns() {

}

func (r *AzureDevOpsRepository) GetPipelineRun() {

}
