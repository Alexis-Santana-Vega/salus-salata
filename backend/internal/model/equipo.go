package model

import (
	"time"
	_ "time/tzdata" // la zona va dentro del binario
)

// ZonaNegocio: las fechas de negocio (vigencia de tarifas, "hoy") son de México.
var ZonaNegocio = func() *time.Location {
	loc, err := time.LoadLocation("America/Mexico_City")
	if err != nil {
		panic(err)
	}
	return loc
}()

// Hoy es la fecha de negocio actual (sin hora).
func Hoy() time.Time {
	y, m, d := time.Now().In(ZonaNegocio).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

// Equipo es el modelo maestro (EQ-001) con sus unidades físicas (EQ-001-01).
type Equipo struct {
	ID            ID        `json:"id"`
	Clave         string    `json:"clave"`
	Nombre        string    `json:"nombre"`
	MarcaModelo   string    `json:"marcaModelo"`
	CicloMax      int       `json:"cicloMax"`
	RequierePurga bool      `json:"requierePurga"`
	PurgaCada     *int      `json:"purgaCada"`
	CostoHora     float64   `json:"costoHora"` // tarifa vigente hoy
	Tarifas       []Tarifa  `json:"tarifas"`   // historial, la más reciente primero (solo en el detalle)
	NumUnidades   int       `json:"numUnidades"`
	Unidades      []Unidad  `json:"unidades"` // activas, por clave (solo en el detalle)
	CreadoEn      time.Time `json:"creadoEn"`
	ActualizadoEn time.Time `json:"actualizadoEn"`
}

// Unidad es una pieza física del equipo. Estado y hospital actual: Fase 2.
type Unidad struct {
	ID          ID     `json:"id"`
	Clave       string `json:"clave"`
	NumSerie    string `json:"numSerie"`
	FechaAlta   string `json:"fechaAlta"` // AAAA-MM-DD
	UsosPrevios int    `json:"usosPrevios"`
}

// Tarifa es el costo por hora vigente desde una fecha.
type Tarifa struct {
	CostoHora    float64 `json:"costoHora"`
	VigenteDesde string  `json:"vigenteDesde"` // AAAA-MM-DD
}
