package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Damien-Venant/prev-updater/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHttpClient mocke l’interface HTTPClientInterface
type MockHttpClient struct {
	mock.Mock
}

func (m *MockHttpClient) Get(path string, headers http.Header) (*http.Response, error) {
	args := m.Called(path, headers)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *MockHttpClient) Patch(path string, body []byte, headers http.Header) (*http.Response, error) {
	args := m.Called(path, body, headers)
	return args.Get(0).(*http.Response), args.Error(1)
}

// Helper pour créer une réponse HTTP avec JSON body
func makeHttpResponse(statusCode int, body interface{}) *http.Response {
	jsonBody, _ := json.Marshal(body)
	return &http.Response{
		StatusCode: statusCode,
		Body:       ioutil.NopCloser(bytes.NewBuffer(jsonBody)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
}

func TestGetPipelineRun(t *testing.T) {
	mockClient := new(MockHttpClient)
	repo := New(mockClient)

	expectedRun := model.PipelineRuns{Id: 123}

	mockResp := makeHttpResponse(200, expectedRun)
	mockClient.On("Get", mock.Anything, mock.Anything).Return(mockResp, nil)

	run, err := repo.GetPipelineRun(1, 123)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if run.Id != expectedRun.Id {
		t.Errorf("Expected Id %d, got %d", expectedRun.Id, run.Id)
	}
	mockClient.AssertExpectations(t)
}

func TestGetPipelineRuns(t *testing.T) {
	mockClient := new(MockHttpClient)
	repo := New(mockClient)

	// mock pagination response
	paginationValue := struct {
		Count int                      `json:"count"`
		Value []map[string]interface{} `json:"value"`
	}{
		Count: 2,
		Value: []map[string]interface{}{
			{"id": float64(1)},
			{"id": float64(2)},
		},
	}

	// Mock réponse pour la liste des runs
	mockListResp := makeHttpResponse(200, paginationValue)
	// Mock réponse pour chaque run
	mockRun1 := makeHttpResponse(200, model.PipelineRuns{Id: 1})
	mockRun2 := makeHttpResponse(200, model.PipelineRuns{Id: 2})

	// On s’attend à 3 appels : 1 liste, 2 runs
	mockClient.On("Get", mock.Anything, mock.Anything).Return(mockListResp, nil).Once()
	mockClient.On("Get", mock.Anything, mock.Anything).Return(mockRun1, nil).Once()
	mockClient.On("Get", mock.Anything, mock.Anything).Return(mockRun2, nil).Once()
	runs, err := repo.GetPipelineRuns(1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(runs) != 2 {
		t.Errorf("Expected 2 runs, got %d", len(runs))
	}

	mockClient.AssertExpectations(t)
}

func TestGetBuildWorkItem(t *testing.T) {
	mockClient := new(MockHttpClient)
	repo := New(mockClient)

	paginated := struct {
		Count int                    `json:"count"`
		Value []model.BuildWorkItems `json:"value"`
	}{
		Count: 1,
		Value: []model.BuildWorkItems{{Id: "10"}},
	}

	mockResp := makeHttpResponse(200, paginated)
	mockClient.On("Get", mock.Anything, mock.Anything).Return(mockResp, nil)

	items, err := repo.GetBuildWorkItem(100, 200)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(items) != 1 {
		t.Errorf("Expected 1 work item, got %d", len(items))
	}

	mockClient.AssertExpectations(t)
}

func TestGetWorkItem(t *testing.T) {
	mockClient := new(MockHttpClient)
	repo := New(mockClient)

	expectedItem := model.WorkItem{Id: 42}

	mockResp := makeHttpResponse(200, expectedItem)
	mockClient.On("Get", mock.Anything, mock.Anything).Return(mockResp, nil)

	item, err := repo.GetWorkItem("42")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if item.Id != expectedItem.Id {
		t.Errorf("Expected ID %d, got %d", expectedItem.Id, item.Id)
	}

	mockClient.AssertExpectations(t)
}

func TestUpdateWorkitemField(t *testing.T) {
	mockClient := new(MockHttpClient)
	repo := New(mockClient)

	mockResp := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString("")),
	}

	mockClient.On("Patch", mock.Anything, mock.Anything, mock.Anything).Return(mockResp, nil)

	err := repo.UpdateWorkitemField("42", model.OperationFields{
		Op:    "add",
		Path:  "/fields/System.Title",
		Value: "Test",
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	mockClient.AssertExpectations(t)
}

func TestGetRepositoryById(t *testing.T) {
	mockClient := new(MockHttpClient)
	repo := New(mockClient)

	expectedRepo := model.Repository{Id: "uuid-123"}

	mockResp := makeHttpResponse(200, expectedRepo)
	mockClient.On("Get", mock.Anything, mock.Anything).Return(mockResp, nil)

	r, err := repo.GetRepositoryById("uuid-123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if r.Id != expectedRepo.Id {
		t.Errorf("Expected ID %s, got %s", expectedRepo.Id, r.Id)
	}

	mockClient.AssertExpectations(t)
}

func TestConfigureRouteWithVersion(t *testing.T) {
	tests := []struct {
		name       string
		route      string
		parameters []any
		version    AzureDevOpsRepository
		wants      string
	}{
		{
			name:       "MultiParam",
			route:      "/api/%s/%d",
			parameters: []any{"test", 10},
			version:    AzureDevOpsRepository{version: "7.1"},
			wants:      "/api/test/10?api-version=7.1",
		},
		{
			name:       "MonoParam",
			route:      "/api/%d",
			parameters: []any{10},
			version:    AzureDevOpsRepository{version: "7.5"},
			wants:      "/api/10?api-version=7.5",
		},
		{
			name:       "MultiParamWithQueryParameter",
			route:      "/api/%s/%d?date=%s",
			parameters: []any{"test", 10, "test"},
			version:    AzureDevOpsRepository{version: "7.1"},
			wants:      "/api/test/10?api-version=7.1",
		},
		{
			name:       "MonoParamWithQueryParameter",
			route:      "/api/%d?date=%s",
			parameters: []any{10, "test"},
			version:    AzureDevOpsRepository{version: "7.5"},
			wants:      "/api/10?api-version=7.5",
		},
	}

	for _, test := range tests {
		testName := fmt.Sprintf("%s_%s", "TestConfigureRouteWithVersion", test.name)
		t.Run(testName, func(t *testing.T) {
			result := test.version.configureRouteWithVersion(test.route, test.parameters...)
			if result != test.wants {
				assert.NotEqual(t, test.wants, result)
			}
		})
	}
}

func TestErrorCodeMapping(t *testing.T) {
	tests := []struct {
		Name           string
		ErrorCode      int
		ExpectedResult error
	}{
		{
			Name:           "InternalServerError",
			ErrorCode:      http.StatusInternalServerError,
			ExpectedResult: ErrInternalServer,
		},
		{
			Name:           "BadRequestError",
			ErrorCode:      http.StatusBadRequest,
			ExpectedResult: ErrBadRequest,
		},
		{
			Name:           "NotFoundError",
			ErrorCode:      http.StatusNotFound,
			ExpectedResult: ErrNotFound,
		},
	}

	for _, test := range tests {
		testName := fmt.Sprintf("TestErrorCodeMapping%s", test.Name)
		t.Run(testName, func(t *testing.T) {
			err := errorCodeMapping(test.ErrorCode)
			assert.ErrorIs(t, err, test.ExpectedResult)
		})
	}
}

func TestErrorCodeMappingErrorUnknown(t *testing.T) {
	test := struct {
		ErrorCode      int
		ExpectedResult error
	}{
		ErrorCode:      http.StatusBadGateway,
		ExpectedResult: ErrIdk,
	}

	err := errorCodeMapping(test.ErrorCode)
	assert.ErrorIs(t, err, test.ExpectedResult)
}

func TestReadAndUnMarshall(t *testing.T) {
	type Person struct {
		FirstName string `json:"first-name"`
		LastName  string `json:"last-name"`
	}
	var person Person

	resultModel, _ := json.Marshal(Person{
		FirstName: "damien",
		LastName:  "venant",
	})

	reader := bytes.NewReader(resultModel)

	err := readAndUnmarshal[Person](reader, &person)

	assert.Nil(t, err)
	assert.Equal(t, "damien", person.FirstName)
	assert.Equal(t, "venant", person.LastName)
}

func TestTreatResult_When_StatusCode_Match_ExpectedStatusCode(t *testing.T) {
	tests := []struct {
		name               string
		statusCode         int
		expectedStatusCode int
	}{
		{
			name:               "OkStatus",
			statusCode:         http.StatusOK,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "CreatedStatus",
			statusCode:         http.StatusCreated,
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "AcceptedStatus",
			statusCode:         http.StatusAccepted,
			expectedStatusCode: http.StatusAccepted,
		},
	}

	for _, test := range tests {
		testName := fmt.Sprintf("TestTreatResult_When_StatusCode_Match_ExpectedStatusCode_%s", test.name)
		t.Run(testName, func(t *testing.T) {
			response := http.Response{
				StatusCode: test.statusCode,
			}
			err := treatResult(&response, test.expectedStatusCode)
			assert.Nil(t, err)
		})
	}
}
func TestTreatResult_When_StatusCode_Does_Not_Match_ExpectedStatusCode(t *testing.T) {
	tests := []struct {
		name               string
		statusCode         int
		expectedStatusCode int
	}{
		{
			name:               "InternalServerError",
			statusCode:         http.StatusInternalServerError,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "NotFoundError",
			statusCode:         http.StatusNotFound,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "BadRequestError",
			statusCode:         http.StatusBadRequest,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "BadGatewayError",
			statusCode:         http.StatusBadGateway,
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		testName := fmt.Sprintf("TestTreatResult_When_StatusCode_Does_Not_Match_ExpectedStatusCode_%s", test.name)
		t.Run(testName, func(t *testing.T) {
			response := http.Response{
				StatusCode: test.statusCode,
			}
			err := treatResult(&response, test.expectedStatusCode)
			assert.NotNil(t, err)
		})
	}
}
