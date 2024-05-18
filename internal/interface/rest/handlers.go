package rest

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/emo2007/block-accounting/examples/license-api/internal/usecases/repository"
	"github.com/google/uuid"
)

type GetMusiciansRequest struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	FromId string `json:"from_id"`
}

type Musician struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (s *Server) handleGetMusicians(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	req, err := CreateRequest[GetMusiciansRequest](r)
	if err != nil {
		return nil, err
	}

	id := uuid.Nil
	if req.ID != "" {
		id = uuid.MustParse(req.ID)
	}

	fromId := uuid.Nil

	if req.FromId != "" {
		fromId = uuid.MustParse(req.FromId)
	}

	mus, err := s.repo.ListMusicians(r.Context(), repository.ListMusiciansParams{
		Ids:    uuid.UUIDs{id},
		Name:   req.Name,
		FromId: fromId,
	})
	if err != nil {
		return nil, err
	}

	muslist := make([]Musician, len(mus))

	for i, m := range mus {
		muslist[i] = Musician{
			ID:   m.ID.String(),
			Name: m.Name,
		}
	}

	data, err := json.Marshal(muslist)
	if err != nil {
		return nil, err
	}

	w.Header().Add("Content-Type", "application/json")

	return data, nil
}

func (s *Server) handleGetPlayout(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	w.Header().Add("Content-Type", "text/plain")

	// todo unmock

	i := rand.Int31n(100)

	return []byte(fmt.Sprintf("%d", i)), nil
}
