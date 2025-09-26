package usescases

import (
	"errors"
	"testing"

	"github.com/Damien-Venant/prev-updater/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock du repository
type MockRepository struct {
	mock.Mock
}

type MockAdoUseCases struct {
	mock.Mock
}

func (m *MockRepository) GetRepositoryById(repositoryId string) (*model.Repository, error) {
	args := m.Called(repositoryId)
	val := args.Get(0).(model.Repository)
	return &val, args.Error(1)
}
func (m *MockRepository) GetPipelineRuns(pipelineId int) ([]model.PipelineRuns, error) {
	return nil, nil
}
func (m *MockRepository) GetPipelineRun(pipelineId, runId int) (*model.PipelineRuns, error) {
	return nil, nil
}
func (m *MockRepository) GetBuildWorkItem(fromBuildId, toBuildId int) ([]model.BuildWorkItems, error) {
	args := m.Called(fromBuildId, toBuildId)
	val := args.Get(0).([]model.BuildWorkItems)
	return val, args.Error(1)
}
func (m *MockRepository) GetWorkItem(workItemId string) (*model.WorkItem, error) {
	args := m.Called(workItemId)
	val := args.Get(0).(model.WorkItem)
	return &val, args.Error(1)
}
func (m *MockRepository) UpdateWorkitemField(workItemId string, operation model.OperationFields) error {
	return nil
}

// Setup d’un run
func createPipelineRun(ref string, name string, id int) model.PipelineRuns {
	return model.PipelineRuns{
		Id:    id,
		State: "completed",
		Name:  name,
		Resources: &struct {
			Repositories *struct {
				Self struct {
					RefName string "json:\"refName\""
					Version string "json:\"Version\""
				} "json:\"self\""
			} "json:\"repositories\""
		}{
			Repositories: &struct {
				Self struct {
					RefName string "json:\"refName\""
					Version string "json:\"Version\""
				} "json:\"self\""
			}{
				Self: struct {
					RefName string "json:\"refName\""
					Version string "json:\"Version\""
				}{RefName: ref},
			},
		},
	}
}

func createWorkItem(id int, fields map[string]interface{}) model.WorkItem {
	return model.WorkItem{
		Id:     id,
		Fields: fields,
	}
}

func (m *MockAdoUseCases) updateFields(id string, value string, path string) error {
	args := m.Called(id, value, path)
	return args.Error(0)
}

func TestGetRunsToUpdate_TwoBuildsOnSameRef(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := &AdoUsesCases{Repository: mockRepo}

	builds := []model.PipelineRuns{
		createPipelineRun("refs/heads/feature-1", "", 1),
		createPipelineRun("refs/heads/main", "", 2),
		createPipelineRun("refs/heads/feature-1", "", 3),
		createPipelineRun("refs/heads/main", "", 4),
		createPipelineRun("refs/heads/feature-1", "", 5),
		createPipelineRun("refs/heads/feature-1", "", 6),
	}

	mockRepo.On("GetRepositoryById", "repo-id").Return(model.Repository{DefaultBranch: "refs/heads/main"}, nil)

	result, err := uc.getRunsToUpdate(builds, "repo-id", 123)

	assert.NoError(t, err)
	assert.Equal(t, builds[0], result[0])
	assert.Equal(t, builds[2], result[1])
}

func TestGetRunsToUpdate_OnlyOneBuildOnRef(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := &AdoUsesCases{Repository: mockRepo}

	builds := []model.PipelineRuns{
		createPipelineRun("refs/heads/feature-1", "", 1), // only one on feature-1
		createPipelineRun("refs/heads/main", "", 2),      // default branch
		createPipelineRun("refs/heads/main", "", 3),
		createPipelineRun("refs/heads/main", "", 4), // default branch
		createPipelineRun("refs/heads/main", "", 5),
	}

	mockRepo.On("GetRepositoryById", "repo-id").Return(model.Repository{DefaultBranch: "refs/heads/main"}, nil)

	result, err := uc.getRunsToUpdate(builds, "repo-id", 123)

	assert.NoError(t, err)
	assert.Equal(t, builds[0], result[0]) // Last on current ref
	assert.Equal(t, builds[1], result[1]) // Last on default branch
}

func TestGetRunsToUpdate_DefaultBranch(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := &AdoUsesCases{Repository: mockRepo}

	builds := []model.PipelineRuns{
		createPipelineRun("refs/heads/main", "", 2),      // default branch
		createPipelineRun("refs/heads/feature-1", "", 1), // only one on feature-1
		createPipelineRun("refs/heads/main", "", 4),      // default branch
		createPipelineRun("refs/heads/main", "", 5),
	}

	mockRepo.On("GetRepositoryById", "repo-id").Return(model.Repository{DefaultBranch: "refs/heads/main"}, nil)

	result, err := uc.getRunsToUpdate(builds, "repo-id", 123)

	assert.NoError(t, err)
	assert.Equal(t, builds[0], result[0]) // Last on default ref
	assert.Equal(t, builds[2], result[1]) // Last on default branch
}

func TestGetRunsToUpdate_RepositoryError(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := &AdoUsesCases{Repository: mockRepo}

	builds := []model.PipelineRuns{
		createPipelineRun("refs/heads/feature-1", "", 1),
		createPipelineRun("refs/heads/main", "", 2),
	}

	mockRepo.On("GetRepositoryById", "repo-id").Return(model.Repository{}, errors.New("db error"))

	result, err := uc.getRunsToUpdate(builds, "repo-id", 123)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetAllWorkItemsToUpdatePrev_OnlyWorkItemWithUpperVersionThanBuild(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := &AdoUsesCases{Repository: mockRepo}

	builds := []model.WorkItem{
		{
			Id: 1,
			Fields: map[string]interface{}{
				"test": "25.5.8.5",
			},
		},
		{
			Id: 2,
			Fields: map[string]interface{}{
				"test": "25.5.5.8",
			},
		},
		{
			Id: 3,
			Fields: map[string]interface{}{
				"test": "25.5.5.3",
			},
		},
	}
	result := uc.getAllWorkItemsToUpdatePrev(builds, "25.5.5.5", "test")
	assert.Equal(t, 2, len(result))
	for _, res := range result {
		assert.Contains(t, []int{1, 2}, res.Id)
	}
}

func TestGetAllWorkItemsToUpdatePrev_ReturnZeroWorkItemWhenWorkItemVersionIsLower(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := &AdoUsesCases{Repository: mockRepo}

	builds := []model.WorkItem{
		{
			Id: 1,
			Fields: map[string]interface{}{
				"test": "25.5.4.5",
			},
		},
		{
			Id: 2,
			Fields: map[string]interface{}{
				"test": "25.5.4.8",
			},
		},
		{
			Id: 3,
			Fields: map[string]interface{}{
				"test": "25.5.4.3",
			},
		},
	}
	result := uc.getAllWorkItemsToUpdatePrev(builds, "25.5.5.5", "test")
	assert.Equal(t, 0, len(result))
}

func TestGetAllWorkItems_ReturnAllWorkItems(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := AdoUsesCases{Repository: mockRepo}

	mockRepo.On("GetBuildWorkItem", mock.Anything, mock.Anything).Return([]model.BuildWorkItems{
		{
			Id:  "1",
			Url: "url",
		},
		{
			Id:  "2",
			Url: "url",
		},
		{
			Id:  "3",
			Url: "url",
		},
	}, nil)

	mockRepo.On("GetWorkItem", mock.Anything).Return(model.WorkItem{}, nil)

	result, err := uc.getAllWorkItems([]model.PipelineRuns{{Id: 1}, {Id: 2}})

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result)
}

func TestGetAllWorkItems_ReturnErr(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := AdoUsesCases{Repository: mockRepo}

	mockRepo.On("GetBuildWorkItem", mock.Anything, mock.Anything).Return([]model.BuildWorkItems{}, errors.New("error"))

	result, err := uc.getAllWorkItems([]model.PipelineRuns{{Id: 1}, {Id: 2}})

	assert.NotNil(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func TestUpdateAdoIntegrationBuild_WithEmptyVersion(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := AdoUsesCases{Repository: mockRepo}

	version := "25.5.5.5"
	workItems := []model.WorkItem{
		createWorkItem(1, map[string]interface{}{AdoIntegrationBuildFieldName: ""}),
	}

	_ = mockRepo.On("UpdateWorkItemFields", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	err := uc.updateAdoIntegrationBuild(workItems, version)

	assert.Nil(t, err)
	mockRepo.AssertNotCalled(t, "UpdatetItemFields", mock.Anything, version, mock.Anything)
}

func TestUpdateAdoIntegrationBuild_WithNotEmptyVerison(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := AdoUsesCases{Repository: mockRepo}

	version := "25.5.5.5"
	result := "25.5.3.5 | 25.5.5.5"
	workItems := []model.WorkItem{
		createWorkItem(1, map[string]interface{}{AdoIntegrationBuildFieldName: "25.5.3.5"}),
	}

	_ = mockRepo.On("UpdateWorkItemFields", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	err := uc.updateAdoIntegrationBuild(workItems, version)

	assert.Nil(t, err)
	mockRepo.AssertNotCalled(t, "UpdatetItemFields", mock.Anything, result, mock.Anything)
}

func TestUpdateAdoIntegration_WithSomeVerion(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := AdoUsesCases{Repository: mockRepo}

	version := "25.5.5.5"
	result := "25.5.3.5 | 25.6.5.5 | 25.5.5.5"
	workItems := []model.WorkItem{
		createWorkItem(1, map[string]interface{}{AdoIntegrationBuildFieldName: "25.5.3.5 | 25.6.5.5"}),
	}

	_ = mockRepo.On("UpdateWorkItemFields", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	err := uc.updateAdoIntegrationBuild(workItems, version)

	assert.Nil(t, err)
	mockRepo.AssertNotCalled(t, "UpdatetItemFields", mock.Anything, result, mock.Anything)
}
