-- 0014 · Los consecutivos de hospital continúan después de las claves que ya
-- existen (por estado). Las claves de insumo capturadas a mano antes no siguen
-- el formato INS-0001: el consecutivo de insumos empieza en 1.
INSERT INTO CONSECUTIVO (NOMBRE, VALOR)
SELECT 'HOSPITAL:' || SUBSTRING(CLAVE FROM 1 FOR POSITION('-' IN CLAVE) - 1),
       MAX(CAST(SUBSTRING(CLAVE FROM POSITION('-' IN CLAVE) + 1) AS INTEGER))
FROM HOSPITAL
GROUP BY 1;
