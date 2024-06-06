package records

import (
	"context"
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

	err := dao.db.QueryRow(`SELECT * FROM records WHERE id = ?`, id).
		Scan(&record.Id, &record.Mood, &record.Date, &record.Notes, &record.UserId)

	if err != nil && err != sql.ErrNoRows {
		return record, fmt.Errorf("records.Find(): %w", err)
	}

	return record, nil
}

func (dao *DAO) Save(record *Record) error {
	result, err := dao.db.Exec(`
		INSERT INTO records
		    (id, mood, date, notes, user_id)
		VALUES 
		    (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		    id = VALUES(id),
		    mood = VALUES(mood),
		    notes = VALUES(notes),
		    date = VALUES(date),
		    user_id = VALUES(user_id)
	`, record.Id, record.Mood, record.Date, record.Notes, record.UserId)

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

func (dao *DAO) SyncRelations(recordId int, activities []int) error {
	tx, err := dao.db.BeginTx(context.Background(), &sql.TxOptions{})
	defer tx.Rollback()
	if err != nil {
		return fmt.Errorf("records.SyncRelations(): %w", err)
	}

	if err = deleteUnusedRelations(tx, recordId, activities); err != nil {
		return fmt.Errorf("records.SyncRelations(): %w", err)
	}

	if err = upsertRelations(tx, recordId, activities); err != nil {
		return fmt.Errorf("records.SyncRelations(): %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("records.SyncRelations(): %w", err)
	}

	return nil
}

func deleteUnusedRelations(tx *sql.Tx, recordId int, activities []int) error {
	var placeholders string
	values := make([]any, len(activities)+1)
	values[0] = recordId

	for i, activityId := range activities {
		values[i+1] = activityId
		placeholders += "?,"
	}
	placeholders = "(" + placeholders[:len(placeholders)-1] + ")"

	_, err := tx.Exec(
		`DELETE FROM records_activities WHERE record_id = ? AND activity_id NOT IN `+placeholders,
		values...,
	)

	if err != nil {
		return fmt.Errorf("deleteUnusedRelations(): %w", err)
	}

	return nil
}

func upsertRelations(tx *sql.Tx, recordId int, activities []int) error {
	const columnsCount = 2
	var placeholders string
	values := make([]any, columnsCount*len(activities))

	for i, activityId := range activities {
		values[columnsCount*i+0] = recordId
		values[columnsCount*i+1] = activityId
		placeholders += "(?, ?),"
	}
	placeholders = placeholders[:len(placeholders)-1]

	_, err := tx.Exec(`
		INSERT INTO records_activities (record_id, activity_id) 
		VALUES `+placeholders+` 
		ON DUPLICATE KEY UPDATE record_id = VALUES(record_id)
	`, values...,
	)

	if err != nil {
		return fmt.Errorf("upsertRelations(): %w", err)
	}

	return nil
}
