<!--
  Acciones al pie de un formulario (mock: FormActions / FormMovil).
  Web: nota de obligatorios + Cancelar (texto) + Guardar (primario).
  Móvil: Cancelar y Guardar a la mitad cada uno.
  El botón Guardar envía el formulario indicado en `form` (atributo HTML).
-->
<template>
  <v-card-actions v-if="!mobile">
    <span class="text-body-small text-text-muted">Los campos con * son obligatorios.</span>
    <v-spacer />
    <v-btn :disabled="guardando" size="large" @click="emit('cancelar')">Cancelar</v-btn>

    <btn-primario
      :disabled="guardando"
      :form="form"
      :loading="guardando"
      size="large"
      type="submit"
    >
      {{ guardando ? 'Cargando' : etiqueta }}
    </btn-primario>
  </v-card-actions>

  <div v-else class="d-flex ga-2 w-100">
    <v-btn class="flex-1-1" :disabled="guardando" size="large" @click="emit('cancelar')">Cancelar</v-btn>

    <btn-primario
      class="flex-1-1"
      :disabled="guardando"
      :form="form"
      :loading="guardando"
      size="large"
      type="submit"
    >
      {{ guardando ? 'Cargando' : etiqueta }}
    </btn-primario>
  </div>
</template>

<script lang="ts" setup>
  import { useDisplay } from 'vuetify'

  defineProps<{ etiqueta: string, guardando: boolean, form: string }>()
  const emit = defineEmits<{ cancelar: [] }>()

  const { mobile } = useDisplay()
</script>
