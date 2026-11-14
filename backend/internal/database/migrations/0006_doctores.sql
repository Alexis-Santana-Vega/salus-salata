-- 0006 · Catálogo de doctores (paso 1.3) y bloque de dirección de México.
-- Los dominios D_ESTADO y D_CP los reusa el hospital (paso 1.4).

-- Las 32 entidades federativas por su abreviatura (GTO, SLP…, como el
-- prefijo de la clave de hospital en el Excel).
CREATE DOMAIN D_ESTADO AS VARCHAR(5)
  CHECK (VALUE IS NULL OR VALUE IN (
    'AGS', 'BC', 'BCS', 'CAMP', 'COAH', 'COL', 'CHIS', 'CHIH', 'CDMX', 'DGO',
    'GTO', 'GRO', 'HGO', 'JAL', 'MEX', 'MICH', 'MOR', 'NAY', 'NL', 'OAX',
    'PUE', 'QRO', 'QROO', 'SLP', 'SIN', 'SON', 'TAB', 'TAMPS', 'TLAX', 'VER',
    'YUC', 'ZAC'));

-- Código postal: exactamente 5 dígitos.
CREATE DOMAIN D_CP AS CHAR(5)
  CHECK (VALUE IS NULL OR VALUE SIMILAR TO '[0-9]{5}');

-- NUMERO es opcional y, si existe, único para siempre (como una clave).
-- CORREO se guarda en minúsculas y es único solo entre los activos.
-- MUNICIPIO es la "ciudad de residencia" del Excel.
CREATE TABLE DOCTOR (
  ID               BINARY(16)   NOT NULL,
  NUMERO           VARCHAR(20),
  NOMBRE           VARCHAR(120) CHARACTER SET UTF8 NOT NULL COLLATE UNICODE_CI_AI,
  CORREO           VARCHAR(254) NOT NULL,
  TEL_PERSONAL     VARCHAR(20),
  TEL_CONSULTORIO  VARCHAR(20),
  CALLE            VARCHAR(120),
  NUM_EXTERIOR     VARCHAR(20),
  NUM_INTERIOR     VARCHAR(20),
  COLONIA          VARCHAR(120),
  CP               D_CP,
  MUNICIPIO        VARCHAR(120) CHARACTER SET UTF8 COLLATE UNICODE_CI_AI,
  ESTADO           D_ESTADO,
  ACTIVO           BOOLEAN      DEFAULT TRUE NOT NULL,
  CREADO_EN        TIMESTAMP WITH TIME ZONE NOT NULL,
  ACTUALIZADO_EN   TIMESTAMP WITH TIME ZONE NOT NULL,
  CONSTRAINT PK_DOCTOR PRIMARY KEY (ID),
  CONSTRAINT UQ_DOCTOR_NUMERO UNIQUE (NUMERO)
);

CREATE UNIQUE INDEX UQ_DOCTOR_CORREO ON DOCTOR (CORREO) WHERE ACTIVO = TRUE;
CREATE INDEX IX_DOCTOR_NOMBRE ON DOCTOR (NOMBRE);
