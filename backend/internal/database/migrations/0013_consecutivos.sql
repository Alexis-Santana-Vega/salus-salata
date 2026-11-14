-- 0013 · Consecutivos para claves generadas por el sistema.
-- El usuario no teclea claves: hospital (HGO-001, por estado), insumo
-- (INS-0001), equipo (EQ-001) y unidad (EQ-001-01) toman su número de aquí.
-- Se incrementa con UPDATE … RETURNING dentro de la transacción del alta: la
-- fila queda bloqueada hasta confirmar (dos altas no toman el mismo número) y
-- si el alta falla el número no se pierde.
CREATE TABLE CONSECUTIVO (
  NOMBRE  VARCHAR(40) NOT NULL,
  VALOR   INTEGER     NOT NULL,
  CONSTRAINT PK_CONSECUTIVO PRIMARY KEY (NOMBRE),
  CONSTRAINT CK_CONSECUTIVO_VALOR CHECK (VALOR >= 0)
);
