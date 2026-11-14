<!--
  Alta y edición de equipo (mock: AltaEquipo, secciones 1 y 2). Unidades: paso
  1.7. Proveedor: Fase 7. Accesorios: Fase 9. Notas: Fase 4.
-->
<template>
  <error-pagina v-if="noEncontrado" codigo="404" />

  <formulario-crud
    v-else
    :cargando="cargando"
    :error="error"
    :etiqueta-guardar="esEdicion ? 'Guardar cambios' : 'Guardar equipo'"
    :guardando="guardando"
    :migas="[{ title: 'Equipos', to: '/equipos' }, { title: titulo }]"
    subtitulo="Se registra el modelo una sola vez; cada unidad llevará su propio uso, mantenimiento y estado."
    :titulo="titulo"
    @cancelar="volver"
    @guardar="guardar"
  >
    <template v-if="esEdicion" #acciones>
      <v-btn :block="mobile" color="error" prepend-icon="$eliminar" @click="pedirBaja({ id: id! })">Eliminar</v-btn>
    </template>

    <seccion-formulario ayuda="Datos que comparten todas las unidades de este equipo. La clave se genera al guardar (EQ-001)." primera titulo="1 · Modelo">
      <v-col v-if="esEdicion" cols="12" md="3">
        <campo-generado label="Clave" :valor="datos.clave" />
      </v-col>

      <v-col cols="12" md="6">
        <v-text-field
          v-model="datos.nombre"
          autofocus
          :error-messages="errores.nombre"
          label="Nombre del equipo *"
          maxlength="120"
          placeholder="Lámpara quirúrgica LED"
          :rules="[requerido, maximo(120), textoGeneral]"
          @update:model-value="limpiarError('nombre')"
        />
      </v-col>

      <v-col cols="12" :md="esEdicion ? 3 : 6">
        <v-text-field
          v-model="datos.marcaModelo"
          :error-messages="errores.marcaModelo"
          label="Marca / modelo *"
          maxlength="120"
          placeholder="Gimmi LX-300"
          :rules="[requerido, maximo(120), textoGeneral]"
          @update:model-value="limpiarError('marcaModelo')"
        />
      </v-col>
    </seccion-formulario>

    <seccion-formulario ayuda="Se cobra por hora de renta. El ciclo y la purga se cuentan por unidad." titulo="2 · Costo de renta y mantenimiento">
      <v-col cols="12" md="4">
        <v-text-field
          v-model="datos.costoHora"
          :error-messages="errores.costoHora"
          :hint="esEdicion ? 'Si lo cambias, se agrega al historial con vigencia desde hoy' : 'Vigente desde hoy'"
          inputmode="decimal"
          label="Costo de renta por hora (MXN) *"
          maxlength="13"
          persistent-hint
          prefix="$"
          :rules="[requerido, dinero]"
          @update:model-value="limpiarError('costoHora')"
        />
      </v-col>

      <v-col cols="12" md="4">
        <v-text-field
          v-model="datos.cicloMax"
          :error-messages="errores.cicloMax"
          inputmode="numeric"
          label="Límite de ciclo (usos) *"
          maxlength="6"
          :rules="[requerido, enteroMayorACero]"
          @update:model-value="limpiarError('cicloMax')"
        />
      </v-col>

      <v-col cols="12" md="4">
        <v-switch v-model="datos.requierePurga" label="Requiere purga" :rules="[]" />
      </v-col>

      <v-expand-transition>
        <v-col v-if="datos.requierePurga" cols="12" md="4">
          <v-text-field
            v-model="datos.purgaCada"
            :error-messages="errores.purgaCada"
            inputmode="numeric"
            label="Purga cada (usos) *"
            maxlength="6"
            :rules="[requerido, enteroMayorACero]"
            @update:model-value="limpiarError('purgaCada')"
          />
        </v-col>
      </v-expand-transition>
    </seccion-formulario>

    <seccion-formulario v-if="esEdicion" ayuda="Las tarifas no se sobrescriben: cada cambio de costo queda aquí." titulo="Historial de tarifas">
      <v-col cols="12">
        <v-table>
          <thead>
            <tr>
              <th>Vigente desde</th>
              <th class="text-end">Costo por hora</th>
              <th />
            </tr>
          </thead>

          <tbody>
            <tr v-for="(t, i) in datos.tarifas" :key="t.vigenteDesde">
              <td>{{ fecha(t.vigenteDesde) }}</td>
              <td class="text-end">{{ moneda(t.costoHora) }}</td>
              <td><v-chip v-if="i === 0" color="success" text="Vigente" /></td>
            </tr>
          </tbody>
        </v-table>
      </v-col>
    </seccion-formulario>
  </formulario-crud>

  <dialogo-eliminar
    v-model="bajaAbierta"
    :eliminando="eliminando"
    :texto="textoBaja"
    titulo="¿Eliminar el equipo?"
    @confirmar="confirmarBaja"
  />
</template>

<script lang="ts" setup>
  import { computed } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { useDisplay } from 'vuetify'
  import CampoGenerado from '@/components/CampoGenerado.vue'
  import DialogoEliminar from '@/components/DialogoEliminar.vue'
  import ErrorPagina from '@/components/ErrorPagina.vue'
  import FormularioCrud from '@/components/FormularioCrud.vue'
  import SeccionFormulario from '@/components/SeccionFormulario.vue'
  import { useBaja } from '@/composables/useBaja'
  import { useFormulario } from '@/composables/useFormulario'
  import { aEnvio, equipos } from '@/services/equipos'
  import { fecha, moneda } from '@/utils/formato'
  import { dinero, enteroMayorACero, maximo, requerido, textoGeneral } from '@/utils/rules'

  const route = useRoute()
  const router = useRouter()
  const { mobile } = useDisplay()
  const id = computed(() => (route.params.id as string | undefined) || undefined)

  const { datos, errores, error, cargando, guardando, noEncontrado, esEdicion, guardar, limpiarError }
    = useFormulario(equipos, {
      id,
      vacio: () => ({ clave: '', nombre: '', marcaModelo: '', cicloMax: '50', requierePurga: false, purgaCada: '5', costoHora: '', tarifas: [] }),
      desde: e => ({
        clave: e.clave, nombre: e.nombre, marcaModelo: e.marcaModelo, cicloMax: String(e.cicloMax),
        requierePurga: e.requierePurga, purgaCada: e.purgaCada == null ? '5' : String(e.purgaCada),
        costoHora: String(e.costoHora), tarifas: e.tarifas,
      }),
      enviar: aEnvio,
      mensajes: { creado: 'Equipo dado de alta', actualizado: 'Equipo actualizado' },
      alGuardar: volver,
    })

  const titulo = computed(() => (esEdicion.value ? 'Editar equipo' : 'Alta de equipo'))
  const textoBaja = computed(() => `${datos.clave} ${datos.nombre} dejará de aparecer al programar cirugías. Su clave no se podrá volver a usar.`)

  // En móvil la baja se hace desde aquí (la lista solo abre la edición).
  const { abierto: bajaAbierta, eliminando, pedir: pedirBaja, confirmar: confirmarBaja }
    = useBaja<{ id: string }>(equipos.eliminar, 'Equipo eliminado', volver)

  function volver () {
    router.push('/equipos')
  }
</script>
