package photos

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/service/fake"
	"github.com/EugeneNail/actum/internal/service/tests/startup"
	"net/http"
	"testing"
)

func TestStoreValidData(t *testing.T) {
	client, database, storage := startup.Photos(t)

	storage.AssertCount("/photos", 0)

	var filename string
	client.
		Post("/api/photos", fmt.Sprintf(`{
			"image": "%s"
		}`, fake.Base64Image())).
		AssertStatus(http.StatusCreated).
		ReadData(&filename)

	database.
		AssertCount("photos", 1).
		AssertHas("photos", map[string]any{
			"id":      1,
			"name":    filename,
			"user_id": 1,
		})

	storage.
		AssertCount("/photos", 1).
		AssertHas("/photos", filename)
}

func TestStoreInvalidData(t *testing.T) {
	client, database, storage := startup.Photos(t)

	storage.AssertCount("/photos", 0)

	client.
		Post("/api/photos", fmt.Sprintf(`{
			"image": "%si"
		}`, fake.Base64Image())).
		AssertStatus(http.StatusBadRequest)

	database.AssertCount("photos", 0)
	storage.AssertCount("/photos", 0)
}
