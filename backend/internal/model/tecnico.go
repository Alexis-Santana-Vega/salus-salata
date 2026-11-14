package model

import "time"

// Puestos del técnico (lista cerrada, igual al CHECK de la tabla TECNICO).
const (
	PuestoAyudante       = "AYUDANTE"
	PuestoEnfermeria     = "ENFERMERIA"
	PuestoAnestesiologia = "ANESTESIOLOGIA"
	PuestoLogistica      = "LOGISTICA"
	PuestoOtro           = "OTRO"
)

// Tecnico es el personal de apoyo que entrega el material y asiste en la
// cirugía (los "ayudantes" del Excel).
type Tecnico struct {
	ID            ID        `json:"id"`
	Nombre        string    `json:"nombre"`
	Puesto        string    `json:"puesto"`
	Telefono      string    `json:"telefono"`
	CreadoEn      time.Time `json:"creadoEn"`
	ActualizadoEn time.Time `json:"actualizadoEn"`
}
