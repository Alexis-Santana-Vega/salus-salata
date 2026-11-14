package model

import "time"

// Doctor es el médico encargado de las cirugías. Su acceso al portal
// (usuario con rol Doctor) llega en la Fase 6.
type Doctor struct {
	ID             ID        `json:"id"`
	Numero         string    `json:"numero"`
	Nombre         string    `json:"nombre"`
	Correo         string    `json:"correo"`
	TelPersonal    string    `json:"telPersonal"`
	TelConsultorio string    `json:"telConsultorio"`
	Direccion      Direccion `json:"direccion"`
	CreadoEn       time.Time `json:"creadoEn"`
	ActualizadoEn  time.Time `json:"actualizadoEn"`
}
