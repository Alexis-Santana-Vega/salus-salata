-- 0005 · Catálogo de técnicos / personal de apoyo (paso 1.2).
-- PUESTO es una lista cerrada (mock: Alta de técnico): no lleva tabla propia.
-- Agregar un puesto = migración nueva que cambie el CHECK y la lista en Go.
-- El nombre no es único: puede haber dos personas con el mismo nombre.

CREATE TABLE TECNICO (
  ID              BINARY(16)   NOT NULL,
  NOMBRE          VARCHAR(120) CHARACTER SET UTF8 NOT NULL COLLATE UNICODE_CI_AI,
  PUESTO          VARCHAR(20)  NOT NULL,
  TELEFONO        VARCHAR(20),
  ACTIVO          BOOLEAN      DEFAULT TRUE NOT NULL,
  CREADO_EN       TIMESTAMP WITH TIME ZONE NOT NULL,
  ACTUALIZADO_EN  TIMESTAMP WITH TIME ZONE NOT NULL,
  CONSTRAINT PK_TECNICO PRIMARY KEY (ID),
  CONSTRAINT CK_TECNICO_PUESTO CHECK (PUESTO IN ('AYUDANTE', 'ENFERMERIA', 'ANESTESIOLOGIA', 'LOGISTICA', 'OTRO'))
);

CREATE INDEX IX_TECNICO_NOMBRE ON TECNICO (NOMBRE);
