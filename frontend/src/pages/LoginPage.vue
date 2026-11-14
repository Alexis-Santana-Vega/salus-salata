<template>
  <v-main>
    <v-container class="fill-height">
      <v-row justify="center">
        <v-col cols="12" lg="4" md="5" sm="8">
          <div class="text-center mb-6">
            <div class="text-h4 font-weight-bold text-primary">SALUS</div>
            <div class="text-body-2 text-text-muted">Renta de equipo médico</div>
          </div>

          <v-card class="pa-6">
            <h1 class="text-h6 font-weight-bold mb-1">Iniciar sesión</h1>
            <p class="text-body-2 text-text-muted mb-6">Entra con tu usuario o correo.</p>

            <v-form ref="form" @submit.prevent="entrar">
              <v-row>
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
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </v-main>
</template>

<script lang="ts" setup>
  import type { VForm } from 'vuetify/components'
  import { ref } from 'vue'
  import { useRouter } from 'vue-router'
  import { contrasena, maximo, requerido, usuarioOCorreo } from '@/utils/rules'

  const router = useRouter()
  const form = ref<VForm>()
  const usuario = ref('')
  const contrasenaValor = ref('')
  const verContrasena = ref(false)
  const cargando = ref(false)

  // Valida en el front; la llamada real a /auth/login llega en el paso 0.6.
  async function entrar () {
    const { valid } = await form.value!.validate()
    if (!valid) return
    cargando.value = true
    try {
      await router.push('/')
    } finally {
      cargando.value = false
    }
  }
</script>
