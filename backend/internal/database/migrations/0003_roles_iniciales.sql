-- 0003 · Roles del sistema
-- Va en su propia migración: Firebird no permite usar en la misma
-- transacción una tabla creada en ella (ROL se crea en la 0002).

INSERT INTO ROL (ID, CLAVE, NOMBRE) VALUES (1, 'SYSADMIN', 'Sysadmin');
INSERT INTO ROL (ID, CLAVE, NOMBRE) VALUES (2, 'ADMIN', 'Admin');
INSERT INTO ROL (ID, CLAVE, NOMBRE) VALUES (3, 'DOCTOR', 'Doctor');
