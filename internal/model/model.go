package model

type (
	PipelineRuns struct {
		Id     int
		Name   string
		State  string
		Result string
	}

	BuildChanges struct {
		Id      string
		Message string
	}

	BuildWorkItems struct {
		Id  int
		Url string
	}

	OperationFields struct {
		Op    string
		Path  string
		Value string
	}
)
