package entity

import "github.com/tuffrabit/BigDTracker/d2api/meta"

type Handler interface {
	NewEntity() Entity
	GetMeta() *meta.Meta
}
