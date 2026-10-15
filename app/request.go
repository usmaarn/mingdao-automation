package app

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type RequestData struct {
	Summary   string `json:"summary"`
	Solution  string `json:"solution"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Type      string `json:"type"`
}

func (d RequestData) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Summary, validation.Required),
		validation.Field(&d.Solution, validation.Required),
		validation.Field(&d.StartDate,
			validation.Required,
			validation.Date("2006-01-02 15:04"),
		),
		validation.Field(&d.EndDate,
			validation.Required,
			validation.Date("2006-01-02 15:04"),
		),
		validation.Field(
			&d.Type,
			validation.Required,
			validation.In(
				"issue locating",
				"documentantion",
				"documentation",
				"phone call",
				"meeting",
			)),
	)
}

type RequestPayload struct {
	MerchantName string        `json:"merchantName"`
	Data         []RequestData `json:"data"`
}

func (r RequestPayload) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.MerchantName),
		validation.Field(
			&r.Data,
			validation.Required,
			validation.Length(1, 0),
			validation.Each(
				validation.By(func(value interface{}) error {
					data, ok := value.(RequestData)
					if !ok {
						return errors.New("invalid data item")
					}

					return data.Validate()
				}),
			),
		),
	)
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req RequestPayload

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		if vErr, ok := errors.AsType[validation.Errors](err); ok {
			json.NewEncoder(w).Encode(vErr)
			return
		}

		http.Error(w, err.Error(), 500)
	}

	go func(req RequestPayload) {
		app.mutex.Lock()
		defer app.mutex.Unlock()

		for _, data := range req.Data {
			if strings.ToLower(data.Type) == "documentation" {
				data.Type = "documentantion"
			}
			record := Record{
				Summary:      data.Summary,
				MerchantName: req.MerchantName,
				StartDate:    data.StartDate,
				EndDate:      data.EndDate,
				Type:         strings.ToLower(data.Type),
				Solution:     data.Solution,
			}

			app.AddRecord(record)
		}
	}(req)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": "queued",
	})
}
