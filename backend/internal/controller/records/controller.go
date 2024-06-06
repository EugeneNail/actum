package records

import (
	"database/sql"
	"github.com/EugeneNail/actum/internal/database/resource/activities"
	"github.com/EugeneNail/actum/internal/database/resource/records"
)

type Controller struct {
	db              *sql.DB
	recordDAO       *records.DAO
	activityDAO     *activities.DAO
	activityService *activities.Service
}

func New(db *sql.DB, recordDAO *records.DAO, activityDAO *activities.DAO, activityService *activities.Service) Controller {
	return Controller{db, recordDAO, activityDAO, activityService}
}
