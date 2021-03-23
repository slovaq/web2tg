package API

import (
	"net/http"

	"github.com/go-chi/chi"
)

// Router...
func Router(r chi.Router) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("access denied"))
	})
	r.Get("/user_create", UserCreate)
	r.Get("/user_get", UserGet)

	r.Get("/record_get", RecordGet)
	r.Get("/record_create", RecordCreate)

	r.Get("/city_create", CityCreate)
	r.Get("/city_get", CityGet)
	r.Get("/city_getAll", CityGetList)
}
