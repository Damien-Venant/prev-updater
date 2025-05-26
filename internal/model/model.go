package model

type (
	PipelineRuns struct {
		Id     int    `json:"id"`
		Name   string `json:"name"`
		State  string `json:"state"`
		Result string `json:"result"`
	}

	BuildChanges struct {
		Id      string `json:"id"`
		Message string `json:"message"`
	}

	BuildWorkItems struct {
		Id  int    `json:"id"`
		Url string `json:"url"`
	}

	OperationFields struct {
		Op    string `json:"op"`
		Path  string `json:"path"`
		Value string `json:"value"`
	}
)
