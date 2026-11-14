<!-- Toasts de la app (guía 05 › Retroalimentación): cola de v-snackbar con icono. -->
<template>
  <v-snackbar-queue v-model="avisos.cola" />
</template>

<script lang="ts" setup>
  import { useAvisosStore } from '@/stores/avisos'

  const avisos = useAvisosStore()
</script>
