package http

import "github.com/kucinghitam/ethereum-watcher/usecase"

type (
	Handler struct {
		usecase usecase.Parser
	}
)

func New(
	usecase usecase.Parser,
) *Handler {
	return &Handler{
		usecase: usecase,
	}
}
