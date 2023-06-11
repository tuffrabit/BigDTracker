package entity

type Entity interface {
	SetRawJson(rawJson string)
	GetRawJson() string
}
