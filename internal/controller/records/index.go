package records

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/service/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/response"
	"github.com/EugeneNail/actum/internal/service/validation"
	"net/http"
	"time"
)

type indexInput struct {
	Cursor string `json:"cursor" rules:"required|date"`
}

func (controller *Controller) Index(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)

	errors, input, err := validation.NewValidator[indexInput]().Validate(request)
	if err != nil {
		response.Send(fmt.Errorf("RecordController.Index(): %w", err), http.StatusBadRequest)
		return
	}

	if len(errors) > 0 {
		response.Send(errors, http.StatusUnprocessableEntity)
		return
	}

	user := jwt.GetUser(request)
	cursor, err := time.Parse("2006-01-02", input.Cursor)
	if err != nil {
		response.Send(fmt.Errorf("RecordController.Index(): %w", err), http.StatusBadRequest)
		return
	}

	records, err := controller.recordService.CollectRecordsForCursor(cursor, user.Id)
	if err != nil {
		response.Send(fmt.Errorf("RecordController.Index(): %w", err), http.StatusInternalServerError)
		return
	}

	response.Send(records, http.StatusOK)
	log.Info("User", user.Id, "fetched", len(records), "records for 2 weeks from cursor", input.Cursor)
}
