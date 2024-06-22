package photos

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/infrastructure/fake"
	"github.com/EugeneNail/actum/internal/infrastructure/tests/startup"
	"net/http"
	"testing"
)

func TestDestroy(t *testing.T) {
	client, database, storage := startup.Photos(t)

	var name string
	client.
		Post("/api/photos", fmt.Sprintf(`{
			"image": "%s"
		}`, fake.Base64Image())).
		AssertStatus(http.StatusCreated).
		ReadData(&name)

	database.
		AssertCount("photos", 1).
		AssertHas("photos", map[string]any{
			"name":    name,
			"user_id": 1,
		})

	storage.
		AssertCount("/photos", 1).
		AssertHas("/photos", name)

	client.
		Delete("/api/photos/" + name).
		AssertStatus(http.StatusNoContent)

	database.AssertCount("photos", 0)
	storage.AssertCount("/photos", 0)
}

func TestDestroyNonexistent(t *testing.T) {
	client, database, storage := startup.Photos(t)

	database.AssertCount("photos", 0)
	storage.AssertCount("/photos", 0)

	client.
		Delete("/api/photos/file.jpg").
		AssertStatus(http.StatusNotFound)
}
