package photos

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/infrastructure/env"
	"github.com/EugeneNail/actum/internal/infrastructure/log"
	"github.com/EugeneNail/actum/internal/infrastructure/middleware/routing"
	"github.com/EugeneNail/actum/internal/infrastructure/response"
	"github.com/EugeneNail/actum/internal/service/auth/jwt"
	"net/http"
	"os"
	"path/filepath"
)

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

	user := jwt.GetUser(request)
	if photo.UserId != user.Id {
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
	log.Info("User", user.Id, "deleted photo", photo.Id)
}
