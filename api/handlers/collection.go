package handlers

import (
	"Vectory/db"
	collectionent "Vectory/entities/collection"
	"Vectory/gen/api/models"
	"Vectory/gen/api/restapi/operations"
	"Vectory/gen/api/restapi/operations/collection"
	"errors"
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

// getCollection handler for getting collection configuration
func (h *CollectionHandler) getCollection(params collection.GetCollectionParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()

	c, err := h.db.GetCollection(ctx, params.CollectionName)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, db.ErrValidationFailed) {
			code = http.StatusBadRequest
		}

		return middleware.Error(code, handleError(err))
	}

	cfg := c.GetConfig()

	col := models.Collection{
		Name:           cfg.Name,
		IndexType:      cfg.IndexType,
		EmbedderType:   cfg.EmbedderType,
		IndexParams:    cfg.IndexParams,
		EmbedderConfig: cfg.EmbedderConfig,
		DataType:       cfg.DataType,
	}

	return collection.NewGetCollectionOK().WithPayload(&col)
}

// addCollection handler for adding a new collection to Vectory
func (h *CollectionHandler) addCollection(params collection.AddCollectionParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()

	cfg := collectionent.Collection{
		Name:           params.Collection.Name,
		IndexType:      params.Collection.IndexType,
		EmbedderType:   params.Collection.EmbedderType,
		IndexParams:    params.Collection.IndexParams,
		EmbedderConfig: params.Collection.EmbedderConfig,
		DataType:       params.Collection.DataType,
	}

	_, err := h.db.CreateCollection(ctx, &cfg)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, db.ErrValidationFailed) {
			code = http.StatusBadRequest
		}

		return middleware.Error(code, handleError(err))
	}

	return collection.NewAddCollectionCreated().WithPayload(&models.CollectionCreated{CollectionName: cfg.Name})
}

// deleteCollection handler for deleting a collection from Vectory
func (h *CollectionHandler) deleteCollection(params collection.DeleteCollectionParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()

	err := h.db.DeleteCollection(ctx, params.CollectionName)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, db.ErrValidationFailed) {
			code = http.StatusBadRequest
		}

		return middleware.Error(code, handleError(err))
	}

	return collection.NewDeleteCollectionOK().WithPayload(&models.APIResponse{Message: "deleted successfully"})
}
