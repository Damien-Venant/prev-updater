package repository

import (
	"encoding/json"
	"net/http"

	"github.com/Damien-Venant/prev-updater/internal/model"
	httpclient "github.com/Damien-Venant/prev-updater/pkg/http-client"
)

type N8nRepository struct {
	client httpclient.HttpClient
}

func NewN8nRepository(client httpclient.HttpClient) *N8nRepository {
	return &N8nRepository{
		client: client,
	}
}

func (repo *N8nRepository) PostWebhook(data model.N8nResult) error {
	model, err := json.Marshal(data)
	if err != nil {
		return err
	}

	httpResponse, err := repo.client.Post("", model, nil)
	if err != nil {
		return err
	}

	if err = treatResult(httpResponse, http.StatusOK); err != nil {
		return err
	}
	return nil
}
