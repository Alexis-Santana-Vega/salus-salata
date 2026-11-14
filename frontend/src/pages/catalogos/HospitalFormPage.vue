<!-- Alta y edición de hospital con sus contactos (mock: AltaHospital). La foto llega en la Fase 9. -->
<template>
  <error-pagina v-if="noEncontrado" codigo="404" />

  <formulario-crud
    v-else
    :cargando="cargando"
    :error="error"
    :etiqueta-guardar="esEdicion ? 'Guardar cambios' : 'Guardar hospital'"
    :guardando="guardando"
    :migas="[{ title: 'Catálogos', to: '/catalogos/hospitales' }, { title: titulo }]"
    subtitulo="Destino de la renta. Su ubicación y contactos aparecen en la ficha de cada cirugía."
    :titulo="titulo"
    @cancelar="volver"
    @guardar="guardar"
  >
    <template v-if="esEdicion" #acciones>
      <v-btn :block="mobile" color="error" prepend-icon="$eliminar" @click="pedirBaja({ id: id! })">Eliminar</v-btn>
    </template>

    <seccion-formulario ayuda="La clave lleva el prefijo del estado de su colonia (HGO-001)." primera titulo="Datos generales">
      <v-col cols="12" md="3">
        <v-text-field
          v-model="datos.clave"
          :error-messages="errores.clave"
          :hint="datos.direccion.estado ? `Debe empezar con ${datos.direccion.estado}-` : 'Ej. HGO-001'"
          label="Clave *"
          maxlength="10"
          persistent-hint
          :rules="[requerido, claveHospital, prefijoEstado]"
          @update:model-value="limpiarError('clave')"
        />
      </v-col>

      <v-col cols="12" md="6">
        <v-text-field
          v-model="datos.nombre"
          autofocus
          :error-messages="errores.nombre"
          label="Nombre del hospital *"
          maxlength="150"
          :rules="[requerido, maximo(150), textoGeneral]"
          @update:model-value="limpiarError('nombre')"
        />
      </v-col>

      <v-col cols="12" md="3">
        <v-text-field
          v-model="datos.telefono"
          :error-messages="errores.telefono"
          inputmode="tel"
          label="Teléfono"
          maxlength="20"
          :rules="[telefono, maximo(20)]"
          type="tel"
          @update:model-value="limpiarError('telefono')"
        />
      </v-col>
    </seccion-formulario>

    <seccion-formulario ayuda="Escribe el C.P. y elige la colonia: el municipio y el estado se llenan solos. La liga del mapa se muestra al equipo que entrega." titulo="Ubicación">
      <campos-direccion v-model="datos.direccion" :errores="errores" obligatoria @limpiar="limpiarError" />

      <v-col cols="6" md="3">
        <v-text-field
          v-model="datos.latitud"
          :error-messages="errores.latitud"
          inputmode="decimal"
          label="Latitud"
          maxlength="12"
          placeholder="20.1011"
          :rules="[rango(-90, 90, 'Latitud')]"
          @update:model-value="limpiarError('latitud')"
        />
      </v-col>

      <v-col cols="6" md="3">
        <v-text-field
          v-model="datos.longitud"
          :error-messages="errores.longitud"
          inputmode="decimal"
          label="Longitud"
          maxlength="12"
          placeholder="-98.7591"
          :rules="[rango(-180, 180, 'Longitud')]"
          @update:model-value="limpiarError('longitud')"
        />
      </v-col>

      <v-col cols="12" md="6">
        <v-text-field
          v-model="datos.ligaMapa"
          :error-messages="errores.ligaMapa"
          inputmode="url"
          label="Liga del mapa"
          maxlength="500"
          placeholder="https://maps.google.com/…"
          :rules="[ligaHttps, maximo(500)]"
          type="url"
          @update:model-value="limpiarError('ligaMapa')"
        />
      </v-col>
    </seccion-formulario>

    <seccion-formulario ayuda="Personas en sitio para coordinar la entrega y la devolución." titulo="Contactos">
      <template v-for="(c, i) in datos.contactos" :key="i">
        <v-col cols="12" md="5">
          <v-text-field
            v-model="c.nombre"
            :error-messages="errores[`contactos[${i}].nombre`]"
            label="Nombre *"
            maxlength="120"
            :rules="[requerido, maximo(120), nombrePersona]"
            @update:model-value="limpiarError(`contactos[${i}].nombre`)"
          />
        </v-col>

        <v-col cols="12" md="3">
          <v-text-field
            v-model="c.puesto"
            :error-messages="errores[`contactos[${i}].puesto`]"
            label="Puesto"
            maxlength="60"
            :rules="[maximo(60), textoGeneral]"
            @update:model-value="limpiarError(`contactos[${i}].puesto`)"
          />
        </v-col>

        <v-col cols="10" md="3">
          <v-text-field
            v-model="c.telefono"
            :error-messages="errores[`contactos[${i}].telefono`]"
            inputmode="tel"
            label="Teléfono"
            maxlength="20"
            :rules="[telefono, maximo(20)]"
            type="tel"
            @update:model-value="limpiarError(`contactos[${i}].telefono`)"
          />
        </v-col>

        <v-col cols="2" md="1">
          <btn-icono :aria-label="`Quitar contacto ${i + 1}`" icon="$eliminar" @click="datos.contactos.splice(i, 1)" />
        </v-col>
      </template>

      <v-col cols="12">
        <v-btn :disabled="datos.contactos.length >= 20" prepend-icon="$agregar" @click="datos.contactos.push({ nombre: '', puesto: '', telefono: '' })">
          Agregar contacto
        </v-btn>
      </v-col>
    </seccion-formulario>
  </formulario-crud>

  <dialogo-eliminar
    v-model="bajaAbierta"
    :eliminando="eliminando"
    :texto="textoBaja"
    titulo="¿Eliminar el hospital?"
    @confirmar="confirmarBaja"
  />
