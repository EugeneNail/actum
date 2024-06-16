package records

import (
	"database/sql"
	"fmt"
	"github.com/EugeneNail/actum/internal/resource/activities"
	"github.com/EugeneNail/actum/internal/resource/photos"
	"github.com/EugeneNail/actum/internal/service/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/middleware/routing"
	"github.com/EugeneNail/actum/internal/service/response"
	"github.com/EugeneNail/actum/internal/service/validation"
	"net/http"
	"strconv"
	"time"
)

type Controller struct {
	db              *sql.DB
	recordDAO       *DAO
	activityDAO     *activities.DAO
	activityService *activities.Service
	recordService   *Service
	photoService    *photos.Service
}

func NewController(
	db *sql.DB,
	recordDAO *DAO,
	activityDAO *activities.DAO,
	activityService *activities.Service,
	recordService *Service,
	photoService *photos.Service,
) Controller {
	return Controller{db, recordDAO, activityDAO, activityService, recordService, photoService}
}

func (controller *Controller) Show(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)

	id, err := strconv.Atoi(routing.GetVariable(request, 0))
	if err != nil {
		response.Send(fmt.Errorf("RecordController.Show(): %w", err), http.StatusBadRequest)
		return
	}

	record, err := controller.recordDAO.Find(id)
	if err != nil {
		response.Send(fmt.Errorf("RecordController.Show(): %w", err), http.StatusInternalServerError)
		return
	}

	if record.Id == 0 {
		response.Send(fmt.Sprintf("Запись %d не найдена", id), http.StatusNotFound)
		return
	}

	user := jwt.GetUser(request)
	if record.UserId != user.Id {
		response.Send("Вы не можете использовать чужую запись.", http.StatusForbidden)
		return
	}

	activities, err := controller.recordService.FetchIdsOfActivities(record.Id)
	if err != nil {
		response.Send(fmt.Errorf("RecordController.Show(): %w", err), http.StatusInternalServerError)
		return
	}

	photos, err := controller.recordService.FetchNamesOfPhotos(record.Id)
	if err != nil {
		response.Send(fmt.Errorf("RecordController.Show(): %w", err), http.StatusInternalServerError)
		return
	}

	output := make(map[string]any, 7)
	output["date"] = record.Date.Format("2006-01-02")
	output["mood"] = record.Mood
	output["weather"] = record.Weather
	output["notes"] = record.Notes
	output["id"] = record.Id
	output["activities"] = activities
	output["photos"] = photos

	response.Send(output, http.StatusOK)
	log.Info("User", user.Id, "fetched record", record.Id)
}

type storeInput struct {
	Mood       int      `json:"mood" rules:"required|integer|min:1|max:5"`
	Weather    int      `json:"weather" rules:"required|integer|min:1|max:9"`
	Notes      string   `json:"notes" rules:"max:5000"`
	Date       string   `json:"date" rules:"required|date|today|after:2020-01-01"`
	Activities []int    `json:"activities" rules:"required"`
	Photos     []string `json:"photos"`
}

