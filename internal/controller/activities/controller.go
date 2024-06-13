package activities

import (
	"database/sql"
	"github.com/EugeneNail/actum/internal/database/resource/activities"
	"github.com/EugeneNail/actum/internal/database/resource/collections"
)

type Controller struct {
	db              *sql.DB
	activityDAO     *activities.DAO
	collectionDAO   *collections.DAO
	activityService *activities.Service
}

func New(db *sql.DB, activityDAO *activities.DAO, collectionDAO *collections.DAO, activityService *activities.Service) (controller Controller) {
	controller.db = db
	controller.activityDAO = activityDAO
	controller.collectionDAO = collectionDAO
	controller.activityService = activityService

	return
}
