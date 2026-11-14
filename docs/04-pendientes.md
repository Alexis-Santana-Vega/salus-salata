# 04 · Pendientes y supuestos

## Pendientes
No hay pendientes abiertos.

### Resuelto: cobro del equipo (2026-11-11)
**El equipo se cobra solo por hora.** Cada equipo tiene un costo de renta por hora y el importe es: horas de la ventana entrega → devolución × costo por hora. Se descartó el cobro por evento; el modelo guarda `TARIFA_EQUIPO.costo_hora` con su fecha de vigencia.

## Supuestos tomados (no requieren respuesta)

| Tema | Supuesto |
|---|---|
| Fracciones de hora | El cálculo usa la ventana real entrega → devolución; la regla de redondeo se fija al implementar la cotización. |
| Renglones de la ficha V2 (Láser, Torre, Insumos, Otros) | La ficha mezcla equipo e insumos. En el sistema cada renglón se clasifica como **equipo** o como **insumo**, y el usuario lo asigna. |
| Plantillas | Pueden existir plantillas por procedimiento o kit con equipo e insumos predefinidos y configurables. Hoy no tienen contenido y el sistema funciona sin ellas. |
| User SALUS / User Doc | En espera de definición, fuera de alcance. |
| Lo económico | Solo cotización correcta, gastos y ganancia por cirugía. Sin facturación, pagos ni contabilidad. |
| Semáforo | Ciclo de 50 usos y purga cada 5 son valores por defecto que se configuran por equipo. |
| Retorno de insumos | Lo retornado vuelve al stock y no se cotiza. |
| Uso de equipo | Un uso es una cirugía en la que se usó el equipo. |
