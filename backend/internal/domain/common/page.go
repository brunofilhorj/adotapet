package common

type PageRequest struct {
	Page int
	Size int
}

type Page[T any] struct {
	Items         []T
	Page          int
	Size          int
	TotalElements int64
	TotalPages    int
}
