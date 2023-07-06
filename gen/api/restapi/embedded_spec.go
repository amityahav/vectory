// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "swagger": "2.0",
  "info": {
    "title": "Vectory",
    "version": "1"
  },
  "paths": {
    "/v1/collection": {
      "post": {
        "description": "Add a new collection to the database",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "collection"
        ],
        "summary": "Add a collection to the database",
        "operationId": "addCollection",
        "parameters": [
          {
            "name": "collection",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/Collection"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Created successfully",
            "schema": {
              "$ref": "#/definitions/CollectionCreated"
            }
          }
        }
      }
    },
    "/v1/collection/{collectionName}": {
      "get": {
        "description": "Get collection information",
        "tags": [
          "collection"
        ],
        "summary": "Get collection information",
        "operationId": "getCollection",
        "parameters": [
          {
            "type": "string",
            "description": "Collection name to delete",
            "name": "collectionName",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "valid operation",
            "schema": {
              "$ref": "#/definitions/Collection"
            }
          },
          "400": {
            "description": "Invalid collection name"
          }
        }
      },
      "delete": {
        "description": "Delete a collection from the database",
        "tags": [
          "collection"
        ],
        "summary": "Delete a collection from the database",
        "operationId": "deleteCollection",
        "parameters": [
          {
            "type": "string",
            "description": "Collection name to delete",
            "name": "collectionName",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "valid operation",
            "schema": {
              "$ref": "#/definitions/Collection"
            }
          },
          "400": {
            "description": "Invalid collection name"
          }
        }
      }
    }
  },
  "definitions": {
    "ApiResponse": {
      "type": "object",
      "properties": {
        "code": {
          "type": "integer",
          "format": "int32",
          "x-order": 0
        },
        "message": {
          "type": "string",
          "x-order": 2
        },
        "type": {
          "type": "string",
          "x-order": 1
        }
      },
      "xml": {
        "name": "##default"
      }
    },
    "Collection": {
      "type": "object",
      "required": [
        "name"
      ],
      "properties": {
        "data_type": {
          "type": "string",
          "x-order": 4,
          "example": "text"
        },
        "distance_metric": {
          "type": "string",
          "x-order": 2,
          "example": "dot"
        },
        "embedder": {
          "type": "string",
          "x-order": 3,
          "example": "text2vec"
        },
        "index_params": {
          "type": "object",
          "x-order": 5
        },
        "index_type": {
          "type": "string",
          "x-order": 1,
          "example": "disk_ann"
        },
        "name": {
          "type": "string",
          "x-order": 0,
          "example": "movie-reviews"
        }
      }
    },
    "CollectionCreated": {
      "type": "object",
      "properties": {
        "collection_id": {
          "type": "string",
          "x-order": 0
        }
      }
    }
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "swagger": "2.0",
  "info": {
    "title": "Vectory",
    "version": "1"
  },
  "paths": {
    "/v1/collection": {
      "post": {
        "description": "Add a new collection to the database",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "collection"
        ],
        "summary": "Add a collection to the database",
        "operationId": "addCollection",
        "parameters": [
          {
            "name": "collection",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/Collection"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Created successfully",
            "schema": {
              "$ref": "#/definitions/CollectionCreated"
            }
          }
        }
      }
    },
    "/v1/collection/{collectionName}": {
      "get": {
        "description": "Get collection information",
        "tags": [
          "collection"
        ],
        "summary": "Get collection information",
        "operationId": "getCollection",
        "parameters": [
          {
            "type": "string",
            "description": "Collection name to delete",
            "name": "collectionName",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "valid operation",
            "schema": {
              "$ref": "#/definitions/Collection"
            }
          },
          "400": {
            "description": "Invalid collection name"
          }
        }
      },
      "delete": {
        "description": "Delete a collection from the database",
        "tags": [
          "collection"
        ],
        "summary": "Delete a collection from the database",
        "operationId": "deleteCollection",
        "parameters": [
          {
            "type": "string",
            "description": "Collection name to delete",
            "name": "collectionName",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "valid operation",
            "schema": {
              "$ref": "#/definitions/Collection"
            }
          },
          "400": {
            "description": "Invalid collection name"
          }
        }
      }
    }
  },
  "definitions": {
    "ApiResponse": {
      "type": "object",
      "properties": {
        "code": {
          "type": "integer",
          "format": "int32",
          "x-order": 0
        },
        "message": {
          "type": "string",
          "x-order": 2
        },
        "type": {
          "type": "string",
          "x-order": 1
        }
      },
      "xml": {
        "name": "##default"
      }
    },
    "Collection": {
      "type": "object",
      "required": [
        "name"
      ],
      "properties": {
        "data_type": {
          "type": "string",
          "x-order": 4,
          "example": "text"
        },
        "distance_metric": {
          "type": "string",
          "x-order": 2,
          "example": "dot"
        },
        "embedder": {
          "type": "string",
          "x-order": 3,
          "example": "text2vec"
        },
        "index_params": {
          "type": "object",
          "x-order": 5
        },
        "index_type": {
          "type": "string",
          "x-order": 1,
          "example": "disk_ann"
        },
        "name": {
          "type": "string",
          "x-order": 0,
          "example": "movie-reviews"
        }
      }
    },
    "CollectionCreated": {
      "type": "object",
      "properties": {
        "collection_id": {
          "type": "string",
          "x-order": 0
        }
      }
    }
  }
}`))
}