package handlers

import (
	"Vectory/api/validators"
	"Vectory/db"
	"Vectory/gen/api/models"
	"Vectory/gen/api/restapi/operations"
	"Vectory/gen/api/restapi/operations/collection"
	"github.com/go-openapi/runtime/middleware"
	"net/http"
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
	ctx := params.HTTPRequest.Context()

	err := validators.ValidateCollection(params.Collection)
	if err != nil {
		return middleware.Error(http.StatusBadRequest, err)
	}

	id, err := h.db.CreateCollection(ctx, params.Collection)
	if err != nil {
		return middleware.Error(http.StatusInternalServerError, err)
	}

	return collection.NewAddCollectionCreated().WithPayload(&models.CollectionCreated{CollectionID: int64(id)})
}

func (h *CollectionHandler) deleteCollection(params collection.DeleteCollectionParams) middleware.Responder {
	return collection.NewDeleteCollectionOK().WithPayload(&models.Collection{})
}
