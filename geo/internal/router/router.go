package router

import (
	"geo/internal/controller"
	"geo/internal/middlw"
	"github.com/go-chi/chi"
	"gitlab.com/ptflp/gopubsub/queue"
)

func StartRout(cont *controller.HandleGeo, queM queue.MessageQueuer) *chi.Mux {
	r := chi.NewRouter()

	lim := middlw.NewLimit(5, queM)
	r.Use(lim.LimitMiddleware)
	r.Post("/api/address/search", cont.SearchHandle)
	r.Post("/api/address/search", cont.GeocodeHandle)

	return r
}
