<!--
  Campos de una dirección de México (doctor y hospital). Va dentro de una
  SeccionFormulario; cada campo reserva su espacio de mensajes.
-->
<template>
  <v-col cols="12" md="6">
    <v-text-field
      v-model="d.calle"
      autocomplete="address-line1"
      :error-messages="errores.calle"
      label="Calle"
      :rules="[maximo(120), textoSeguro]"
      @update:model-value="emit('limpiar', 'calle')"
    />
  </v-col>

  <v-col cols="6" md="3">
    <v-text-field
      v-model="d.numExterior"
      :error-messages="errores.numExterior"
      label="Núm. exterior"
      :rules="[maximo(20), textoSeguro]"
      @update:model-value="emit('limpiar', 'numExterior')"
    />
  </v-col>

  <v-col cols="6" md="3">
    <v-text-field
      v-model="d.numInterior"
      :error-messages="errores.numInterior"
      label="Núm. interior"
      :rules="[maximo(20), textoSeguro]"
      @update:model-value="emit('limpiar', 'numInterior')"
    />
  </v-col>

  <v-col cols="12" md="6">
    <v-text-field
      v-model="d.colonia"
      :error-messages="errores.colonia"
      label="Colonia"
      :rules="[maximo(120), textoSeguro]"
      @update:model-value="emit('limpiar', 'colonia')"
    />
  </v-col>

  <v-col cols="12" md="6">
    <v-text-field
      v-model="d.cp"
      autocomplete="postal-code"
      :error-messages="errores.cp"
      inputmode="numeric"
      label="C.P."
      maxlength="5"
      :rules="[cp]"
      @update:model-value="emit('limpiar', 'cp')"
    />
  </v-col>

  <v-col cols="12" md="6">
    <v-text-field
      v-model="d.municipio"
      :error-messages="errores.municipio"
      label="Ciudad o municipio"
      :rules="[maximo(120), textoSeguro]"
      @update:model-value="emit('limpiar', 'municipio')"
    />
  </v-col>

  <v-col cols="12" md="6">
    <v-autocomplete
      v-model="d.estado"
      clearable
      :error-messages="errores.estado"
      :items="ESTADOS"
      label="Estado"
      :rules="[enLista(ESTADOS)]"
      @update:model-value="emit('limpiar', 'estado')"
    />
  </v-col>
</template>

<script lang="ts" setup>
  import { type Direccion, ESTADOS } from '@/utils/estados'
  import { cp, enLista, maximo, textoSeguro } from '@/utils/rules'

  defineProps<{ errores: Record<string, string[] | undefined> }>()
  const emit = defineEmits<{ limpiar: [campo: string] }>()
  const d = defineModel<Direccion>({ required: true })
</script>
