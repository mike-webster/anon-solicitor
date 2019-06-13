package data

import (
	"errors"
	"fmt"
	"log"
	"time"

	// TODO: do I need both of these sql drivers...? Probs not...
	"github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
	domain "github.com/mike-webster/anon-solicitor/app"
	"github.com/mike-webster/anon-solicitor/env"
)

type dbEvent struct {
	ID          int64
	Title       string
	Description string
	Time        time.Time      `binding:"required"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   mysql.NullTime `db:"updated_at"`
	DeletedAt   mysql.NullTime `db:"deleted_at"`
}

type EventService struct {
	// TODO: make this private
	DB *sqlx.DB
}

func (es *EventService) CreateEvent(event *domain.Event) error {
	if event == nil {
		return errors.New("must pass event in order to create")
	}

	createdAt := time.Now().UTC()
	event.CreatedAt = createdAt
	event.UpdatedAt = &createdAt

	res, err := es.Conn().Exec("INSERT INTO events (title, description, time, created_at, updated_at) VALUES (?,?,?,?,?)",
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

	rows, err := es.Conn().Queryx("SELECT * FROM events WHERE ID = ?", id)
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
	var dbe []dbEvent
	err := es.Conn().Select(&dbe, "SELECT * FROM events")
	if err != nil {
		return nil, err
	}

	for _, e := range dbe {
		var del *time.Time
		var upd *time.Time
		if e.DeletedAt.Valid {
			val1, err := e.DeletedAt.Value()
			if err != nil {
				log.Println(fmt.Sprint("Error encountered: ", err))
			}
			t, ok := val1.(time.Time)
			if !ok {
				log.Println(fmt.Sprint("Couldnt cast value to time.Time, val: ", val1))
			}
			del = &t

			val2, err := e.UpdatedAt.Value()
			if err != nil {
				log.Println(fmt.Sprint("Error encountered: ", err))
			}
			t2, ok := val2.(time.Time)
			if !ok {
				log.Println(fmt.Sprint("Couldnt cast value to time.Time, val: ", val2))
			}
			upd = &t2
		}
		ret = append(ret, domain.Event{
			ID:          e.ID,
			Title:       e.Title,
			Description: e.Description,
			Time:        e.Time,
			CreatedAt:   e.CreatedAt,
			UpdatedAt:   upd,
			DeletedAt:   del,
		})
	}

	return &ret, nil
}

func (es *EventService) AddQuestion(q *domain.Question) error {
	if q == nil {
		return errors.New("must pass question in order to add")
	}

	createdAt := time.Now().UTC()
	q.CreatedAt = createdAt
	q.UpdatedAt = mysql.NullTime{Time: createdAt}

	res, err := es.Conn().Exec("INSERT INTO questions (event_id, content, answers, created_at, updated_at) VALUES (?,?,?,?,?)",
		q.EventID,
		q.Content,
		q.Answers,
		q.CreatedAt,
		q.UpdatedAt,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	log.Printf("-- newly created question id: %v", id)

	q.ID = id

	log.Printf("-- assigning question id: %v", q.ID)

	return nil
}
func (es *EventService) GetQuestion(ID int64) *domain.Question {
	if ID < 1 {
		return nil
	}

	rows, err := es.Conn().Queryx("SELECT * FROM questions WHERE ID = ?", ID)
	if err != nil {
		log.Printf("query error: %v", err)
		return nil
	}

	if rows.Next() {
		var ret domain.Question
		err = rows.StructScan(&ret)
		if err != nil {
			log.Printf("struct scan error: %v", err)
			return nil
		}

		return &ret
	}

	return nil
}

func (es *EventService) CanUserAnswerQuestion(ID int64, tok string) bool {
	if ID < 1 {
		return false
	}

	if len(tok) < 1 {
		return false
	}

	// this query will check to see if an answer record exists
	// tied to the given token for the given question
	query := `SELECT a.id, a.question_id, a.content, a.created_at
	    FROM feedback f 
		INNER JOIN events e 
		ON f.event_id = e.id 
		INNER JOIN questions q 
		ON q.event_id = e.id 
		INNER JOIN answers a 
		ON q.id = a.question_id 
		LEFT JOIN feedback_answers fa 
		ON f.id = fa.feedback_id 
		AND q.id = fa.question_id 
		WHERE f.tok = ? 
		AND q.id = ?
		AND fa.feedback_id is null`

	var ret []domain.Answer
	err := es.Conn().Select(&ret, query, tok, ID)
	if err != nil {
		log.Printf("query error: %v", err)
		return false
	}

	if len(ret) > 0 {
		// if there's an answer... we can't add another

		return false
	}

	return true
}

func (es *EventService) AddAnswer(a *domain.Answer, feedbackID int64) error {
	if a == nil {
		return errors.New("must pass answer in order to add")
	}

	createdAt := time.Now().UTC()
	a.CreatedAt = createdAt

	res, err := es.Conn().Exec("INSERT INTO answers (question_id, content, created_at) VALUES (?,?,?)",
		a.QuestionID,
		a.Content,
		a.CreatedAt,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	log.Printf("-- newly created answer id: %v", id)

	a.ID = id

	log.Printf("-- assigning question id: %v", a.ID)

	_, err = es.Conn().Exec("INSERT INTO feedback_answers (feedback_id, question_id, answer_id) VALUES (?,?,?)",
		feedbackID,
		a.QuestionID,
		a.ID,
	)
	if err != nil {
		return err
	}

	log.Println("-- newly created feedback_answer record")

	return nil
}

// Conn will get a database connection.
// This uses a cached pointer.
func (es *EventService) Conn() *sqlx.DB {
	if es.DB == nil {
		log.Println("No DB connection - establishing...")
		cfg := env.Config()
		db, err := sqlx.Open("mysql", cfg.ConnectionString)
		if err != nil {
			panic(fmt.Sprint("Couldn't load database; error", err))
		}
		es.DB = db
		return es.DB
	}

	err := es.DB.Ping()
	if err != nil {
		log.Println("No DB connection - ping failed/ establishing...")

		cfg := env.Config()
		db, err := sqlx.Open("mysql", cfg.ConnectionString)
		if err != nil {
			panic(fmt.Sprint("Couldn't establish database connection; err: ", err))
		}
		es.DB = db
		return es.DB
	}
	log.Println("DB connection - ping success!")

	return es.DB
}

type FeedbackService struct {
	DB *sqlx.DB
}

func (fs *FeedbackService) Conn() *sqlx.DB {
	if fs.DB == nil {
		log.Println("No DB connection - establishing...")
		cfg := env.Config()
		db, err := sqlx.Open("mysql", cfg.ConnectionString)
		if err != nil {
			panic(fmt.Sprint("Couldn't load database; error", err))
		}
		fs.DB = db
		return fs.DB
	}

	err := fs.DB.Ping()
	if err != nil {
		log.Println("No DB connection - ping failed/ establishing...")

		cfg := env.Config()
		db, err := sqlx.Open("mysql", cfg.ConnectionString)
		if err != nil {
			panic(fmt.Sprint("Couldn't establish database connection; err: ", err))
		}
		fs.DB = db
		return fs.DB
	}
	log.Println("DB connection - ping success!")

	return fs.DB
}

func (fs FeedbackService) CreateFeedback(feedback *domain.Feedback) error {
	if feedback == nil {
		return errors.New("must pass event in order to create")
	}

	res, err := fs.DB.Exec("INSERT INTO feedback  tok, event_id) VALUES (?,?)",
		feedback.Tok, feedback.EventID)
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

	err := fs.Conn().Select(&ret, "SELECT * FROM feedback where tok = ?", tok)
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

func (fs FeedbackService) GetQuestionsForTok(tok string) *[]domain.Question {
	if len(tok) < 1 {
		return nil
	}

	ret := []domain.Question{}

	query := `
		SELECT 
			q.id, 
			q.event_id,
			q.content,
			q.answers,
			q.created_at,
			q.updated_at,
			q.deleted_at
		FROM feedback f 
		INNER JOIN events e
		ON f.event_id = e.id 
		INNER JOIN questions q 
		ON q.event_id = e.id 
		LEFT JOIN answers a 
		ON a.question_id = q.id 
		WHERE f.tok = ?
		AND a.id IS NULL;`

	err := fs.Conn().Select(&ret, query, tok)
	if err != nil {
		log.Println("Error encountered: ", err, "\nquery: ", query)
	}

	return &ret
}
