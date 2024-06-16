package photos

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/EugeneNail/actum/internal/service/env"
	"github.com/EugeneNail/actum/internal/service/jwt"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/middleware/routing"
	"github.com/EugeneNail/actum/internal/service/response"
	"github.com/EugeneNail/actum/internal/service/uuid"
	"github.com/EugeneNail/actum/internal/service/validation"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Controller struct {
	dao *DAO
}

func NewController(dao *DAO) *Controller {
	return &Controller{dao}
}

type storeInput struct {
	Image string `json:"image" rules:"required"`
}

func (controller *Controller) Store(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)

	validationErrors, input, err := validation.NewValidator[storeInput]().Validate(request)
	if err != nil {
		response.Send(fmt.Errorf("PhotoController.Store: malformed input: %w", err), http.StatusBadRequest)
		return
	}

	if len(validationErrors) > 0 {
		response.Send(validationErrors, http.StatusUnprocessableEntity)
		return
	}

	directory := filepath.Join(env.Get("APP_PATH"), "storage", "photos")
	err = os.MkdirAll(directory, 0755)
	if err != nil {
		response.Send(fmt.Errorf("PhotoController.Store: failed to create directory %s: %w", directory, err), http.StatusInternalServerError)
		return
	}

	name := uuid.New() + ".png"
	filePath := filepath.Join(directory, name)
	file, err := os.Create(filePath)
	defer file.Close()

	imageBytes, err := base64.StdEncoding.DecodeString(input.Image)
	if err != nil {
		response.Send(fmt.Errorf("photoController.Store: failed to convert base64 to binary: %w", err), http.StatusBadRequest)
		if err := os.Remove(filePath); err != nil {
			response.Send(fmt.Errorf("photoController.Store: failed to delete unused file %s: %w", filePath, err), http.StatusInternalServerError)
		}
		return
	}

	_, err = file.Write(imageBytes)
	if err != nil {
		response.Send(fmt.Errorf("photoController.Store: failed to write to file %s: %w", filePath, err), http.StatusInternalServerError)
		return
	}

	userId := jwt.GetUserId(request)
	photo := New(name, nil, userId)
	if err := controller.dao.Save(&photo); err != nil {
		response.Send(fmt.Errorf("photoController.Store: failed to save to database: %w", err), http.StatusInternalServerError)
		return
	}

	response.Send(name, http.StatusCreated)
	log.Info("User", userId, "uploaded photo", photo.Id, "with name", name)
}

func (controller *Controller) Destroy(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)

	name := routing.GetVariable(request, 0)
	photo, err := controller.dao.FindBy("name", name)
	if err != nil {
		response.Send(fmt.Errorf("photoController.Destroy: failed to get the photo: %w", err), http.StatusInternalServerError)
		return
	}

	if photo.Id == 0 {
		response.Send(fmt.Sprintf("Фотография %d не найдена", photo.Id), http.StatusNotFound)
		return
	}

	userId := jwt.GetUserId(request)
	if photo.UserId != userId {
		response.Send("Вы не можете удалять чужие фотографии", http.StatusForbidden)
		return
	}

	if err := controller.dao.Delete(photo.Id); err != nil {
		response.Send(fmt.Errorf("photoController.Destroy: failed to delete the photo: %w", err), http.StatusInternalServerError)
	}

	filePath := filepath.Join(env.Get("APP_PATH"), "storage", "photos", name)
	if err := os.Remove(filePath); err != nil {
		response.Send(fmt.Errorf("photoController.Destroy: failed to delete file: %w", err), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
	log.Info("User", userId, "deleted photo", photo.Id)
}
func (controller *Controller) Show(writer http.ResponseWriter, request *http.Request) {
	response := response.NewSender(writer)
	name := routing.GetVariable(request, 0)

	parts := strings.Split(name, ".")
	if len(parts) <= 1 {
		response.Send("У файла отсутствует расширение", http.StatusBadRequest)
		return
	}

	path := filepath.Join(env.Get("APP_PATH"), "storage", "photos", name)
	if _, err := os.Stat(path); err != nil && errors.Is(err, os.ErrNotExist) {
		response.Send(fmt.Sprintf("Фотография %s не найдена", name), http.StatusNotFound)
		return
	}

	contentType := "image/" + parts[len(parts)-1]
	writer.Header().Set("Content-Type", contentType)
	http.ServeFile(writer, request, path)
}
