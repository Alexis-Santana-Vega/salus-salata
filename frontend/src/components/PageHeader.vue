<!--
  Encabezado de página (mock: Componentes › PageHeader).
  Web: título + subtítulo a la izquierda, acciones a la derecha.
  Móvil: acciones debajo, a todo el ancho.
-->
<template>
  <v-row align="center">
    <v-col cols="12" sm>
      <h1 class="font-weight-semibold" :class="mobile ? 'text-title-large' : 'text-headline-small'">{{ titulo }}</h1>
      <p v-if="subtitulo" class="text-body-medium text-text-muted">{{ subtitulo }}</p>
    </v-col>

    <v-col v-if="$slots.acciones" class="d-flex flex-column flex-sm-row" cols="12" sm="auto">
      <slot name="acciones" />
    </v-col>
  </v-row>
</template>

<script lang="ts" setup>
  import { useDisplay } from 'vuetify'

  defineProps<{ titulo: string, subtitulo?: string }>()

  const { mobile } = useDisplay()
</script>
