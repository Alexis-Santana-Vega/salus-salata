-- 0002 · Roles, usuarios, sesiones y auditoría de accesos (paso 0.4)
-- IDs: UUID v7 generado en Go, guardado en 16 bytes (BINARY(16)).
-- Fechas como TIMESTAMP WITH TIME ZONE: guardan el instante exacto sin
-- depender de la zona horaria del servidor.

-- El código se refiere a cada rol por su CLAVE, nunca por su id.
CREATE TABLE ROL (
  ID      BINARY(16)  NOT NULL,
  CLAVE   VARCHAR(20) NOT NULL,
  NOMBRE  VARCHAR(50) NOT NULL,
  CONSTRAINT PK_ROL PRIMARY KEY (ID),
  CONSTRAINT UQ_ROL_CLAVE UNIQUE (CLAVE)
);

-- USUARIO y CORREO se guardan en minúsculas para que el acceso no distinga
-- mayúsculas. La vinculación con DOCTOR llega con el catálogo de doctores.
CREATE TABLE USUARIO (
  ID                      BINARY(16)   NOT NULL,
  USUARIO                 VARCHAR(50)  NOT NULL,
  CORREO                  VARCHAR(254) NOT NULL,
  NOMBRE                  VARCHAR(120) NOT NULL,
  HASH_CONTRASENA         VARCHAR(255) NOT NULL,
  ROL_ID                  BINARY(16)   NOT NULL,
  ACTIVO                  BOOLEAN      DEFAULT TRUE NOT NULL,
  DEBE_CAMBIAR_CONTRASENA BOOLEAN      DEFAULT FALSE NOT NULL,
  ULTIMO_ACCESO           TIMESTAMP WITH TIME ZONE,
  CREADO_EN               TIMESTAMP WITH TIME ZONE NOT NULL,
  ACTUALIZADO_EN          TIMESTAMP WITH TIME ZONE NOT NULL,
  CONSTRAINT PK_USUARIO PRIMARY KEY (ID),
  CONSTRAINT FK_USUARIO_ROL FOREIGN KEY (ROL_ID) REFERENCES ROL (ID),
  CONSTRAINT UQ_USUARIO_USUARIO UNIQUE (USUARIO),
  CONSTRAINT UQ_USUARIO_CORREO UNIQUE (CORREO)
);

-- Una fila por dispositivo con sesión abierta. Del refresh token solo se
-- guarda su SHA-256: si se filtra la tabla, los tokens no sirven.
-- FAMILIA agrupa las rotaciones de una misma sesión para detectar reuso.
CREATE TABLE SESION (
  ID           BINARY(16)   NOT NULL,
  USUARIO_ID   BINARY(16)   NOT NULL,
  FAMILIA      CHAR(32)     NOT NULL,
  TOKEN_HASH   CHAR(64)     NOT NULL,
  CREADA_EN    TIMESTAMP WITH TIME ZONE NOT NULL,
  EXPIRA_EN    TIMESTAMP WITH TIME ZONE NOT NULL,
  REVOCADA_EN  TIMESTAMP WITH TIME ZONE,
  IP           VARCHAR(45),
  USER_AGENT   VARCHAR(255),
  CONSTRAINT PK_SESION PRIMARY KEY (ID),
  CONSTRAINT FK_SESION_USUARIO FOREIGN KEY (USUARIO_ID) REFERENCES USUARIO (ID),
  CONSTRAINT UQ_SESION_TOKEN UNIQUE (TOKEN_HASH)
);

CREATE INDEX IX_SESION_FAMILIA ON SESION (FAMILIA);

-- Bitácora de accesos. También sirve para contar intentos fallidos por
-- identificador y por IP en una ventana de tiempo (anti fuerza bruta).
-- USUARIO_ID es nulo cuando el usuario no existe (intento con otro nombre).
CREATE TABLE AUDITORIA_ACCESO (
  ID             BINARY(16)   NOT NULL,
  FECHA          TIMESTAMP WITH TIME ZONE NOT NULL,
  EVENTO         VARCHAR(30)  NOT NULL,
  USUARIO_ID     BINARY(16),
  IDENTIFICADOR  VARCHAR(254),
  IP             VARCHAR(45),
  USER_AGENT     VARCHAR(255),
  CONSTRAINT PK_AUDITORIA_ACCESO PRIMARY KEY (ID),
  CONSTRAINT FK_AUDITORIA_USUARIO FOREIGN KEY (USUARIO_ID) REFERENCES USUARIO (ID)
);

CREATE INDEX IX_AUDITORIA_IDENT_FECHA ON AUDITORIA_ACCESO (IDENTIFICADOR, FECHA);
CREATE INDEX IX_AUDITORIA_IP_FECHA ON AUDITORIA_ACCESO (IP, FECHA);
