package profile

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
	return &Profile{}
}

func (handler *Handler) DoGet(membershipId string, membershipType int) (*Profile, error) {
	log.Printf("Getting Profile data from Bungie for: %v\n", membershipId)

	entity, err := http.DoGet(
		fmt.Sprintf("%v/Profile/%v/?components=100", strconv.Itoa(membershipType), membershipId),
		handler,
	)
	if err != nil {
		return nil, fmt.Errorf("d2api/entity/profile.DoGet: could not get profile data from bungie api: %w", err)
	}

	return entity.(*Profile), nil
}