func (controller *Controller) Store(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)

	errors, input, err := validation.NewValidator[storeInput]().Validate(request)
	if err != nil {
		response.Send(fmt.Errorf("RecordController.Store(): %w", err), http.StatusBadRequest)
		return
	}

	if len(errors) > 0 {
		response.Send(errors, http.StatusUnprocessableEntity)
		return
	}

	user := jwt.GetUser(request)
	isDateTaken, err := controller.recordService.IsDateTaken(input.Date, user.Id)
	if err != nil {
		response.Send(fmt.Errorf("RecordController.Store(): %w", err), http.StatusInternalServerError)
		return
	}

	if isDateTaken {
		message := fmt.Sprintf("У вас уже есть запись на дату %s.", input.Date)
		response.Send(map[string]any{"date": message}, http.StatusConflict)
		return
	}

	allExist, missingActivities, err := controller.activityService.CheckExistence(input.Activities, user.Id)
	if err != nil {
		response.Send(fmt.Errorf("RecordController.Store(): %w", err), http.StatusInternalServerError)
		return
	}

	if !allExist {
		errors := map[string]any{"activities": fmt.Sprintf("Активности %v не найдены.", missingActivities)}
		response.Send(errors, http.StatusNotFound)
		return
	}

	allExist, missingPhotos, err := controller.photoService.CheckExistence(input.Photos, user.Id)
	if err != nil {
		response.Send(fmt.Errorf("RecordController.Store(): %w", err), http.StatusInternalServerError)
		return
	}

	if !allExist {
		errors := map[string]any{"photos": fmt.Sprintf("Фотографии %v не найдены.", missingPhotos)}
		response.Send(errors, http.StatusNotFound)
		return
	}

	record, err := New(input.Mood, input.Weather, input.Date, input.Notes, user.Id)
	if err != nil {
		response.Send(fmt.Errorf("RecordController.Store(): %w", err), http.StatusInternalServerError)
		return
	}

	if err := controller.recordDAO.Save(&record); err != nil {
		response.Send(fmt.Errorf("RecordController.Store(): %w", err), http.StatusInternalServerError)
		return
	}

	if err = controller.activityService.SyncRelations(record.Id, input.Activities); err != nil {
		response.Send(fmt.Errorf("RecordController.Store(): %w", err), http.StatusInternalServerError)
		return
	}

	if err = controller.photoService.SyncRelations(record.Id, input.Photos); err != nil {
		response.Send(fmt.Errorf("RecordController.Store(): %w", err), http.StatusInternalServerError)
		return
	}

	response.Send(record.Id, http.StatusCreated)
	log.Info("User", user.Id, "created record", record.Id, "for date", input.Date, "with", len(input.Activities), "activities")
}

type updateInput struct {
	Mood       int      `json:"mood" rules:"required|integer|min:1|max:5"`
	Weather    int      `json:"weather" rules:"required|integer|min:1|max:9"`
	Notes      string   `json:"notes" rules:"max:5000"`
	Activities []int    `json:"activities" rules:"required"`
	Photos     []string `json:"photos"`
}

func (controller *Controller) Update(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)

	id, err := strconv.Atoi(routing.GetVariable(request, 0))
	if err != nil {
		response.Send(fmt.Errorf("RecordController.Update(): %w", err), http.StatusBadRequest)
		return
	}

	record, err := controller.recordDAO.Find(id)
	if err != nil {
		response.Send(fmt.Errorf("RecordController.Update(): %w", err), http.StatusInternalServerError)
		return
	}

	if record.Id == 0 {
		response.Send(fmt.Sprintf("Запись %d не найдена", id), http.StatusNotFound)
		return
	}

	errors, input, err := validation.NewValidator[updateInput]().Validate(request)
	if err != nil {
		response.Send(fmt.Errorf("RecordController.Update(): %w", err), http.StatusBadRequest)
		return
	}

	if len(errors) > 0 {
		response.Send(errors, http.StatusUnprocessableEntity)
		return
	}

	user := jwt.GetUser(request)
	allExist, missingActivities, err := controller.activityService.CheckExistence(input.Activities, user.Id)
	if err != nil {
		response.Send(fmt.Errorf("RecordController.Update(): %w", err), http.StatusInternalServerError)
		return
	}

	if !allExist {
		errors := map[string]any{"activities": fmt.Sprintf("Активности %v не найдены", missingActivities)}
		response.Send(errors, http.StatusNotFound)
		return
	}

	allExist, missingPhotos, err := controller.photoService.CheckExistence(input.Photos, user.Id)
	if err != nil {
		response.Send(fmt.Errorf("RecordController.Store(): %w", err), http.StatusInternalServerError)
		return
	}

	if !allExist {
		errors := map[string]any{"photos": fmt.Sprintf("Фотографии %v не найдены.", missingPhotos)}
		response.Send(errors, http.StatusNotFound)
		return
	}

	recordBefore := record
	record.Notes = input.Notes
	record.Mood = input.Mood
	record.Weather = input.Weather
	if err := controller.recordDAO.Save(&record); err != nil {
		response.Send(fmt.Errorf("RecordController.Update(): %w", err), http.StatusInternalServerError)
		return
	}

	if err = controller.activityService.SyncRelations(record.Id, input.Activities); err != nil {
		response.Send(fmt.Errorf("RecordController.Store(): %w", err), http.StatusInternalServerError)
		return
	}

	if err = controller.photoService.SyncRelations(record.Id, input.Photos); err != nil {
		response.Send(fmt.Errorf("RecordController.Store(): %w", err), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
	log.Info("User", user.Id, "changed record", fmt.Sprintf("%+v", recordBefore), "to", fmt.Sprintf("%+v", record))
}

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
