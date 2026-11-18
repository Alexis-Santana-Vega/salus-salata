<!--
  Sección de un formulario: título, ayuda y sus campos en rejilla
  (mock: "1 · Modelo", "Datos generales"…). Entre secciones va un divisor.
  El slot "extra" va junto al título (p. ej. el chip con el total de unidades).
-->
<template>
  <v-divider v-if="!primera" />

  <v-card-text>
    <div class="d-flex flex-wrap align-center ga-2">
      <h2 class="text-title-medium font-weight-semibold">{{ titulo }}</h2>
      <slot name="extra" />
    </div>

    <p v-if="ayuda" class="text-body-medium text-text-muted">{{ ayuda }}</p>

    <v-row>
      <slot />
    </v-row>
  </v-card-text>
</template>

<script lang="ts" setup>
  defineProps<{ titulo: string, ayuda?: string, primera?: boolean }>()
</script>
