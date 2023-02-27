package http

import "github.com/kucinghitam/ethereum-watcher/usecase"

type (
	Handler struct {
		usecase usecase.Watcher
	}
)

func New(
	usecase usecase.Watcher,
) *Handler {
	return &Handler{
		usecase: usecase,
	}
}
