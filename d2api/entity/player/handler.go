package player

import (
	"fmt"
	"log"
	"net/url"

	"github.com/tuffrabit/BigDTracker/d2api/entity"
	"github.com/tuffrabit/BigDTracker/d2api/http"
	"github.com/tuffrabit/BigDTracker/d2api/meta"
)

type Handler struct {
	Meta *meta.Meta
}

func (handler *Handler) GetMeta() *meta.Meta {
	return handler.Meta
}

func (handler *Handler) NewEntity() entity.Entity {
	return &Player{}
}

func (handler *Handler) DoGet(user string) (*Player, error) {
	log.Printf("Getting Player data from Bungie for: %v\n", user)

	entity, err := http.DoGet(
		fmt.Sprintf("SearchDestinyPlayer/-1/%v/", url.QueryEscape(user)),
		handler,
	)
	if err != nil {
		return nil, fmt.Errorf("d2api/entity/player.DoGet: could not get player data from bungie api: %w", err)
	}

	return entity.(*Player), nil
}
