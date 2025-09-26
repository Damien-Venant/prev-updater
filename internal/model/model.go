package model

type (
	PaginatedValue[T any] struct {
		Count int `json:"count"`
		Value []T `json:"value"`
	}

	PipelineRuns struct {
		Id        int    `json:"id"`
		Name      string `json:"name"`
		State     string `json:"state"`
		Result    string `json:"result"`
		Resources *struct {
			Repositories *struct {
				Self struct {
					RefName string `json:"refName"`
					Version string `json:"Version"`
				} `json:"self"`
			} `json:"repositories"`
		} `json:"resources"`
	}

	BuildChanges struct {
		Id      string `json:"id"`
		Message string `json:"message"`
	}

	BuildWorkItems struct {
		Id  string `json:"id"`
		Url string `json:"url"`
	}

	OperationFields struct {
		Op    string `json:"op"`
		Path  string `json:"path"`
		Value string `json:"value"`
	}

	Repository struct {
		Id            string `json:"id"`
		Name          string `json:"name"`
		DefaultBranch string `json:"defaultBranch"`
		Url           string `json:"url"`
		RemoteUrl     string `json:"remoteUrl"`
	}
)
