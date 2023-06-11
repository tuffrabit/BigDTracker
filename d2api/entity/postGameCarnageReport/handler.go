package postgamecarnagereport

import (
	"fmt"
	"log"

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
	return &PostGameCarnageReport{}
}

func (handler *Handler) DoGet(instanceId string) (*PostGameCarnageReport, error) {
	log.Printf("Getting PostGameCarnageReport data from Bungie for: %v\n", instanceId)

	entity, err := http.DoGet(
		fmt.Sprintf("Stats/PostGameCarnageReport/%v/", instanceId),
		handler,
	)
	if err != nil {
		return nil, fmt.Errorf("Api.GetPostGameCarnageReport: could not create bungie api request: %w", err)
	}

	return entity.(*PostGameCarnageReport), nil
}
