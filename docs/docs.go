// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/ad": {
            "get": {
                "description": "Get a list of ads with queries",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "ad"
                ],
                "summary": "Public API",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Get ads starting from offset",
                        "name": "offset",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Get how many ads, default is 5",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Target age",
                        "name": "age",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Target gender",
                        "name": "gender",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Target country",
                        "name": "country",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Target platform",
                        "name": "platform",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\\\"items\\\": [ad, ...]}",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "array",
                                "items": {
                                    "$ref": "#/definitions/domain.Ad"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/domain.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Create an ad",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "ad"
                ],
                "summary": "Admin API",
                "parameters": [
                    {
                        "description": "Add an ad",
                        "name": "ad",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.Ad"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/domain.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/domain.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.Ad": {
            "type": "object",
            "properties": {
                "condition": {
                    "$ref": "#/definitions/domain.Condition"
                },
                "endAt": {
                    "type": "string"
                },
                "startAt": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "domain.Condition": {
            "type": "object",
            "properties": {
                "ageEnd": {
                    "type": "integer"
                },
                "ageStart": {
                    "type": "integer"
                },
                "country": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "gender": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "platform": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "domain.ErrorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "domain.SuccessResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "127.0.0.1:3000",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Dcard AD API",
	Description:      "The server for AD services",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
