package d2api

import (
	"strconv"

	"github.com/tuffrabit/BigDTracker/d2api/entity/activities"
	"github.com/tuffrabit/BigDTracker/d2api/entity/player"
	postgamecarnagereport "github.com/tuffrabit/BigDTracker/d2api/entity/postGameCarnageReport"
	"github.com/tuffrabit/BigDTracker/d2api/entity/profile"
	"github.com/tuffrabit/BigDTracker/d2api/meta"
)

const baseUrl string = "https://www.bungie.net/Platform/Destiny2/"
const ActivitiesPageSize int = 100
const BigDApiHash int = 3824106094

type Api struct {
	PlayerHandler                *player.Handler
	ProfileHandler               *profile.Handler
	ActivitiesHandler            *activities.Handler
	PostGameCarnageReportHandler *postgamecarnagereport.Handler
}

func NewApi(apiKey string) *Api {
	meta := NewMeta(apiKey)

	return &Api{
		PlayerHandler:                NewPlayerHandler(meta),
		ProfileHandler:               NewProfileHandler(meta),
		ActivitiesHandler:            NewActivitiesHandler(meta),
		PostGameCarnageReportHandler: NewPostGameCarnageReportHandler(meta),
	}
}

func NewMeta(apiKey string) *meta.Meta {
	return &meta.Meta{
		ApiKey:                   apiKey,
		BaseUrl:                  baseUrl,
		ActivitiesPageSizeString: strconv.Itoa(ActivitiesPageSize),
	}
}

func NewPlayerHandler(meta *meta.Meta) *player.Handler {
	return &player.Handler{
		Meta: meta,
	}
}

func NewProfileHandler(meta *meta.Meta) *profile.Handler {
	return &profile.Handler{
		Meta: meta,
	}
}

func NewActivitiesHandler(meta *meta.Meta) *activities.Handler {
	return &activities.Handler{
		Meta: meta,
	}
}

func NewPostGameCarnageReportHandler(meta *meta.Meta) *postgamecarnagereport.Handler {
	return &postgamecarnagereport.Handler{
		Meta: meta,
	}
}
