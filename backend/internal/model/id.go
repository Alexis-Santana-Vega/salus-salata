package model

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
)

// ID es la llave de todas las tablas: UUID v7 guardado como BINARY(16).
// En JSON y en la URL viaja como texto ("0193f1c2-…").
// El valor cero (ID{}) se guarda como NULL: sirve para relaciones opcionales.
type ID uuid.UUID

// NuevoID genera un UUID v7 (empieza con la hora: los índices crecen en orden).
func NuevoID() ID { return ID(uuid.Must(uuid.NewV7())) }

// ParseID lee un id en texto; falla si no es un UUID válido.
func ParseID(s string) (ID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return ID{}, err
	}
	return ID(u), nil
}

func (id ID) IsZero() bool   { return id == ID{} }
func (id ID) String() string { return uuid.UUID(id).String() }

// Value manda los 16 bytes al driver (NULL si es cero).
func (id ID) Value() (driver.Value, error) {
	if id.IsZero() {
		return nil, nil
	}
	return id[:], nil
}

// Scan lee los 16 bytes de Firebird (NULL queda en cero).
func (id *ID) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*id = ID{}
		return nil
	case []byte:
		if len(v) != 16 {
			return fmt.Errorf("id: se esperaban 16 bytes, llegaron %d", len(v))
		}
		copy(id[:], v)
		return nil
	}
	return fmt.Errorf("id: tipo no soportado %T", src)
}

func (id ID) MarshalText() ([]byte, error) { return []byte(id.String()), nil }

func (id *ID) UnmarshalText(b []byte) error {
	v, err := ParseID(string(b))
	if err != nil {
		return err
	}
	*id = v
	return nil
}
