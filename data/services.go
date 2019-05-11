package data

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	domain "github.com/mike-webster/anon-solicitor/app"
)

type EventService struct {
	DB *sqlx.DB
}

func (es *EventService) CreateEvent(event *domain.Event) error {
	if event == nil {
		return errors.New("must pass event in order to create")
	}

	createdAt := time.Now().UTC()
	event.CreatedAt = createdAt
	event.UpdatedAt = createdAt

	res, err := es.DB.Exec("INSERT INTO events (title, description, time, created_at, updated_at) VALUES (?,?,?,?,?)",
		event.Title,
		event.Description,
		event.Time,
		event.CreatedAt,
		event.UpdatedAt,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	log.Printf("-- newly created event id: %v", id)

	event.ID = id

	log.Printf("-- assigning event id: %v", event.ID)

	return nil
}

func (es *EventService) GetEvent(id int64) *domain.Event {
	if id < 1 {
		log.Print("id less than 1")
		return nil
	}

	rows, err := es.DB.Queryx("SELECT * FROM events WHERE ID = ?", id)
	if err != nil {
		log.Printf("query error: %v", err)
		return nil
	}

	if rows.Next() {
		var ret domain.Event
		err = rows.StructScan(&ret)
		if err != nil {
			log.Printf("struct scan error: %v", err)
			return nil
		}

		return &ret
	}

	return nil
}

func (es *EventService) GetEvents() (*[]domain.Event, error) {
	var ret []domain.Event

	err := es.DB.Select(&ret, "SELECT * FROM events")
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

type FeedbackService struct {
	DB *sqlx.DB
}

func (fs FeedbackService) CreateFeedback(feedback *domain.Feedback) error {
	if feedback == nil {
		return errors.New("must pass event in order to create")
	}

	res, err := fs.DB.Exec("INSERT INTO feedback (content, tok, event_id) VALUES (?,?,?)",
		feedback.Content, feedback.Tok, feedback.EventID)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	feedback.ID = id
	return nil
}

func (fs FeedbackService) GetFeedbackByTok(tok string) (*domain.Feedback, error) {
	if len(tok) < 1 {
		return nil, errors.New("please provide a token")
	}

	var ret []domain.Feedback

	err := fs.DB.Select(&ret, "SELECT * FROM feedback where tok = '?'", tok)
	if err != nil {
		return nil, err
	}

	if len(ret) < 1 {
		return nil, fmt.Errorf("no record found for tok: %v", tok)
	}

	return &ret[0], nil
}

func (fs FeedbackService) MarkFeedbackAbsent(f *domain.Feedback) error {
	if f == nil {
		return errors.New("please provide a record to mark as absent")
	}

	res, err := fs.DB.Exec("UPDATE feedback SET tok = '', absent=1 WHERE id = ?", f.ID)
	if err != nil {
		return err
	}

	num, _ := res.RowsAffected()
	if num != 1 {
		return fmt.Errorf("Unexpected number of results updated.... %v", num)
	}

	return nil
}
