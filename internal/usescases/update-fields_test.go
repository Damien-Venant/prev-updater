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
	return nil, nil
}
func (m *MockRepository) GetWorkitem(workItemId int) (*model.BuildWorkItems, error) {
	return nil, nil
}
func (m *MockRepository) UpdateWorkitemField(workItemId string, operation model.OperationFields) error {
	return nil
}

// Setup d’un run
func createPipelineRun(ref string) model.PipelineRuns {
	return model.PipelineRuns{
		State: "completed",
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

func TestGetRunsToUpdate_TwoBuildsOnSameRef(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := &AdoUsesCases{Repository: mockRepo}

	builds := []model.PipelineRuns{
		createPipelineRun("refs/heads/feature-1"),
		createPipelineRun("refs/heads/feature-1"),
		createPipelineRun("refs/heads/main"),
	}

	mockRepo.On("GetRepositoryById", "repo-id").Return(model.Repository{DefaultBranch: "refs/heads/main"}, nil)

	result, err := uc.getRunsToUpdate(builds, "repo-id", 123)

	assert.NoError(t, err)
	assert.Equal(t, builds[0], result[0])
	assert.Equal(t, builds[1], result[1])
}

func TestGetRunsToUpdate_OnlyOneBuildOnRef(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := &AdoUsesCases{Repository: mockRepo}

	builds := []model.PipelineRuns{
		createPipelineRun("refs/heads/feature-1"), // only one on feature-1
		createPipelineRun("refs/heads/main"),      // default branch
		createPipelineRun("refs/heads/main"),
	}

	mockRepo.On("GetRepositoryById", "repo-id").Return(model.Repository{DefaultBranch: "refs/heads/main"}, nil)

	result, err := uc.getRunsToUpdate(builds, "repo-id", 123)

	assert.NoError(t, err)
	assert.Equal(t, builds[0], result[0]) // Last on current ref
	assert.Equal(t, builds[1], result[1]) // Last on default branch
}

func TestGetRunsToUpdate_RepositoryError(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := &AdoUsesCases{Repository: mockRepo}

	builds := []model.PipelineRuns{
		createPipelineRun("refs/heads/feature-1"),
		createPipelineRun("refs/heads/main"),
	}

	mockRepo.On("GetRepositoryById", "repo-id").Return(model.Repository{}, errors.New("db error"))

	result, err := uc.getRunsToUpdate(builds, "repo-id", 123)

	assert.Error(t, err)
	assert.Nil(t, result)
}
