package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	cfg "github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/repository"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/utils"
	_ "github.com/jackc/pgx/v4/stdlib" // golint
	"github.com/pressly/goose/v3"
)

const (
	LeftSquare  = "["
	RightSquare = "]"
	COMMA       = ","
	FindAll     = `
SELECT id, title, lower(duration) as start_time, upper(duration) as end_time, description, user_id,
       notify_before_event FROM events
`
	FindByID = `
SELECT id, title, lower(duration) as start_time, upper(duration) as end_time, description, user_id,
       notify_before_event FROM events
WHERE id = $1
`
	DELETE = `
DELETE FROM events
WHERE id = $1
`
	INSERT = `
INSERT INTO events (title, duration, description, user_id, notify_before_event)
VALUES ($1, $2, $3, $4, $5) RETURNING id
`
	UPDATE = `
UPDATE events 
    SET title = $1, 
        duration = $2, 
        description = $3, 
        user_id = $4, 
        notify_before_event = $5
WHERE id = $6 RETURNING id
`
	CheckForIntersection = `
select count(t.b)
from (
select $1::tsrange && duration as b
from events
)t
where b = true
`
	CheckForIntersectionUpdate = `
select count(t.b)
from (
select $1::tsrange && duration as b
from events
where id != $2
)t
where b = true
`
	FindByRange = `
SELECT id, title, lower(duration) as start_time, upper(duration) as end_time, description, user_id,
       notify_before_event
FROM events
WHERE lower(duration) <@ $1::tsrange
`
)

type Storage struct { // TODO
	db        *sql.DB
	dsn       string
	migration string
}

func New(conf *cfg.DBConf) *Storage {
	return &Storage{dsn: getDsn(conf), migration: conf.Migration}
}

func (s *Storage) Connect(ctx context.Context) (err error) {
	s.db, err = sql.Open("pgx", s.dsn)
	if err != nil {
		return fmt.Errorf("cannot open pgx driver: %w", err)
	}
	return s.db.PingContext(ctx)
}

func (s *Storage) Migrate(ctx context.Context) (err error) {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("cannot set dialect: %w", err)
	}

	if err := goose.UpContext(ctx, s.db, s.migration); err != nil {
		return fmt.Errorf("cannot do up migration: %w", err)
	}

	return nil
}

func (s *Storage) Close() error {
	err := s.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Remove(ctx context.Context, id string) (err error) {
	fail := func(e error) error {
		return fmt.Errorf("DeleteEvent %w", e)
	}
	_, err = s.db.ExecContext(ctx, DELETE, id)
	return fail(err)
}

func (s *Storage) Find(ctx context.Context, id string) (*model.DBEvent, error) {
	fail := func(e error) (*model.DBEvent, error) {
		return nil, fmt.Errorf("FindByIdEvent %w", e)
	}

	row := s.db.QueryRowContext(ctx, FindByID, id)
	var startTime, endTime sql.NullTime
	var notify sql.NullInt64
	var descr sql.NullString
	var ev model.DBEvent

	err := row.Scan(
		&ev.ID,
		&ev.Title,
		&startTime,
		&endTime,
		&descr,
		&ev.UserID,
		&notify)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return fail(err)
	}
	if descr.Valid {
		ev.Description = descr.String
	}
	if startTime.Valid {
		ev.StartTime = startTime.Time
	}
	if endTime.Valid {
		ev.EndTime = endTime.Time
	}
	if notify.Valid {
		ev.NotifyBeforeEvent = time.Duration(notify.Int64)
	}
	return &ev, nil
}

func (s *Storage) FindAll(ctx context.Context) ([]*model.DBEvent, error) {
	fail := func(e error) ([]*model.DBEvent, error) {
		return nil, fmt.Errorf("FindAllEvents %w", e)
	}

	rows, err := s.db.QueryContext(ctx, FindAll)
	if err != nil {
		return fail(err)
	}
	defer rows.Close()

	events, er := findRequestHandler(rows)
	if er != nil {
		return fail(er)
	}
	return events, nil
}

func (s *Storage) FindAllByDay(ctx context.Context, date time.Time) ([]*model.DBEvent, error) {
	fail := func(e error) ([]*model.DBEvent, error) {
		return nil, fmt.Errorf("FindByDayRange %w", e)
	}
	b, e := utils.DayRange(date)
	r := makePostgresDuration(b, e)
	rows, err := s.db.QueryContext(ctx, FindByRange, r)
	if err != nil {
		return fail(err)
	}
	defer rows.Close()

	events, er := findRequestHandler(rows)
	if er != nil {
		return fail(er)
	}
	return events, nil
}

func (s *Storage) FindAllByWeek(ctx context.Context, date time.Time) ([]*model.DBEvent, error) {
	fail := func(e error) ([]*model.DBEvent, error) {
		return nil, fmt.Errorf("FindByDayRange %w", e)
	}
	b, e := utils.WeekRange(date)
	r := makePostgresDuration(b, e)
	rows, err := s.db.QueryContext(ctx, FindByRange, r)
	if err != nil {
		return fail(err)
	}
	defer rows.Close()

	events, er := findRequestHandler(rows)
	if er != nil {
		return fail(er)
	}
	return events, nil
}

