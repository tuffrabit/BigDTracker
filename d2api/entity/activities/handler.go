package activities

import (
	"fmt"
	"log"
	"strconv"

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
	return &Activities{}
}

func (handler *Handler) DoGet(page int, membershipId string, membershipType int, characterId string) (*Activities, error) {
	log.Printf("Getting Activity data from Bungie for: MID:%v CID:%v P:%v\n", membershipId, characterId, page)

	entity, err := http.DoGet(
		fmt.Sprintf("%v/Account/%v/Character/%v/Stats/Activities/?mode=AllPvP&count=%v&page=%v", strconv.Itoa(membershipType), membershipId, characterId, handler.GetMeta().ActivitiesPageSizeString, strconv.Itoa(page)),
		handler,
	)
	if err != nil {
		return nil, fmt.Errorf("Api.GetActivities: could not create bungie api request: %w", err)
	}

	return entity.(*Activities), nil
}
