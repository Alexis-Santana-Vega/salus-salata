<!-- Alta y edición de técnico (mock: AltaTecnico). -->
<template>
  <error-pagina v-if="noEncontrado" codigo="404" />

  <formulario-crud
    v-else
    :cargando="cargando"
    :error="error"
    :etiqueta-guardar="esEdicion ? 'Guardar cambios' : 'Guardar técnico'"
    :guardando="guardando"
    :migas="[{ title: 'Catálogos', to: '/catalogos/tecnicos' }, { title: titulo }]"
    subtitulo="Personal de apoyo que entrega el material y asiste en la cirugía."
    :titulo="titulo"
    @cancelar="volver"
    @guardar="guardar"
  >
    <template v-if="esEdicion" #acciones>
      <v-btn :block="mobile" color="error" prepend-icon="$eliminar" @click="baja.pedir({ id: id! })">Eliminar</v-btn>
    </template>

    <seccion-formulario ayuda="Se usa para asignar personal y detectar choques de horario." primera titulo="Datos del técnico">
      <v-col cols="12" md="6">
        <v-text-field
          v-model="datos.nombre"
          autofocus
          counter="120"
          :error-messages="errores.nombre"
          label="Nombre *"
          :rules="[requerido, maximo(120), textoSeguro]"
          @update:model-value="limpiarError('nombre')"
        />
      </v-col>

      <v-col cols="12" md="6">
        <v-select
          v-model="datos.puesto"
          :error-messages="errores.puesto"
          :items="PUESTOS"
          label="Puesto *"
          :rules="[requerido, enLista(PUESTOS)]"
          @update:model-value="limpiarError('puesto')"
        />
      </v-col>

      <v-col cols="12" md="6">
        <v-text-field
          v-model="datos.telefono"
          autocomplete="tel"
          :error-messages="errores.telefono"
          inputmode="tel"
          label="Teléfono"
          :rules="[telefono, maximo(20)]"
          type="tel"
          @update:model-value="limpiarError('telefono')"
        />
      </v-col>
    </seccion-formulario>
  </formulario-crud>

  <dialogo-eliminar
    v-model="baja.abierto.value"
    :eliminando="baja.eliminando.value"
    :texto="`${datos.nombre} dejará de aparecer al asignar personal a las cirugías. Las cirugías donde ya participó lo conservan.`"
    titulo="¿Eliminar al técnico?"
    @confirmar="baja.confirmar"
  />
</template>

<script lang="ts" setup>
  import { computed } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { useDisplay } from 'vuetify'
  import DialogoEliminar from '@/components/DialogoEliminar.vue'
  import ErrorPagina from '@/components/ErrorPagina.vue'
  import FormularioCrud from '@/components/FormularioCrud.vue'
  import SeccionFormulario from '@/components/SeccionFormulario.vue'
  import { useBaja } from '@/composables/useBaja'
  import { useFormulario } from '@/composables/useFormulario'
  import { PUESTOS, tecnicos } from '@/services/tecnicos'
  import { enLista, maximo, requerido, telefono, textoSeguro } from '@/utils/rules'

  const route = useRoute()
  const router = useRouter()
  const { mobile } = useDisplay()
  const id = computed(() => (route.params.id as string | undefined) || undefined)

  const { datos, errores, error, cargando, guardando, noEncontrado, esEdicion, guardar, limpiarError }
    = useFormulario(tecnicos, {
      id,
      vacio: () => ({ nombre: '', puesto: null, telefono: '' }),
      desde: t => ({ nombre: t.nombre, puesto: t.puesto, telefono: t.telefono }),
      mensajes: { creado: 'Técnico dado de alta', actualizado: 'Técnico actualizado' },
      alGuardar: volver,
    })

  const titulo = computed(() => (esEdicion.value ? 'Editar técnico' : 'Alta de técnico'))

  // En móvil la baja se hace desde aquí (la lista solo abre la edición).
  const baja = useBaja<{ id: string }>(tecnicos.eliminar, 'Técnico eliminado', volver)

  function volver () {
    router.push('/catalogos/tecnicos')
  }
</script>