func (s *Storage) FindAllByMonth(ctx context.Context, date time.Time) ([]*model.DBEvent, error) {
	fail := func(e error) ([]*model.DBEvent, error) {
		return nil, fmt.Errorf("FindByDayRange %w", e)
	}
	b, e := utils.MonthRange(date)
	r := makePostgresDuration(b, e)
	rows, err := s.db.QueryContext(ctx, FindByRange, r)
	if err != nil {
		return fail(err)
	}
	defer rows.Close()

	events, er := findRequestHandler(rows)
	if er != nil {
		return fail(er)
	}
	return events, nil
}

func (s *Storage) Update(ctx context.Context, event *model.DBEvent) (updatedEvent *model.DBEvent, err error) {
	fail := func(e error) error {
		return fmt.Errorf("UpdateEvent %w", e)
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fail(repository.ErrBeginTx)
	}
	defer tx.Rollback()
	var b bool
	var dur string
	b, dur, err = checkForIntersection(ctx, tx, event, CheckForIntersectionUpdate)
	if err != nil {
		return nil, fail(err)
	}
	if b {
		return nil, fail(repository.ErrDateBusy)
	}

	var eid string
	var notify interface{}
	var descr interface{}
	if event.NotifyBeforeEvent == 0 {
		notify = nil
	} else {
		notify = event.NotifyBeforeEvent
	}
	if len(event.Description) == 0 {
		descr = nil
	} else {
		descr = event.Description
	}
	err = tx.QueryRowContext(ctx, UPDATE, event.Title, dur, descr, event.UserID, notify,
		event.ID).Scan(&eid)
	if err != nil {
		return nil, fail(err)
	}
	var ev *model.DBEvent

	err = tx.Commit()
	if err != nil {
		return nil, fail(repository.ErrCommitTx)
	}

	ev, err = s.Find(ctx, eid)
	if err != nil {
		return nil, fail(err)
	}

	return ev, nil
}

func (s *Storage) Insert(ctx context.Context, event *model.DBEvent) (newEvent *model.DBEvent, err error) {
	fail := func(e error) error {
		return fmt.Errorf("CreateEvent %w", e)
	}

	if len(event.ID) != 0 {
		return nil, fail(repository.ErrIDExists)
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fail(repository.ErrBeginTx)
	}
	defer tx.Rollback()

	var b bool
	var dur string
	b, dur, err = checkForIntersection(ctx, tx, event, CheckForIntersection)
	if err != nil {
		return nil, fail(err)
	}
	if b {
		return nil, fail(repository.ErrDateBusy)
	}

	var eid string
	var notify interface{}
	if event.NotifyBeforeEvent == 0 {
		notify = nil
	} else {
		notify = event.NotifyBeforeEvent
	}
	err = tx.QueryRowContext(ctx, INSERT, event.Title, dur, event.Description, event.UserID, notify).Scan(&eid)
	if err != nil {
		return nil, fail(err)
	}
	var ev *model.DBEvent

	err = tx.Commit()
	if err != nil {
		return nil, fail(repository.ErrCommitTx)
	}

	ev, err = s.Find(ctx, eid)
	if err != nil {
		return nil, fail(err)
	}

	return ev, nil
}

func getDsn(conf *cfg.DBConf) string {
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		conf.User, conf.Password, conf.Dbname, conf.Host, conf.Port)
}

func checkForIntersection(ctx context.Context, tx *sql.Tx, event *model.DBEvent, query string) (bool, string, error) {
	dur := makePostgresDuration(event.StartTime, event.EndTime)
	var row *sql.Row
	if query == CheckForIntersection {
		row = tx.QueryRowContext(ctx, query, dur)
	} else {
		row = tx.QueryRowContext(ctx, query, dur, event.ID)
	}
	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, dur, err
	}
	if count > 0 {
		return true, dur, nil
	}
	return false, dur, nil
}

func makePostgresDuration(start time.Time, end time.Time) string {
	return fmt.Sprintf("%s%s%s%s%s", LeftSquare, start.Format(time.RFC3339), COMMA,
		end.Format(time.RFC3339), RightSquare)
}

func findRequestHandler(rows *sql.Rows) ([]*model.DBEvent, error) {
	var events []*model.DBEvent

	for rows.Next() {
		ev := model.NewDBEvent()

		var startTime, endTime sql.NullTime
		var descr sql.NullString
		var notify sql.NullInt64

		if err := rows.Scan(
			&ev.ID,
			&ev.Title,
			&startTime,
			&endTime,
			&descr,
			&ev.UserID,
			&notify,
		); err != nil {
			return nil, err
		}

		if descr.Valid {
			ev.Description = descr.String
		}
		if startTime.Valid {
			ev.StartTime = startTime.Time
		}
		if endTime.Valid {
			ev.EndTime = endTime.Time
		}
		if notify.Valid {
			ev.NotifyBeforeEvent = time.Duration(notify.Int64)
		}
		events = append(events, ev)
	}
	return events, nil
}
