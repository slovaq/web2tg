package API

import "github.com/go-chi/chi"

func Router(r chi.Router) {
	r.Get("/user_create", UserCreate)
	r.Get("/user_get", UserGet)

	r.Get("/record_get", RecordGet)
	r.Get("/record_create", RecordCreate)

	r.Get("/city_create", CityCreate)
	r.Get("/city_get", CityGet)
}
