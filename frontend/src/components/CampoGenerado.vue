<!--
  Identificador que genera el sistema (clave, número interno). Solo lectura:
  el usuario nunca lo teclea. Se muestra en la edición; en el alta no existe aún.
-->
<template>
  <v-text-field
    :counter="false"
    hint="Lo asigna el sistema"
    :label="label"
    :model-value="valor"
    persistent-hint
    readonly
    :rules="[]"
  />
</template>

<script lang="ts" setup>
  defineProps<{ label: string, valor: string }>()
</script>
