package records

import (
	"database/sql"
	"fmt"
)

type DAO struct {
	db *sql.DB
}

func NewDAO(db *sql.DB) *DAO {
	return &DAO{db}
}

func (dao *DAO) Find(id int) (Record, error) {
	var record Record

	err := dao.db.QueryRow(`SELECT id, mood, weather, date, notes, user_id FROM records WHERE id = ?`, id).
		Scan(&record.Id, &record.Mood, &record.Weather, &record.Date, &record.Notes, &record.UserId)

	if err != nil && err != sql.ErrNoRows {
		return record, fmt.Errorf("records.Find(): %w", err)
	}

	return record, nil
}

func (dao *DAO) Save(record *Record) error {
	result, err := dao.db.Exec(`
		INSERT INTO records
		    (id, mood, weather, date, notes, user_id)
		VALUES 
		    (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		    id = VALUES(id),
		    mood = VALUES(mood),
		    weather = VALUES(weather),
		    notes = VALUES(notes),
		    date = VALUES(date),
		    user_id = VALUES(user_id)
	`, record.Id, record.Mood, record.Weather, record.Date, record.Notes, record.UserId)

	if err != nil {
		return fmt.Errorf("records.Save(): %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("records.Save(): %w", err)
	}

	if id != 0 {
		record.Id = int(id)
	}

	return nil
}

func (dao *DAO) Delete(id int) error {
	_, err := dao.db.Exec(`DELETE FROM records WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("records.Delete(): %w", err)
	}

	return nil
}
