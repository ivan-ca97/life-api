package types

type PaginationParams struct {
	Limit  int
	Offset int
}

type Page[T any] struct {
	Items  []T
	Total  int64
	Limit  int
	Offset int
}
