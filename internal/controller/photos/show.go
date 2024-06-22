package photos

import (
	"errors"
	"fmt"
	"github.com/EugeneNail/actum/internal/infrastructure/env"
	"github.com/EugeneNail/actum/internal/infrastructure/middleware/routing"
	"github.com/EugeneNail/actum/internal/infrastructure/response"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

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
