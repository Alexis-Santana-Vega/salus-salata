<!-- Login (mock: Login / LoginMovil): página blanca sin card, marca arriba y formulario centrado. -->
<template>
  <v-main>
    <v-container class="d-flex flex-column fill-height">
      <v-row class="flex-grow-0" :justify="mobile ? 'center' : 'start'" no-gutters>
        <app-marca />
      </v-row>

      <v-row align="center" class="w-100" justify="center">
        <v-col
          cols="12"
          lg="5"
          md="6"
          sm="10"
          xl="4"
        >
          <v-form ref="form" @submit.prevent="entrar">
            <v-row>
              <v-col class="text-center" cols="12">
                <h1 class="font-weight-semibold" :class="mobile ? 'text-headline-medium' : 'text-display-small'">Bienvenido</h1>
                <p class="text-body-large text-text-muted">Control de renta de equipo médico</p>
              </v-col>

              <v-expand-transition>
                <v-col v-if="error" cols="12">
                  <v-alert color="error" icon="$error" :text="error" />
                </v-col>
              </v-expand-transition>

              <v-col cols="12">
                <v-text-field
                  v-model.trim="usuario"
                  autocomplete="username"
                  label="Usuario o correo"
                  maxlength="254"
                  :rules="[requerido, maximo(254), usuarioOCorreo]"
                />
              </v-col>

              <v-col cols="12">
                <v-text-field
                  v-model="contrasenaValor"
                  :append-inner-icon="verContrasena ? '$ocultar' : '$ver'"
                  autocomplete="current-password"
                  label="Contraseña"
                  maxlength="128"
                  :rules="contrasena"
                  :type="verContrasena ? 'text' : 'password'"
                  @click:append-inner="verContrasena = !verContrasena"
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
                  {{ cargando ? 'Cargando' : 'Entrar' }}
                </btn-primario>
              </v-col>
            </v-row>
          </v-form>
        </v-col>
      </v-row>
    </v-container>
  </v-main>
</template>

<script lang="ts" setup>
  import type { VForm } from 'vuetify/components'
  import { ref } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { useDisplay } from 'vuetify'
  import AppMarca from '@/components/AppMarca.vue'
  import { ApiError } from '@/services/api'
  import { useAuthStore } from '@/stores/auth'
  import { contrasena, maximo, requerido, usuarioOCorreo } from '@/utils/rules'

  const router = useRouter()
  const route = useRoute()
  const auth = useAuthStore()
  const { mobile } = useDisplay()
  const form = ref<VForm>()
  const usuario = ref('')
  const contrasenaValor = ref('')
  const verContrasena = ref(false)
  const cargando = ref(false)
  const error = ref('')

  async function entrar () {
    error.value = ''
    const { valid } = await form.value!.validate()
    if (!valid) return
    cargando.value = true
    try {
      await auth.login(usuario.value, contrasenaValor.value)
      contrasenaValor.value = ''
      // Solo rutas internas: evita redirigir a otro sitio con ?volver=.
      const volver = typeof route.query.volver === 'string' && route.query.volver.startsWith('/') && !route.query.volver.startsWith('//')
        ? route.query.volver
        : '/'
      await router.replace(volver)
    } catch (error_) {
      error.value = error_ instanceof ApiError ? error_.mensaje : 'Algo salió mal. Intenta de nuevo.'
    } finally {
      cargando.value = false
    }
  }
</script>
