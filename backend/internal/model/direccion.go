package model

// Estados de la república por su abreviatura (dominio D_ESTADO de la base).
var Estados = []string{
	"AGS", "BC", "BCS", "CAMP", "COAH", "COL", "CHIS", "CHIH", "CDMX", "DGO",
	"GTO", "GRO", "HGO", "JAL", "MEX", "MICH", "MOR", "NAY", "NL", "OAX",
	"PUE", "QRO", "QROO", "SLP", "SIN", "SON", "TAB", "TAMPS", "TLAX", "VER",
	"YUC", "ZAC",
}

// Direccion de México, desglosada. La usan doctor y hospital. Todos los
// campos son opcionales; vacío = no capturado.
type Direccion struct {
	Calle       string `json:"calle"`
	NumExterior string `json:"numExterior"`
	NumInterior string `json:"numInterior"`
	Colonia     string `json:"colonia"`
	CP          string `json:"cp"`
	Municipio   string `json:"municipio"`
	Estado      string `json:"estado"`
}
