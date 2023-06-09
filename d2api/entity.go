package d2api

type Entity interface {
	UnmarshalHttpResponseBody([]byte) error
}
