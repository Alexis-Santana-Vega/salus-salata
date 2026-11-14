// Package response define el formato único de respuestas de la API.
//
// Éxito:   {"data": ...}                      (un recurso)
//
//	{"data": [...], "meta": {...}}     (un listado)
//
// Error:   {"error": {"code", "message", "details", "requestId"}}
package response

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
)

// Meta acompaña a los listados paginados.
type Meta struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

type dataBody struct {
	Data any   `json:"data"`
	Meta *Meta `json:"meta,omitempty"`
}

// OK responde 200 con {"data": ...}.
func OK(c fiber.Ctx, data any) error {
	return c.Status(http.StatusOK).JSON(dataBody{Data: data})
}

// Created responde 201 con {"data": ...}.
func Created(c fiber.Ctx, data any) error {
	return c.Status(http.StatusCreated).JSON(dataBody{Data: data})
}

// List responde 200 con {"data": [...], "meta": {...}}.
func List(c fiber.Ctx, data any, meta Meta) error {
	return c.Status(http.StatusOK).JSON(dataBody{Data: data, Meta: &meta})
}

// NoContent responde 204 sin cuerpo.
func NoContent(c fiber.Ctx) error {
	return c.SendStatus(http.StatusNoContent)
}
