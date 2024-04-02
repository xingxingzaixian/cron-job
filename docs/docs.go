// Code generated by swaggo/swag. DO NOT EDIT
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
        "/api/task/edit": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "新增/修改任务",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务管理"
                ],
                "summary": "新增/修改任务",
                "operationId": "/api/task/edit",
                "parameters": [
                    {
                        "description": "body",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/schemas.TaskEditHTTPInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/schemas.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/task/list": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "任务列表",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务管理"
                ],
                "summary": "任务列表",
                "operationId": "/api/task/list",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "页数",
                        "name": "pageNo",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "每页条数",
                        "name": "pageSize",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "任务名称",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "任务协议",
                        "name": "protocol",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "任务标签",
                        "name": "tag",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/schemas.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/schemas.SearchTaskResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/task/op": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "操作任务",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务管理"
                ],
                "summary": "操作任务",
                "operationId": "/api/task/op",
                "parameters": [
                    {
                        "description": "操作信息",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/schemas.TaskOptionInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/schemas.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "boolean"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/taskLog/list": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "任务执行日志列表",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务日志管理"
                ],
                "summary": "任务执行日志列表",
                "operationId": "/api/taskLog/list",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "页数",
                        "name": "pageNo",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "每页条数",
                        "name": "pageSize",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "任务ID",
                        "name": "taskId",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/schemas.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/schemas.TaskLogOutput"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/taskLog/query": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "任务执行日志",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "任务日志管理"
                ],
                "summary": "任务执行日志",
                "operationId": "/api/taskLog/query",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "任务日志ID",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/schemas.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/schemas.TaskLogItemOutput"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "global.TaskPolicy": {
            "type": "integer",
            "enum": [
                1,
                2,
                3,
                4
            ],
            "x-enum-comments": {
                "TaskPolicyMulti": "并行策略",
                "TaskPolicyOnce": "单词策略",
                "TaskPolicySingle": "单利策略",
                "TaskPolicyTimes": "多次策略"
            },
            "x-enum-varnames": [
                "TaskPolicyMulti",
                "TaskPolicyOnce",
                "TaskPolicySingle",
                "TaskPolicyTimes"
            ]
        },
        "global.TaskProtocol": {
            "type": "integer",
            "enum": [
                1,
                2,
                3
            ],
            "x-enum-varnames": [
                "TaskProtocolHttp",
                "TaskProtocolShell",
                "TaskProtocolGrpc"
            ]
        },
        "global.TaskStatus": {
            "type": "integer",
            "enum": [
                0,
                1,
                2,
                3,
                4,
                5,
                6
            ],
            "x-enum-comments": {
                "TaskStatusCancel": "取消",
                "TaskStatusDisabled": "禁用",
                "TaskStatusEnabled": "启用",
                "TaskStatusFailure": "失败",
                "TaskStatusFinish": "成功",
                "TaskStatusRunning": "运行中",
                "TaskStatusTimeout": "超时"
            },
            "x-enum-varnames": [
                "TaskStatusDisabled",
                "TaskStatusEnabled",
                "TaskStatusFailure",
                "TaskStatusRunning",
                "TaskStatusFinish",
                "TaskStatusCancel",
                "TaskStatusTimeout"
            ]
        },
        "schemas.Response": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {},
                "message": {
                    "type": "string"
                }
            }
        },
        "schemas.SearchTaskResponse": {
            "type": "object",
            "properties": {
                "list": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/schemas.TaskItemOutput"
                    }
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "schemas.TaskEditHTTPInput": {
            "type": "object",
            "required": [
                "command",
                "name",
                "protocol",
                "spec"
            ],
            "properties": {
                "command": {
                    "type": "string",
                    "example": ""
                },
                "count": {
                    "type": "integer",
                    "default": 0,
                    "example": 0
                },
                "delay": {
                    "type": "integer",
                    "default": 0,
                    "example": 0
                },
                "id": {
                    "type": "integer",
                    "default": 0
                },
                "name": {
                    "type": "string",
                    "example": ""
                },
                "params": {
                    "type": "string",
                    "example": ""
                },
                "policy": {
                    "default": 1,
                    "allOf": [
                        {
                            "$ref": "#/definitions/global.TaskPolicy"
                        }
                    ],
                    "example": 1
                },
                "protocol": {
                    "allOf": [
                        {
                            "$ref": "#/definitions/global.TaskProtocol"
                        }
                    ],
                    "example": 1
                },
                "remark": {
                    "type": "string",
                    "example": ""
                },
                "retry_interval": {
                    "type": "integer",
                    "default": 0,
                    "example": 0
                },
                "retry_times": {
                    "type": "integer",
                    "default": 0,
                    "example": 0
                },
                "spec": {
                    "type": "string",
                    "example": ""
                },
                "status": {
                    "default": 0,
                    "allOf": [
                        {
                            "$ref": "#/definitions/global.TaskStatus"
                        }
                    ]
                },
                "tag": {
                    "type": "string",
                    "example": ""
                },
                "timeout": {
                    "type": "integer",
                    "default": 0,
                    "example": 0
                }
            }
        },
        "schemas.TaskItemOutput": {
            "type": "object",
            "required": [
                "command",
                "name",
                "protocol",
                "spec"
            ],
            "properties": {
                "command": {
                    "type": "string",
                    "example": ""
                },
                "count": {
                    "type": "integer",
                    "default": 0,
                    "example": 0
                },
                "delay": {
                    "type": "integer",
                    "default": 0,
                    "example": 0
                },
                "id": {
                    "type": "integer",
                    "default": 0
                },
                "name": {
                    "type": "string",
                    "example": ""
                },
                "params": {
                    "type": "string",
                    "example": ""
                },
                "policy": {
                    "default": 1,
                    "allOf": [
                        {
                            "$ref": "#/definitions/global.TaskPolicy"
                        }
                    ],
                    "example": 1
                },
                "protocol": {
                    "allOf": [
                        {
                            "$ref": "#/definitions/global.TaskProtocol"
                        }
                    ],
                    "example": 1
                },
                "remark": {
                    "type": "string",
                    "example": ""
                },
                "retry_interval": {
                    "type": "integer",
                    "default": 0,
                    "example": 0
                },
                "retry_times": {
                    "type": "integer",
                    "default": 0,
                    "example": 0
                },
                "spec": {
                    "type": "string",
                    "example": ""
                },
                "status": {
                    "default": 0,
                    "allOf": [
                        {
                            "$ref": "#/definitions/global.TaskStatus"
                        }
                    ]
                },
                "tag": {
                    "type": "string",
                    "example": ""
                },
                "timeout": {
                    "type": "integer",
                    "default": 0,
                    "example": 0
                }
            }
        },
        "schemas.TaskLogItemOutput": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "result": {
                    "description": "执行结果",
                    "type": "string"
                },
                "retry_times": {
                    "description": "任务重试次数",
                    "type": "integer"
                },
                "status": {
                    "description": "状态 0:执行失败 1:执行中  2:执行完毕 3:任务取消(上次任务未执行完成) 4:异步执行",
                    "allOf": [
                        {
                            "$ref": "#/definitions/global.TaskStatus"
                        }
                    ]
                },
                "task_id": {
                    "description": "任务id",
                    "type": "integer"
                },
                "total_time": {
                    "description": "执行总时长",
                    "type": "integer"
                }
            }
        },
        "schemas.TaskLogOutput": {
            "type": "object",
            "properties": {
                "list": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/schemas.TaskLogItemOutput"
                    }
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "schemas.TaskOptionInput": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "op": {
                    "type": "string",
                    "example": "stop"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}