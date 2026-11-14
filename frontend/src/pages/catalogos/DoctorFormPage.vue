<!-- Alta y edición de doctor (mock: AltaDoctor). Hospitales donde opera y acceso al portal: Fases 2 y 6. -->
<template>
  <error-pagina v-if="noEncontrado" codigo="404" />

  <formulario-crud
    v-else
    :cargando="cargando"
    :error="error"
    :etiqueta-guardar="esEdicion ? 'Guardar cambios' : 'Guardar doctor'"
    :guardando="guardando"
    :migas="[{ title: 'Catálogos', to: '/catalogos/doctores' }, { title: titulo }]"
    subtitulo="Médico encargado de las cirugías."
    :titulo="titulo"
    @cancelar="volver"
    @guardar="guardar"
  >
    <template v-if="esEdicion" #acciones>
      <v-btn :block="mobile" color="error" prepend-icon="$eliminar" @click="baja.pedir({ id: id! })">Eliminar</v-btn>
    </template>

    <seccion-formulario ayuda="El correo se usará para darle acceso a su agenda más adelante." primera titulo="Datos del doctor">
      <v-col cols="12" md="6">
        <v-text-field
          v-model="datos.nombre"
          autofocus
          :error-messages="errores.nombre"
          label="Nombre del doctor *"
          maxlength="120"
          :rules="[requerido, maximo(120), nombrePersona]"
          @update:model-value="limpiarError('nombre')"
        />
      </v-col>

      <!-- El número interno lo genera el sistema al dar de alta (DOC-0001). -->
      <v-col v-if="esEdicion" cols="12" md="6">
        <v-text-field
          :counter="false"
          hint="Lo asigna el sistema"
          label="Número interno"
          :model-value="datos.numero"
          persistent-hint
          readonly
          :rules="[]"
        />
      </v-col>

      <v-col cols="12" md="6">
        <v-text-field
          v-model="datos.correo"
          autocomplete="email"
          :error-messages="errores.correo"
          inputmode="email"
          label="Correo *"
          maxlength="254"
          :rules="[requerido, maximo(254), correo]"
          type="email"
          @update:model-value="limpiarError('correo')"
        />
      </v-col>

      <v-col cols="12" md="3">
        <v-text-field
          v-model="datos.telPersonal"
          :error-messages="errores.telPersonal"
          inputmode="tel"
          label="Teléfono personal"
          maxlength="20"
          :rules="[telefono, maximo(20)]"
          type="tel"
          @update:model-value="limpiarError('telPersonal')"
        />
      </v-col>

      <v-col cols="12" md="3">
        <v-text-field
          v-model="datos.telConsultorio"
          :error-messages="errores.telConsultorio"
          inputmode="tel"
          label="Teléfono de consultorio"
          maxlength="20"
          :rules="[telefono, maximo(20)]"
          type="tel"
          @update:model-value="limpiarError('telConsultorio')"
        />
      </v-col>
    </seccion-formulario>

    <seccion-formulario ayuda="Opcional. Escribe el C.P. y elige la colonia: el municipio y el estado se llenan solos." titulo="Dirección">
      <campos-direccion v-model="datos.direccion" :errores="errores" @limpiar="limpiarError" />
    </seccion-formulario>
  </formulario-crud>

  <dialogo-eliminar
    v-model="baja.abierto.value"
    :eliminando="baja.eliminando.value"
    :texto="`${datos.nombre} dejará de aparecer al programar cirugías. Las cirugías que ya tiene lo conservan.`"
    titulo="¿Eliminar al doctor?"
    @confirmar="baja.confirmar"
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
  import { doctores } from '@/services/doctores'
  import { direccionVacia } from '@/utils/direccion'
  import { correo, maximo, nombrePersona, requerido, telefono } from '@/utils/rules'

  const route = useRoute()
  const router = useRouter()
  const { mobile } = useDisplay()
  const id = computed(() => (route.params.id as string | undefined) || undefined)

  const { datos, errores, error, cargando, guardando, noEncontrado, esEdicion, guardar, limpiarError }
    = useFormulario(doctores, {
      id,
      vacio: () => ({ numero: '', nombre: '', correo: '', telPersonal: '', telConsultorio: '', direccion: direccionVacia() }),
      desde: d => ({
        numero: d.numero, nombre: d.nombre, correo: d.correo,
        telPersonal: d.telPersonal, telConsultorio: d.telConsultorio,
        direccion: { ...d.direccion },
      }),
      mensajes: { creado: 'Doctor dado de alta', actualizado: 'Doctor actualizado' },
      alGuardar: volver,
    })

  const titulo = computed(() => (esEdicion.value ? 'Editar doctor' : 'Alta de doctor'))

  // En móvil la baja se hace desde aquí (la lista solo abre la edición).
  const baja = useBaja<{ id: string }>(doctores.eliminar, 'Doctor eliminado', volver)

  function volver () {
    router.push('/catalogos/doctores')
  }
</script>
