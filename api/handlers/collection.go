package handlers

import (
	"Vectory/db"
	"Vectory/entities"
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
		Name:        cfg.Name,
		IndexType:   cfg.IndexType,
		Embedder:    cfg.Embedder,
		DataType:    cfg.DataType,
		IndexParams: cfg.IndexParams,
	}

	return collection.NewGetCollectionOK().WithPayload(&col)
}

func (h *CollectionHandler) addCollection(params collection.AddCollectionParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()

	cfg := entities.Collection{
		Name:        params.Collection.Name,
		IndexType:   params.Collection.IndexType,
		Embedder:    params.Collection.Embedder,
		DataType:    params.Collection.DataType,
		IndexParams: params.Collection.IndexParams,
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
