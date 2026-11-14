-- 0004 · Catálogo de procedimientos (paso 1.1).
-- Patrón de todos los catálogos: ID UUID v7, ACTIVO para la baja lógica y
-- fechas de creación y cambio. NOMBRE con colación UNICODE_CI_AI: "holep",
-- "Holep" y "HOLEP" son el mismo nombre al validar duplicados y al buscar.

CREATE TABLE PROCEDIMIENTO (
  ID              BINARY(16)   NOT NULL,
  NOMBRE          VARCHAR(120) CHARACTER SET UTF8 NOT NULL COLLATE UNICODE_CI_AI,
  ACTIVO          BOOLEAN      DEFAULT TRUE NOT NULL,
  CREADO_EN       TIMESTAMP WITH TIME ZONE NOT NULL,
  ACTUALIZADO_EN  TIMESTAMP WITH TIME ZONE NOT NULL,
  CONSTRAINT PK_PROCEDIMIENTO PRIMARY KEY (ID)
);

-- El nombre no se repite entre los activos; uno dado de baja libera el nombre.
CREATE UNIQUE INDEX UQ_PROCEDIMIENTO_NOMBRE ON PROCEDIMIENTO (NOMBRE) WHERE ACTIVO = TRUE;
