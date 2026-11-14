<!--
  Página de alta y edición (guía 05 › CRUD, mock: Alta* / FormMovil).
  Migas + PageHeader, card con secciones (slot) y FormActions al pie.
  Móvil: una columna y acciones fijas abajo. Valida con v-form antes de
  emitir "guardar"; los errores del servidor llegan por la prop error.
-->
<template>
  <v-breadcrumbs :items="migas" />

  <page-header :subtitulo="subtitulo" :titulo="titulo">
    <template v-if="$slots.acciones" #acciones>
      <slot name="acciones" />
    </template>
  </page-header>

  <v-skeleton-loader v-if="cargando" type="heading, text, text, actions" />

  <v-card v-else>
    <v-form :id="idForm" ref="form" @submit.prevent="enviar">
      <v-expand-transition>
        <v-card-text v-if="error">
          <alerta-formulario :text="error" />
        </v-card-text>
      </v-expand-transition>

      <slot />
    </v-form>

    <form-actions
      v-if="!mobile"
      :etiqueta="etiquetaGuardar"
      :form="idForm"
      :guardando="guardando"
      @cancelar="emit('cancelar')"
    />
  </v-card>

  <!-- Móvil: acciones fijas abajo -->
  <v-footer v-if="mobile && !cargando" app border="t">
    <form-actions
      :etiqueta="etiquetaGuardar"
      :form="idForm"
      :guardando="guardando"
      @cancelar="emit('cancelar')"
    />
  </v-footer>
</template>

<script lang="ts" setup>
  import type { VForm } from 'vuetify/components'
  import { ref, useId } from 'vue'
  import { useDisplay } from 'vuetify'
  import FormActions from '@/components/FormActions.vue'
  import PageHeader from '@/components/PageHeader.vue'

  defineProps<{
    titulo: string
    subtitulo?: string
    migas: { title: string, to?: string }[]
    etiquetaGuardar: string
    cargando: boolean
    guardando: boolean
    error: string
  }>()

  const emit = defineEmits<{ guardar: [], cancelar: [] }>()

  const { mobile } = useDisplay()
  const form = ref<VForm>()
  const idForm = `form-${useId()}`

  async function enviar () {
    const { valid } = await form.value!.validate()
    if (valid) emit('guardar')
  }
</script>
