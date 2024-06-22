package photos

import (
	"encoding/base64"
	"fmt"
	"github.com/EugeneNail/actum/internal/database/resource/photos"
	"github.com/EugeneNail/actum/internal/service/auth/jwt"
	"github.com/EugeneNail/actum/internal/service/env"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/EugeneNail/actum/internal/service/response"
	"github.com/EugeneNail/actum/internal/service/uuid"
	"github.com/EugeneNail/actum/internal/service/validation"
	"net/http"
	"os"
	"path/filepath"
)

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

	user := jwt.GetUser(request)
	photo := photos.New(name, nil, user.Id)
	if err := controller.dao.Save(&photo); err != nil {
		response.Send(fmt.Errorf("photoController.Store: failed to save to database: %w", err), http.StatusInternalServerError)
		return
	}

	response.Send(name, http.StatusCreated)
	log.Info("User", user.Id, "uploaded photo", photo.Id, "with name", name)
}
