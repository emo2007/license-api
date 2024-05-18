package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/emo2007/block-accounting/examples/license-api/internal/pkg/models"
	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

type ListMusiciansParams struct {
	Ids    uuid.UUIDs
	Name   string
	FromId uuid.UUID
}

func (r *Repository) ListMusicians(ctx context.Context, params ListMusiciansParams) ([]models.Musician, error) {
	ms := make([]models.Musician, 0, len(params.Ids))
	if err := r.Transaction(ctx, func(ctx context.Context) error {
		query := sq.Select(
			"id",
			"name",
		).From("musicians").PlaceholderFormat(sq.Dollar)

		if len(params.Ids) > 0 {
			query = query.Where(sq.Eq{
				"id": params.Ids,
			})
		}

		if params.FromId != uuid.Nil {
			query = query.Where(sq.Gt{
				"id": params.FromId,
			})
		}

		if params.Name != "" {
			query = query.Where(sq.Eq{
				"name": params.Name,
			})
		}

		rows, err := query.RunWith(r.Conn(ctx)).QueryContext(ctx)
		if err != nil {
			return err
		}

		defer rows.Close()

		for rows.Next() {
			var (
				id   uuid.UUID
				name string
			)

			if err := rows.Scan(&id, &name); err != nil {
				return err
			}

			ms = append(ms, models.Musician{
				ID:   id,
				Name: name,
			})
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ms, nil
}

func (r *Repository) New(ctx context.Context, m models.Musician) error {
	if err := r.Transaction(ctx, func(ctx context.Context) error {
		query := sq.Insert("musicians").Columns("id", "name").Values(
			m.ID,
			m.Name,
		).PlaceholderFormat(sq.Dollar)

		if _, err := query.RunWith(r.Conn(ctx)).ExecContext(ctx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (r *Repository) AddTrack(ctx context.Context, t models.Track) error {
	if err := r.Transaction(ctx, func(ctx context.Context) error {
		query := sq.Insert("musicians_tracks").Columns("id", "title", "played_times").Values(
			t.ID,
			t.Title,
			t.Played,
		).PlaceholderFormat(sq.Dollar)

		if _, err := query.RunWith(r.Conn(ctx)).ExecContext(ctx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

type ListTracksParams struct {
	Ids    uuid.UUIDs
	Title  string
	FromId uuid.UUID
}

func (r *Repository) ListTracks(ctx context.Context, params ListTracksParams) ([]models.Track, error) {
	t := make([]models.Track, 0, len(params.Ids))
	if err := r.Transaction(ctx, func(ctx context.Context) error {
		query := sq.Select(
			"id",
			"title",
			"played_times",
		).From("musicians_tracks").PlaceholderFormat(sq.Dollar)

		if len(params.Ids) > 0 {
			query = query.Where(sq.Eq{
				"id": params.Ids,
			})
		}

		if params.FromId != uuid.Nil {
			query = query.Where(sq.Gt{
				"id": params.FromId,
			})
		}

		if params.Title != "" {
			query = query.Where(sq.Eq{
				"title": params.Title,
			})
		}

		rows, err := query.RunWith(r.Conn(ctx)).QueryContext(ctx)
		if err != nil {
			return err
		}

		defer rows.Close()

		for rows.Next() {
			var (
				id     uuid.UUID
				title  string
				played int64
			)

			if err := rows.Scan(&id, &title, &played); err != nil {
				return err
			}

			t = append(t, models.Track{
				ID:     id,
				Title:  title,
				Played: played,
			})
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return t, nil
}

type AddPlaysParams struct {
	TrackID    uuid.UUID
	MusicianID uuid.UUID
	Plays      int64
	Month      uint8
	Year       uint32
}

func (r *Repository) AddPlays(ctx context.Context, params AddPlaysParams) error {
	if err := r.Transaction(ctx, func(ctx context.Context) error {
		lockQuery1 := sq.Select("*").From("musicians_plays_by_month").Where(sq.Eq{
			"musician_id": params.MusicianID,
		}).Suffix("for update").PlaceholderFormat(sq.Dollar)

		if rows, err := lockQuery1.RunWith(r.Conn(ctx)).QueryContext(ctx); err != nil {
			return err
		} else {
			rows.Close()
		}

		lockQuery2 := sq.Select("*").From("musicians_tracks").Where(sq.Eq{
			"track_id":    params.TrackID,
			"musician_id": params.MusicianID,
		}).Suffix("for update").PlaceholderFormat(sq.Dollar)

		if rows, err := lockQuery2.RunWith(r.Conn(ctx)).QueryContext(ctx); err != nil {
			return err
		} else {
			rows.Close()
		}

		query := sq.Update("musicians_plays_by_month").
			Set("plays_total", sq.Expr("plays_total + 1")).
			Where(sq.Eq{
				"musician_id": params.MusicianID,
				"month":       params.Month,
				"year":        params.Year,
			}).PlaceholderFormat(sq.Dollar)

		if _, err := query.RunWith(r.Conn(ctx)).ExecContext(ctx); err != nil {
			return err
		}

		query = sq.Update("musicians_tracks").
			Set("plays_total", sq.Expr("played_times + 1")).
			Where(sq.Eq{
				"id":          params.TrackID,
				"musician_id": params.MusicianID,
			}).PlaceholderFormat(sq.Dollar)

		if _, err := query.RunWith(r.Conn(ctx)).ExecContext(ctx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

type PlaysByDateParams struct {
	MusicianID uuid.UUID
	Month      uint8
	Year       uint32
}

func (r *Repository) PlaysByDate(ctx context.Context, params PlaysByDateParams) (int64, error) {
	var p int64

	if err := r.Transaction(ctx, func(ctx context.Context) error {
		q := sq.Select("plays_total").From("musicians_plays_by_month").
			Where(sq.Eq{
				"musician_id": params.MusicianID,
				"month":       params.Month,
				"year":        params.Year,
			}).
			PlaceholderFormat(sq.Dollar)

		rows, err := q.RunWith(r.Conn(ctx)).QueryContext(ctx)
		if err != nil {
			return err
		}

		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(&p); err != nil {
				return err
			}

			break
		}

		return nil
	}); err != nil {
		return -1, err
	}

	return p, nil
}

type DBTX interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type txCtxKey struct{}

var TxCtxKey = txCtxKey{}

func (r *Repository) Transaction(ctx context.Context, fn func(context.Context) error) (err error) {
	var tx *sql.Tx = new(sql.Tx)

	hasExternalTx := hasExternalTransaction(ctx)

	defer func() {
		if hasExternalTx {
			if err != nil {
				err = fmt.Errorf("error perform operation. %w", err)
				return
			}

			return
		}

		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = errors.Join(fmt.Errorf("error rollback transaction. %w", rbErr), err)
				return
			}

			err = fmt.Errorf("error execute transactional operation. %w", err)

			return
		}

		if commitErr := tx.Commit(); commitErr != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = errors.Join(fmt.Errorf("error rollback transaction. %w", rbErr), commitErr, err)

				return
			}

			err = fmt.Errorf("error commit transaction. %w", err)
		}
	}()

	if !hasExternalTx {
		tx, err = r.db.BeginTx(ctx, &sql.TxOptions{
			Isolation: sql.LevelRepeatableRead,
		})
		if err != nil {
			return fmt.Errorf("error begin transaction. %w", err)
		}

		ctx = context.WithValue(ctx, TxCtxKey, tx)
	}

	return fn(ctx)
}

func (s *Repository) Conn(ctx context.Context) DBTX {
	if tx, ok := ctx.Value(TxCtxKey).(*sql.Tx); ok {
		return tx
	}

	return s.db
}

func hasExternalTransaction(ctx context.Context) bool {
	if _, ok := ctx.Value(TxCtxKey).(*sql.Tx); ok {
		return true
	}

	return false
}