</template>

<script lang="ts" setup>
  import { computed } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { useDisplay } from 'vuetify'
  import CamposDireccion from '@/components/CamposDireccion.vue'
  import DialogoEliminar from '@/components/DialogoEliminar.vue'
  import ErrorPagina from '@/components/ErrorPagina.vue'
  import FormularioCrud from '@/components/FormularioCrud.vue'
  import SeccionFormulario from '@/components/SeccionFormulario.vue'
  import { useBaja } from '@/composables/useBaja'
  import { useFormulario } from '@/composables/useFormulario'
  import { aEnvio, hospitales } from '@/services/hospitales'
  import { direccionVacia } from '@/utils/direccion'
  import { claveHospital, ligaHttps, maximo, nombrePersona, rango, requerido, telefono, textoGeneral } from '@/utils/rules'

  const route = useRoute()
  const router = useRouter()
  const { mobile } = useDisplay()
  const id = computed(() => (route.params.id as string | undefined) || undefined)

  const { datos, errores, error, cargando, guardando, noEncontrado, esEdicion, guardar, limpiarError }
    = useFormulario(hospitales, {
      id,
      vacio: () => ({ clave: '', nombre: '', telefono: '', direccion: direccionVacia(), latitud: '', longitud: '', ligaMapa: '', contactos: [] }),
      desde: h => ({
        clave: h.clave, nombre: h.nombre, telefono: h.telefono, direccion: { ...h.direccion },
        latitud: h.latitud == null ? '' : String(h.latitud), longitud: h.longitud == null ? '' : String(h.longitud),
        ligaMapa: h.ligaMapa, contactos: h.contactos.map(c => ({ ...c })),
      }),
      enviar: aEnvio,
      mensajes: { creado: 'Hospital dado de alta', actualizado: 'Hospital actualizado' },
      alGuardar: volver,
    })

  const titulo = computed(() => (esEdicion.value ? 'Editar hospital' : 'Alta de hospital'))
  const textoBaja = computed(() => `${datos.nombre} dejará de aparecer al programar cirugías. Las cirugías que ya tiene lo conservan.`)

  // Con colonia elegida, la clave debe empezar con el estado de esa colonia (igual que el back).
  function prefijoEstado (v: unknown) {
    const estado = datos.direccion.estado
    const clave = String(v ?? '').trim().toUpperCase()
    return !estado || !clave || clave.startsWith(`${estado}-`) || `La clave debe empezar con ${estado}-`
  }

  // En móvil la baja se hace desde aquí (la lista solo abre la edición).
  const { abierto: bajaAbierta, eliminando, pedir: pedirBaja, confirmar: confirmarBaja }
    = useBaja<{ id: string }>(hospitales.eliminar, 'Hospital eliminado', volver)

  function volver () {
    router.push('/catalogos/hospitales')
  }
</script>
