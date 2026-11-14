package model

// Direccion de México. Se captura calle, números y el asentamiento (colonia
// del catálogo SEPOMEX); C.P., colonia, municipio y estado salen del catálogo
// y en la respuesta solo son de lectura. La usan doctor y hospital.
type Direccion struct {
	Calle          string `json:"calle"`
	NumExterior    string `json:"numExterior"`
	NumInterior    string `json:"numInterior"`
	AsentamientoID *ID    `json:"asentamientoId"` // nil = sin colonia capturada
	CP             string `json:"cp"`
	Colonia        string `json:"colonia"`
	Municipio      string `json:"municipio"`
	Estado         string `json:"estado"` // abreviatura (HGO)
	EstadoNombre   string `json:"estadoNombre"`
}

// Asentamiento es una colonia, barrio, fraccionamiento… con su código postal.
type Asentamiento struct {
	ID           ID     `json:"id"`
	CP           string `json:"cp"`
	Nombre       string `json:"nombre"`
	Tipo         string `json:"tipo"`
	Zona         string `json:"zona"`
	Ciudad       string `json:"ciudad"`
	Municipio    string `json:"municipio"`
	Estado       string `json:"estado"`
	EstadoNombre string `json:"estadoNombre"`
}
