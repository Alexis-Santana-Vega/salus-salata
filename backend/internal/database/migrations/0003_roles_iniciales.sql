-- 0003 · Roles del sistema
-- Va en su propia migración: Firebird no permite usar en la misma
-- transacción una tabla creada en ella (ROL se crea en la 0002).

INSERT INTO ROL (ID, CLAVE, NOMBRE) VALUES (GEN_UUID(), 'SYSADMIN', 'Sysadmin');
INSERT INTO ROL (ID, CLAVE, NOMBRE) VALUES (GEN_UUID(), 'ADMIN', 'Admin');
INSERT INTO ROL (ID, CLAVE, NOMBRE) VALUES (GEN_UUID(), 'DOCTOR', 'Doctor');
