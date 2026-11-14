-- 0009 · Número interno del doctor automático (1 de 3): secuencia.
-- El número ya no se teclea: se genera al dar de alta (DOC-0001, DOC-0002…).
-- Una secuencia no se revierte con la transacción: si un alta falla puede
-- quedar un hueco en la numeración, y eso está bien para un número interno.
CREATE SEQUENCE SEQ_DOCTOR_NUMERO;
