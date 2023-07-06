package handlers

import (
	"Vectory/db"
	"Vectory/gen/api/models"
	"Vectory/gen/api/restapi/operations"
	"Vectory/gen/api/restapi/operations/collection"
	"github.com/go-openapi/runtime/middleware"
)

type CollectionHandler struct {
	db *db.DB
}

func (h *CollectionHandler) initHandlers(api *operations.VectoryAPI) {
	api.CollectionGetCollectionHandler = collection.GetCollectionHandlerFunc(h.getCollection)
	api.CollectionAddCollectionHandler = collection.AddCollectionHandlerFunc(h.addCollection)
	api.CollectionDeleteCollectionHandler = collection.DeleteCollectionHandlerFunc(h.deleteCollection)
}

func (h *CollectionHandler) getCollection(params collection.GetCollectionParams) middleware.Responder {
	return collection.NewGetCollectionOK().WithPayload(&models.Collection{})
}

func (h *CollectionHandler) addCollection(params collection.AddCollectionParams) middleware.Responder {
	return collection.NewAddCollectionCreated().WithPayload(&models.CollectionCreated{})
}

func (h *CollectionHandler) deleteCollection(params collection.DeleteCollectionParams) middleware.Responder {
	return collection.NewDeleteCollectionOK().WithPayload(&models.Collection{})
}
