package collections

import (
	"github.com/EugeneNail/actum/internal/database/resource/collections"
	"github.com/EugeneNail/actum/internal/service/tests/startup"
	"net/http"
	"sort"
	"testing"
)

func TestIndexEmpty(t *testing.T) {
	client, database := startup.Collections(t)

	database.AssertCount("collections", 0)

	var endpointCollections []collections.Collection
	client.
		Get("/api/collections").
		AssertStatus(http.StatusOK).
		ReadData(&endpointCollections)

	if len(endpointCollections) != 0 {
		t.Errorf("Expected %d collections from endpoint, got %d", 0, len(endpointCollections))
	}
}

func TestIndexFew(t *testing.T) {
	performByCount(10, t)
}

func TestIndexSome(t *testing.T) {
	performByCount(100, t)
}

func TestIndexMany(t *testing.T) {
	performByCount(1000, t)
}

func performByCount(count int, t *testing.T) {
	client, database := startup.Collections(t)

	newCollections := collections.NewFactory(1).
		Make(count).
		Insert().
		List()

	database.AssertCount("collections", count)

	var endpointCollections []collections.Collection
	client.
		Get("/api/collections").
		AssertStatus(http.StatusOK).
		ReadData(&endpointCollections)

	if len(endpointCollections) != count {
		t.Errorf("Expected %d collections from endpoint, got %d", count, len(endpointCollections))
	}

	sort.Slice(endpointCollections, func(i int, j int) bool {
		return endpointCollections[i].Id < endpointCollections[j].Id
	})

	for i, collection := range endpointCollections {
		if collection.Name != newCollections[i].Name ||
			collection.UserId != 1 ||
			collection.Color != newCollections[i].Color {
			t.Errorf("Collection %+v must be %+v", collection, newCollections[i])
		}
	}
}
