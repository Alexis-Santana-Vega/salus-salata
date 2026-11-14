<!--
  Cambio de contraseña obligatorio (primer acceso). El mock no tiene esta
  pantalla: sigue el mismo patrón del login (marca arriba, formulario centrado).
-->
<template>
  <auth-layout>
    <v-form ref="form" @submit.prevent="guardar">
      <v-row>
        <v-col class="text-center" cols="12">
          <h1 class="font-weight-semibold" :class="mobile ? 'text-headline-medium' : 'text-display-small'">Cambia tu contraseña</h1>
          <p class="text-body-large text-text-muted">Por seguridad, define una contraseña nueva antes de continuar.</p>
        </v-col>

        <v-expand-transition>
          <v-col v-if="error" cols="12">
            <v-alert color="error" icon="$error" :text="error" />
          </v-col>
        </v-expand-transition>

        <v-col cols="12">
          <campo-contrasena
            v-model="actual"
            autocomplete="current-password"
            :error-messages="errores.contrasenaActual"
            label="Contraseña actual"
            :rules="contrasena"
            @update:model-value="delete errores.contrasenaActual"
          />
        </v-col>

        <v-col cols="12">
          <campo-contrasena
            v-model="nueva"
            autocomplete="new-password"
            :error-messages="errores.contrasenaNueva"
            hint="Al menos 12 caracteres, sin tu usuario ni tu correo"
            label="Contraseña nueva"
            persistent-hint
            :rules="reglasNueva"
            @update:model-value="delete errores.contrasenaNueva"
          />
        </v-col>

        <v-col cols="12">
          <campo-contrasena
            v-model="confirmacion"
            autocomplete="new-password"
            label="Confirma la contraseña nueva"
            :rules="[requerido, (v: string) => v === nueva || 'Las contraseñas no coinciden']"
          />
        </v-col>

        <v-col cols="12">
          <btn-primario
            block
            :disabled="cargando"
            :loading="cargando"
            size="large"
            type="submit"
          >
            {{ cargando ? 'Cargando' : 'Guardar contraseña' }}
          </btn-primario>
        </v-col>

        <v-col cols="12">
          <v-btn block size="large" @click="salir">Cerrar sesión</v-btn>
        </v-col>
      </v-row>
    </v-form>
  </auth-layout>
</template>

<script lang="ts" setup>
  import type { VForm } from 'vuetify/components'
  import { computed, reactive, ref } from 'vue'
  import { useRouter } from 'vue-router'
  import { useDisplay } from 'vuetify'
  import CampoContrasena from '@/components/CampoContrasena.vue'
  import AuthLayout from '@/layouts/AuthLayout.vue'
  import { ApiError } from '@/services/api'
  import { useAuthStore } from '@/stores/auth'
  import { contrasena, contrasenaNueva, requerido, sinDatosPersonales } from '@/utils/rules'

  const router = useRouter()
  const auth = useAuthStore()
  const { mobile } = useDisplay()
  const form = ref<VForm>()
  const actual = ref('')
  const nueva = ref('')
  const confirmacion = ref('')
  const cargando = ref(false)
  const error = ref('')
  const errores = reactive<Record<string, string[]>>({})

  // Mismas reglas que el backend, más que la nueva sea distinta de la actual.
  const reglasNueva = computed(() => [
    ...contrasenaNueva,
    (v: unknown) => v !== actual.value || 'Debe ser distinta de la contraseña actual',
    (v: unknown) => sinDatosPersonales(String(v ?? ''), auth.usuario?.usuario, auth.usuario?.correo),
  ])

  async function guardar () {
    error.value = ''
    for (const k of Object.keys(errores)) delete errores[k]
    const { valid } = await form.value!.validate()
    if (!valid) return
    cargando.value = true
    try {
      await auth.cambiarContrasena(actual.value, nueva.value)
      await router.replace('/')
    } catch (error_) {
      // Errores por campo del servidor junto a su campo; el resto, arriba.
      if (error_ instanceof ApiError && error_.detalles.length > 0) {
        for (const d of error_.detalles) errores[d.field] = [d.message]
        error.value = 'Revisa los campos marcados.'
      } else {
        error.value = error_ instanceof ApiError ? error_.mensaje : 'Algo salió mal. Intenta de nuevo.'
      }
    } finally {
      cargando.value = false
    }
  }

  async function salir () {
    await auth.logout()
    await router.replace('/login')
  }
</script>
