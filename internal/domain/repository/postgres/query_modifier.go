package postgres

type QueryModifier interface {
	Apply(query interface{}) interface{}
}

type PaginationOptions struct {
	Limit  int
	Offset int
}

type SortingOptions struct {
	Column string
	Asc    bool
}

type SearchOptions struct {
	Query   string
	Columns []string
}

type QueryOptions struct {
	Pagination *PaginationOptions
	Sorting    *SortingOptions
	Search     *SearchOptions
}

func NewQueryOptions() *QueryOptions {
	return &QueryOptions{}
}

func (q *QueryOptions) WithPagination(limit, offset int) *QueryOptions {
	q.Pagination = &PaginationOptions{
		Limit:  limit,
		Offset: offset,
	}
	return q
}

func (q *QueryOptions) WithSorting(column string, asc bool) *QueryOptions {
	q.Sorting = &SortingOptions{
		Column: column,
		Asc:    asc,
	}
	return q
}

func (q *QueryOptions) HasPagination() bool {
	return q.Pagination != nil
}

func (q *QueryOptions) HasSorting() bool {
	return q.Sorting != nil
}

func (q *QueryOptions) WithSearch(query string, columns []string) *QueryOptions {
	if query != "" && len(columns) > 0 {
		q.Search = &SearchOptions{
			Query:   query,
			Columns: columns,
		}
	}
	return q
}

func (q *QueryOptions) HasSearch() bool {
	return q.Search != nil && q.Search.Query != ""
}

func ValidateSortColumn(column string, allowed []string, defaultColumn string) string {
	for _, a := range allowed {
		if column == a {
			return column
		}
	}
	return defaultColumn
}
