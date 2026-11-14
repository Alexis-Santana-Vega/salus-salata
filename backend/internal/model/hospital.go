package model

import "time"

// Hospital es el destino de la renta. La foto llega en la Fase 9.
type Hospital struct {
	ID            ID         `json:"id"`
	Clave         string     `json:"clave"`
	Nombre        string     `json:"nombre"`
	Telefono      string     `json:"telefono"`
	Direccion     Direccion  `json:"direccion"`
	Latitud       *float64   `json:"latitud"`
	Longitud      *float64   `json:"longitud"`
	LigaMapa      string     `json:"ligaMapa"`
	Contactos     []Contacto `json:"contactos"`
	CreadoEn      time.Time  `json:"creadoEn"`
	ActualizadoEn time.Time  `json:"actualizadoEn"`
}

// Contacto es una persona en el hospital para coordinar entregas y devoluciones.
type Contacto struct {
	Nombre   string `json:"nombre"`
	Puesto   string `json:"puesto"`
	Telefono string `json:"telefono"`
}
