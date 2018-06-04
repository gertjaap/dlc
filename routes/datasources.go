package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gertjaap/dlcoracle/datasources"
)

type DataSourceResponse struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Id           uint64 `json:"id"`
	CurrentValue uint64 `json:"currentValue"`
	ValueError   string `json:"valueError,omitempty"`
}

func ListDataSourcesHandler(w http.ResponseWriter, r *http.Request) {

	var ds = datasources.GetAllDatasources()

	response := []DataSourceResponse{}
	for _, src := range ds {
		value := src.Value

		jsonSrc := DataSourceResponse{
			Name:         src.Name,
			Description:  src.Description,
			Id:           src.Id,
			CurrentValue: value}

		response = append(response, jsonSrc)
	}

	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
