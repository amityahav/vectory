package handlers

import (
	"Vectory/db"
	"Vectory/db/core/indexes"
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

	err := validateCollection(params.Collection)
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

// TODO: validate
func validateCollection(cfg *models.Collection) error {
	if cfg.Name == "" {
		return ErrCollectionNameEmpty
	}

	if _, ok := indexes.SupportedIndexes[cfg.IndexType]; !ok {
		return ErrIndexTypeUnsupported
	}

	if _, ok := db.SupportedDataTypes[cfg.DataType]; !ok {
		return ErrDataTypeUnsupported
	}

	return validateIndexParams(cfg.IndexType, cfg.IndexParams)
}

// TODO: validate
func validateIndexParams(indexType string, params interface{}) error {
	switch indexType {
	case "disk_ann":

	case "hnsw":

	}

	return nil
}
