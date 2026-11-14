-- 0010 · Número interno del doctor automático (2 de 3): numera a los doctores
-- que ya existen, en orden de alta, con el formato nuevo.
UPDATE DOCTOR
SET NUMERO = 'DOC-' || LPAD(CAST(NEXT VALUE FOR SEQ_DOCTOR_NUMERO AS VARCHAR(10)), 4, '0')
ORDER BY CREADO_EN;
