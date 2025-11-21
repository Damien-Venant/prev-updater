package usescases

import (
	"errors"
	"fmt"
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

type MockN8N struct {
	mock.Mock
}

func (m *MockRepository) GetRepositoryById(repositoryId string) (*model.Repository, error) {
	args := m.Called(repositoryId)
	val := args.Get(0).(model.Repository)
	return &val, args.Error(1)
}
func (m *MockRepository) GetPipelineRuns(pipelineId int) ([]model.PipelineRuns, error) {
	args := m.Called(pipelineId)
	val := args.Get(0).([]model.PipelineRuns)
	return val, args.Error(1)
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
	var val model.WorkItem
	args := m.Called(workItemId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	val, _ = args.Get(0).(model.WorkItem)
	return &val, args.Error(1)
}
func (m *MockRepository) UpdateWorkitemField(workItemId string, operation model.OperationFields) error {
	return nil
}

func (m *MockN8N) PostWebhook(data model.N8nResult) error {
	args := m.Called(data)
	return args.Error(0)
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

	result, err := uc.getRunsToUpdate(builds, "repo-id", 123, "")

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

	result, err := uc.getRunsToUpdate(builds, "repo-id", 123, "")

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

	result, err := uc.getRunsToUpdate(builds, "repo-id", 123, "")

	assert.NoError(t, err)
	assert.Equal(t, builds[0], result[0]) // Last on default ref
	assert.Equal(t, builds[2], result[1]) // Last on default branch
}
func TestGetRunsToUpdate_SpecificBranchName(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := &AdoUsesCases{Repository: mockRepo}

	builds := []model.PipelineRuns{
		createPipelineRun("refs/heads/main", "", 1),      // default branch
		createPipelineRun("refs/heads/feature-1", "", 2), // only one on feature-1
		createPipelineRun("refs/heads/feature-1", "", 3), // only one on feature-1
		createPipelineRun("refs/heads/main", "", 4),      // default branch
		createPipelineRun("refs/heads/main", "", 5),
	}

	mockRepo.On("GetRepositoryById", "repo-id").Return(model.Repository{DefaultBranch: "refs/heads/main"}, nil)

	result, err := uc.getRunsToUpdate(builds, "repo-id", 123, "feature-1")

	assert.NoError(t, err)
	assert.Equal(t, builds[1], result[0]) // Last on default ref
	assert.Equal(t, builds[2], result[1]) // Last on default branch
}
func TestGetRunsToUpdate_SpecificBranchNameWithOneRun(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := &AdoUsesCases{Repository: mockRepo}

	builds := []model.PipelineRuns{
		createPipelineRun("refs/heads/main", "", 1),      // default branch
		createPipelineRun("refs/heads/feature-1", "", 2), // only one on feature-1
		createPipelineRun("refs/heads/main", "", 3),      // default branch
		createPipelineRun("refs/heads/main", "", 4),
	}

	mockRepo.On("GetRepositoryById", "repo-id").Return(model.Repository{DefaultBranch: "refs/heads/main"}, nil)

	result, err := uc.getRunsToUpdate(builds, "repo-id", 123, "feature-1")

	assert.NoError(t, err)
	assert.Equal(t, builds[1], result[0]) // Last on default ref
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

	result, err := uc.getRunsToUpdate(builds, "repo-id", 123, "")

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
		{
			Id: 4,
			Fields: map[string]interface{}{
				"test": "25.5.10.3",
			},
		},
	}
	result := uc.getAllWorkItemsToUpdatePrev(builds, "25.5.5.5", "test")
	assert.Equal(t, 3, len(result))
	for _, res := range result {
		assert.Contains(t, []int{1, 2, 4}, res.Id)
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

func TestUpdateAdoIntegration_AlreadyContainsVersion_ShouldNotAddIt(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := AdoUsesCases{Repository: mockRepo}

	version := "25.5.5.5"
	result := "25.5.3.5 | 25.6.5.5 | 25.5.5.5"
	workItems := []model.WorkItem{
		createWorkItem(1, map[string]interface{}{AdoIntegrationBuildFieldName: "25.5.3.5 | 25.6.5.5 | 25.5.5.5"}),
	}

	_ = mockRepo.On("UpdateWorkItemFields", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	err := uc.updateAdoIntegrationBuild(workItems, version)

	assert.Nil(t, err)
	mockRepo.AssertNotCalled(t, "UpdatetItemFields", mock.Anything, result, mock.Anything)
}

func TestUpdateFieldsByLastRuns_WhenPipelineRunIsEmpty(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := AdoUsesCases{Repository: mockRepo}

	mockRepo.On("GetPipelineRuns", mock.Anything).Return([]model.PipelineRuns{}, nil)

	err := uc.UpdateFieldsByLastRuns(UpdateFieldsParams{
		PipelineId:   862,
		RepositoryId: "62",
		FieldName:    "Custom",
	})

	assert.Nil(t, err)
}

func TestUpdateFieldsByLastRuns_WhenPipelineRunReturnAnError(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := AdoUsesCases{Repository: mockRepo}

	mockRepo.On("GetPipelineRuns", mock.Anything).Return([]model.PipelineRuns{}, errors.New("err"))

	err := uc.UpdateFieldsByLastRuns(UpdateFieldsParams{
		PipelineId:   862,
		RepositoryId: "62",
		FieldName:    "Custom",
	})

	assert.NotNil(t, err)
	assert.Equal(t, "err", err.Error())
}

func TestUpdateFieldsByLastRuns(t *testing.T) {
	mockRepo := new(MockRepository)
	mockN8N := new(MockN8N)
	uc := AdoUsesCases{Repository: mockRepo, N8nRepo: mockN8N}

	pipelineRuns := []model.PipelineRuns{
		createPipelineRun("main", "25.6.5.0", 4),
		createPipelineRun("main", "25.6.5.1", 3),
		createPipelineRun("main", "25.6.5.2", 2),
		createPipelineRun("main", "25.6.5.3", 1),
	}
	buildWorkItems := []model.BuildWorkItems{
		{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"},
	}
	workItems := []model.WorkItem{
		createWorkItem(1, map[string]interface{}{"Custom": "", "System.Title": "", "System.Tags": "", "Microsoft.VSTS.Build.IntegrationBuild": ""}),
		createWorkItem(2, map[string]interface{}{"Custom": "", "System.Title": "", "System.Tags": "", "Microsoft.VSTS.Build.IntegrationBuild": ""}),
		createWorkItem(3, map[string]interface{}{"Custom": "", "System.Title": "", "System.Tags": "", "Microsoft.VSTS.Build.IntegrationBuild": ""}),
		createWorkItem(4, map[string]interface{}{"Custom": "", "System.Title": "", "System.Tags": "", "Microsoft.VSTS.Build.IntegrationBuild": ""}),
	}
	mockRepo.On("GetPipelineRuns", mock.Anything).Return(pipelineRuns, nil)
	mockRepo.On("GetRepositoryById", mock.Anything).Return(model.Repository{Id: "1", DefaultBranch: "main", Url: ""}, nil)
	mockRepo.On("GetBuildWorkItem", 3, 4).Return(buildWorkItems, nil)
	for _, workItem := range workItems {
		mockRepo.On("GetWorkItem", fmt.Sprintf("%d", workItem.Id)).Return(workItem, nil)
	}
	mockRepo.On("UpdateWorkItemField", mock.Anything, mock.Anything)
	mockN8N.On("PostWebhook", mock.Anything).Return(nil)

	err := uc.UpdateFieldsByLastRuns(UpdateFieldsParams{
		PipelineId:   862,
		RepositoryId: "62",
		FieldName:    "Custom",
	})

	assert.Nil(t, err)
}

func TestUpdateFieldsByLastRuns_ShouldReturnError_OnPipelineRuns(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := AdoUsesCases{Repository: mockRepo}

	pipelineRuns := []model.PipelineRuns{
		createPipelineRun("main", "25.6.5.0", 4),
		createPipelineRun("main", "25.6.5.1", 3),
		createPipelineRun("main", "25.6.5.2", 2),
		createPipelineRun("main", "25.6.5.3", 1),
	}
	buildWorkItems := []model.BuildWorkItems{
		{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"},
	}
	workItems := []model.WorkItem{
		createWorkItem(1, map[string]interface{}{"Custom": ""}),
		createWorkItem(2, map[string]interface{}{"Custom": ""}),
		createWorkItem(3, map[string]interface{}{"Custom": ""}),
		createWorkItem(4, map[string]interface{}{"Custom": ""}),
	}
	mockRepo.On("GetPipelineRuns", mock.Anything).Return(pipelineRuns, errors.New("error"))
	mockRepo.On("GetRepositoryById", mock.Anything).Return(model.Repository{Id: "1", DefaultBranch: "main", Url: ""}, nil)
	mockRepo.On("GetBuildWorkItem", 3, 4).Return(buildWorkItems, nil)
	for _, workItem := range workItems {
		mockRepo.On("GetWorkItem", fmt.Sprintf("%d", workItem.Id)).Return(workItem, nil)
	}
	mockRepo.On("UpdateWorkItemField", mock.Anything, mock.Anything)

	err := uc.UpdateFieldsByLastRuns(UpdateFieldsParams{
		PipelineId:   862,
		RepositoryId: "62",
		FieldName:    "Custom",
	})

	assert.NotNil(t, err)
	assert.Equal(t, "error", err.Error())
}
func TestUpdateFieldsByLastRuns_ShouldReturnError_OnRepositoryId(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := AdoUsesCases{Repository: mockRepo}

	pipelineRuns := []model.PipelineRuns{
		createPipelineRun("main", "25.6.5.0", 4),
		createPipelineRun("main", "25.6.5.1", 3),
		createPipelineRun("main", "25.6.5.2", 2),
		createPipelineRun("main", "25.6.5.3", 1),
	}
	buildWorkItems := []model.BuildWorkItems{
		{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"},
	}
	workItems := []model.WorkItem{
		createWorkItem(1, map[string]interface{}{"Custom": ""}),
		createWorkItem(2, map[string]interface{}{"Custom": ""}),
		createWorkItem(3, map[string]interface{}{"Custom": ""}),
		createWorkItem(4, map[string]interface{}{"Custom": ""}),
	}
	mockRepo.On("GetPipelineRuns", mock.Anything).Return(pipelineRuns, nil)
	mockRepo.On("GetRepositoryById", mock.Anything).Return(model.Repository{Id: "1", DefaultBranch: "main", Url: ""}, errors.New("error"))
	mockRepo.On("GetBuildWorkItem", 3, 4).Return(buildWorkItems, nil)
	for _, workItem := range workItems {
		mockRepo.On("GetWorkItem", fmt.Sprintf("%d", workItem.Id)).Return(workItem, nil)
	}
	mockRepo.On("UpdateWorkItemField", mock.Anything, mock.Anything)

	err := uc.UpdateFieldsByLastRuns(UpdateFieldsParams{
		PipelineId:   862,
		RepositoryId: "62",
		FieldName:    "Custom",
	})

	assert.NotNil(t, err)
	assert.Equal(t, "error", err.Error())
}
func TestUpdateFieldsByLastRuns_ShouldReturnError_OnBuildWorkItem(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := AdoUsesCases{Repository: mockRepo}

	pipelineRuns := []model.PipelineRuns{
		createPipelineRun("main", "25.6.5.0", 4),
		createPipelineRun("main", "25.6.5.1", 3),
		createPipelineRun("main", "25.6.5.2", 2),
		createPipelineRun("main", "25.6.5.3", 1),
	}
	buildWorkItems := []model.BuildWorkItems{
		{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"},
	}
	mockRepo.On("GetPipelineRuns", mock.Anything).Return(pipelineRuns, nil)
	mockRepo.On("GetRepositoryById", mock.Anything).Return(model.Repository{Id: "1", DefaultBranch: "main", Url: ""}, nil)
	mockRepo.On("GetBuildWorkItem", 3, 4).Return(buildWorkItems, errors.New("error"))
	mockRepo.On("UpdateWorkItemField", mock.Anything, mock.Anything)

	err := uc.UpdateFieldsByLastRuns(UpdateFieldsParams{
		PipelineId:   862,
		RepositoryId: "62",
		FieldName:    "Custom",
	})

	assert.NotNil(t, err)
	assert.Equal(t, "error", err.Error())
}

func TestIsSmallerThan_ShouldReturnOne_WhenIsSmaller(t *testing.T) {
	tests := []struct {
		Actual Version
		Target Version
		Result int
	}{
		{
			Actual: newVersion("25.5.5.5"),
			Target: newVersion("25.5.5.6"),
			Result: 1,
		},
		{
			Actual: newVersion("25.5.5.5"),
			Target: newVersion("25.5.6.5"),
			Result: 1,
		},
		{
			Actual: newVersion("25.5.5.5"),
			Target: newVersion("25.6.5.5"),
			Result: 1,
		},
		{
			Actual: newVersion("25.5.5.5"),
			Target: newVersion("26.5.5.5"),
			Result: 1,
		},
	}

	for index, test := range tests {
		name := fmt.Sprintf("TestIsSmallerThan_ShouldReturnOne_WhenIsSmaller_%d", index)
		t.Run(name, func(t *testing.T) {
			result := test.Actual.isSmallerThan(test.Target)
			assert.Equal(t, test.Result, result)
		})
	}
}

func TestIsSmallerThan_ShouldReturnMinusOne_WhenIsUpper(t *testing.T) {
	tests := []struct {
		Actual Version
		Target Version
		Result int
	}{
		{
			Actual: newVersion("25.5.5.6"),
			Target: newVersion("25.5.5.5"),
			Result: -1,
		},
		{
			Actual: newVersion("25.5.6.5"),
			Target: newVersion("25.5.5.5"),
			Result: -1,
		},
		{
			Actual: newVersion("25.6.5.5"),
			Target: newVersion("25.5.5.5"),
			Result: -1,
		},
		{
			Actual: newVersion("26.5.5.5"),
			Target: newVersion("25.5.5.5"),
			Result: -1,
		},
	}

	for index, test := range tests {
		name := fmt.Sprintf("TestIsSmallerThan_ShouldReturnMinusOne_WhenIsUpper_%d", index)
		t.Run(name, func(t *testing.T) {
			result := test.Actual.isSmallerThan(test.Target)
			assert.Equal(t, test.Result, result)
		})
	}
}

func TestIsSmallerThan_ShouldReturnZero_WhenIsEqual(t *testing.T) {
	actual := newVersion("25.5.5.5")
	target := newVersion("25.5.5.5")

	result := actual.isSmallerThan(target)
	assert.Equal(t, 0, result)
}

func TestIsHigherThan_ShouldReturnMinusOne_WhenIsSmaller(t *testing.T) {
	tests := []struct {
		Actual Version
		Target Version
		Result int
	}{
		{
			Actual: newVersion("25.5.5.5"),
			Target: newVersion("25.5.5.6"),
			Result: -1,
		},
		{
			Actual: newVersion("25.5.5.5"),
			Target: newVersion("25.5.6.5"),
			Result: -1,
		},
		{
			Actual: newVersion("25.5.5.5"),
			Target: newVersion("25.6.5.5"),
			Result: -1,
		},
		{
			Actual: newVersion("25.5.5.5"),
			Target: newVersion("26.5.5.5"),
			Result: -1,
		},
	}

	for index, test := range tests {
		name := fmt.Sprintf("TestIsHigherThan_ShouldReturnMinusOne_WhenIsSmaller_%d", index)
		t.Run(name, func(t *testing.T) {
			result := test.Actual.isHigherThan(test.Target)
			assert.Equal(t, test.Result, result)
		})
	}
}

func TestIsHigherThan_ShouldReturnZero_WhenIsUpper(t *testing.T) {
	tests := []struct {
		Actual Version
		Target Version
		Result int
	}{
		{
			Actual: newVersion("25.5.5.6"),
			Target: newVersion("25.5.5.5"),
			Result: 1,
		},
		{
			Actual: newVersion("25.5.6.5"),
			Target: newVersion("25.5.5.5"),
			Result: 1,
		},
		{
			Actual: newVersion("25.6.5.5"),
			Target: newVersion("25.5.5.5"),
			Result: 1,
		},
		{
			Actual: newVersion("26.5.5.5"),
			Target: newVersion("25.5.5.5"),
			Result: 1,
		},
	}

	for index, test := range tests {
		name := fmt.Sprintf("TestIsHigherThan_ShouldReturnOne_WhenIsUpper_%d", index)
		t.Run(name, func(t *testing.T) {
			result := test.Actual.isHigherThan(test.Target)
			assert.Equal(t, test.Result, result)
		})
	}
}

func TestIsHigherThan_ShouldReturnZero_WhenIsEqual(t *testing.T) {
	version1 := newVersion("25.5.5.5")
	version2 := newVersion("25.5.5.5")

	result := version1.isSmallerThan(version2)
	assert.Equal(t, 0, result)
}
