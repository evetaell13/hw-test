package sqlstorage

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Storage struct {
	dsn string // TODO close file
	db  *sql.DB
}

func New() *Storage {
	return &Storage{}
}
func (s *Storage) Connect(ctx context.Context) error { // TODO use ctx
	db, err := sql.Open("pgx", s.dsn)
	if err != nil {
		log.Fatalf("failed to load driver: %v", err)
		return err
	}
	s.db = db
	return nil
}

func (s *Storage) Close(ctx context.Context) error { // TODO use ctx
	s.db.Close()
	return nil
}

func (s *Storage) GetEventsByPeriod(start time.Time, finish time.Time) ([]storage.Event, error) {
	query := `SELECT events.id, events.title, events.start,
						events.finish, events.description, events.user_id
				FROM events 
				WHERE events.start BETWEEN $1 AND $2 
						AND 
						events.finish BETWEEN $1 AND $2`
	rows, err := s.db.Query(query, start, finish)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	events := make([]storage.Event, 0, 1)
	for rows.Next() {
		var event storage.Event
		if err := rows.Scan(&event.ID, &event.Title, &event.Start,
			&event.Finish, &event.Description, &event.UserID); err != nil {
			return events, err
		}
		events = append(events, event)
	}
	if err = rows.Err(); err != nil {
		return events, err
	}
	return events, nil
}

func (s *Storage) GetEvent(id string) (storage.Event, error) {
	var event storage.Event
	query := `SELECT events.id, events.title, events.start,
				events.finish, events.description, events.user_id
				FROM events 
				WHERE events.id = $1`
	if err := s.db.QueryRow(query, id).Scan(
		&event.ID, &event.Title, &event.Start,
		&event.Finish, &event.Description, &event.UserID,
	); err != nil {
		return event, err
	}
	return event, nil
}

func (s *Storage) GetUser(id string) (storage.User, error) {
	var user storage.User
	query := `SELECT users.id, users.name
				FROM users 
				WHERE users.id = $1`
	if err := s.db.QueryRow(query, id).Scan(
		&user.ID, &user.Name,
	); err != nil {
		return user, err
	}
	return user, nil
}

func (s *Storage) WriteEvent(event storage.Event) error {
	sqlStatement := `
	INSERT INTO event (title, start, finish, description, user_id)
	VALUES ($1, $2, $3, $4, $5)`
	_, err := s.db.Exec(sqlStatement, event.Title, event.Start, event.Finish, event.Description, event.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) WriteUser(user storage.User) error {
	sqlStatement := `
	INSERT INTO user (name)
	VALUES ($1)`
	_, err := s.db.Exec(sqlStatement, user.Name)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) UpdateEvent(event storage.Event) error {
	sqlStatement := `
	UPDATE users
	SET title = $2, start = $3, finish = $4, description = $5, user_id = $6
	WHERE id = $1;`
	_, err := s.db.Exec(sqlStatement, event.ID, event.Title, event.Start,
		event.Finish, event.Description, event.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) UpdateUser(user storage.User) error {
	sqlStatement := `
	UPDATE users
	SET name = $2
	WHERE id = $1;`
	_, err := s.db.Exec(sqlStatement, user.ID, user.Name)
	if err != nil {
		return err
	}
	return nil
}
