package handlers

import (
	"Vectory/db"
	"Vectory/gen/api/restapi/operations"
)

// InitHandlers registers api handlers for the api
func InitHandlers(api *operations.VectoryAPI, db *db.DB) {
	collectionHandler := CollectionHandler{db: db}
	collectionHandler.initHandlers(api)
}
