<!--
  Dirección de México (doctor y hospital). Al escribir el C.P. se consultan sus
  colonias en el catálogo SEPOMEX: si hay una sola se elige sola, y el
  municipio y el estado se llenan solos (no se escriben a mano).
-->
<template>
  <v-col cols="12" md="6">
    <v-text-field
      v-model="d.calle"
      autocomplete="address-line1"
      :error-messages="errores.calle"
      :label="obligatoria ? 'Calle *' : 'Calle'"
      maxlength="120"
      :rules="[...(obligatoria ? [requerido] : []), maximo(120), textoGeneral]"
      @update:model-value="emit('limpiar', 'calle')"
    />
  </v-col>

  <v-col cols="6" md="3">
    <v-text-field
      v-model="d.numExterior"
      :error-messages="errores.numExterior"
      label="Núm. exterior"
      maxlength="20"
      :rules="[maximo(20), numDomicilio]"
      @update:model-value="emit('limpiar', 'numExterior')"
    />
  </v-col>

  <v-col cols="6" md="3">
    <v-text-field
      v-model="d.numInterior"
      :error-messages="errores.numInterior"
      label="Núm. interior"
      maxlength="20"
      :rules="[maximo(20), numDomicilio]"
      @update:model-value="emit('limpiar', 'numInterior')"
    />
  </v-col>

  <v-col cols="12" md="3">
    <v-text-field
      v-model="cpCaptura"
      autocomplete="postal-code"
      :counter="false"
      :error-messages="errorCP"
      hint="Llena la colonia, el municipio y el estado"
      inputmode="numeric"
      :label="obligatoria ? 'C.P. *' : 'C.P.'"
      :loading="buscando"
      maxlength="5"
      :rules="[...(obligatoria ? [requerido] : []), cp]"
    />
  </v-col>

  <v-col cols="12" md="9">
    <v-autocomplete
      v-model="d.asentamientoId"
      :disabled="colonias.length === 0"
      :error-messages="errores.asentamientoId"
      item-title="nombre"
      item-value="id"
      :items="colonias"
      :label="obligatoria ? 'Colonia *' : 'Colonia'"
      no-data-text="Sin colonias para ese C.P."
      :rules="[reglaColonia]"
      @update:model-value="elegir"
    >
      <template #item="{ props: p, item }">
        <v-list-item v-bind="p" :subtitle="item.tipo" />
      </template>
    </v-autocomplete>
  </v-col>

  <v-col cols="12" md="6">
    <v-text-field
      :counter="false"
      label="Ciudad o municipio"
      :model-value="d.municipio"
      readonly
      :rules="[]"
    />
  </v-col>

  <v-col cols="12" md="6">
    <v-text-field
      :counter="false"
      label="Estado"
      :model-value="d.estadoNombre"
      readonly
      :rules="[]"
    />
  </v-col>
</template>

<script lang="ts" setup>
  import type { Asentamiento } from '@/services/direcciones'
  import type { Direccion } from '@/utils/direccion'
  import { onMounted, ref, shallowRef, watch } from 'vue'
  import { ApiError } from '@/services/api'
  import { coloniasPorCP } from '@/services/direcciones'
  import { cp, maximo, numDomicilio, requerido, textoGeneral } from '@/utils/rules'

  const props = defineProps<{
    errores: Record<string, string[] | undefined>
    obligatoria?: boolean // hospital: calle y colonia obligatorias
  }>()
  const emit = defineEmits<{ limpiar: [campo: string] }>()
  const d = defineModel<Direccion>({ required: true })

  const cpCaptura = ref(d.value.cp)
  const colonias = shallowRef<Asentamiento[]>([])
  const buscando = ref(false)
  const errorCP = ref<string[]>([])

  // Con C.P. capturado, la colonia es obligatoria (es lo que se guarda).
  const reglaColonia = () => ((!props.obligatoria && !cpCaptura.value) || !!d.value.asentamientoId) || 'Elige la colonia'

  let turno = 0
  async function buscar (valor: string) {
    const mio = ++turno
    errorCP.value = []
    if (!/^\d{5}$/.test(valor)) {
      colonias.value = []
      return
    }
    buscando.value = true
    try {
      const r = await coloniasPorCP(valor)
      if (mio !== turno) return
      colonias.value = r.asentamientos
      // Si el C.P. tiene una sola colonia, se elige sola.
      if (r.asentamientos.length === 1) elegir(r.asentamientos[0]!.id)
    } catch (error_) {
      if (mio !== turno) return
      colonias.value = []
      errorCP.value = [error_ instanceof ApiError && error_.status === 404
        ? 'Ese C.P. no está en el catálogo (Hidalgo y estados vecinos)'
        : 'No pudimos consultar el C.P. Intenta de nuevo.']
    } finally {
      if (mio === turno) buscando.value = false
    }
  }

  // Al elegir una colonia se copian del catálogo C.P., municipio y estado.
  function elegir (id: string | null) {
    const a = colonias.value.find(x => x.id === id)
    d.value.asentamientoId = a?.id ?? null
    d.value.cp = a?.cp ?? ''
    d.value.colonia = a?.nombre ?? ''
    d.value.municipio = a?.municipio ?? ''
    d.value.estado = a?.estado ?? ''
    d.value.estadoNombre = a?.estadoNombre ?? ''
    emit('limpiar', 'asentamientoId')
  }

  // Cambiar el C.P. quita la colonia elegida y busca las nuevas.
  watch(cpCaptura, valor => {
    if (valor !== d.value.cp) elegir(null)
    buscar(valor)
  })

  // Al editar: se cargan las colonias del C.P. guardado para poder cambiarla.
  onMounted(() => {
    if (d.value.cp) buscar(d.value.cp)
  })
</script>
