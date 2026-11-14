package model

import "time"

// Unidades de insumo (lista cerrada, igual al CHECK de la tabla INSUMO).
const (
	UnidadPieza   = "PIEZA"
	UnidadCaja    = "CAJA"
	UnidadPaquete = "PAQUETE"
)

// Insumo es un consumible desechable que se controla por stock.
type Insumo struct {
	ID            ID        `json:"id"`
	Clave         string    `json:"clave"`
	Descripcion   string    `json:"descripcion"`
	Unidad        string    `json:"unidad"`
	Precio        float64   `json:"precio"`
	CodigoBarras  string    `json:"codigoBarras"`
	StockActual   int       `json:"stockActual"`
	StockMinimo   int       `json:"stockMinimo"`
	BajoMinimo    bool      `json:"bajoMinimo"` // stock actual menor al mínimo
	CreadoEn      time.Time `json:"creadoEn"`
	ActualizadoEn time.Time `json:"actualizadoEn"`
}
