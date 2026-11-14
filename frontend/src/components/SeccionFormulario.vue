<!--
  Sección de un formulario: título, ayuda y sus campos en rejilla
  (mock: "1 · Modelo", "Datos generales"…). Entre secciones va un divisor.
-->
<template>
  <v-divider v-if="!primera" />

  <v-card-text>
    <h2 class="text-title-medium font-weight-semibold">{{ titulo }}</h2>
    <p v-if="ayuda" class="text-body-medium text-text-muted">{{ ayuda }}</p>

    <v-row>
      <slot />
    </v-row>
  </v-card-text>
</template>

<script lang="ts" setup>
  defineProps<{ titulo: string, ayuda?: string, primera?: boolean }>()
</script>
