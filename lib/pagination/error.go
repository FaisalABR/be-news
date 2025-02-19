package pagination

import "errors"

var (
	ErrorMaxPage     = errors.New("choosen page morethan  max page")
	ErrorPage        = errors.New("page must be greather than 0")
	ErrorPageEmpty   = errors.New("page cannot be empty")
	ErrorPageInvalid = errors.New("page invalid, must be an integer")
)
