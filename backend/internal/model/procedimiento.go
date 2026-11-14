package model

import "time"

// Procedimiento es el tipo de cirugía (Holep, Artroscopia…). Su plantilla de
// material llega en la Fase 9.
type Procedimiento struct {
	ID            ID        `json:"id"`
	Nombre        string    `json:"nombre"`
	CreadoEn      time.Time `json:"creadoEn"`
	ActualizadoEn time.Time `json:"actualizadoEn"`
}
