<!-- Encabezado (mock: AppHeader). Móvil: menú, marca centrada y avatar. Web: menú, buscador, tema y usuario. -->
<template>
  <v-app-bar>
    <template #prepend>
      <btn-icono aria-label="Abrir menú" icon="$menu" @click="app.drawer = !app.drawer" />
    </template>

    <v-app-bar-title v-if="mobile">
      <div class="d-flex justify-center">
        <app-marca tamano="32" />
      </div>
    </v-app-bar-title>

    <campo-busqueda v-else aria-label="Buscar" max-width="460" placeholder="Buscar cirugía, equipo o insumo…">
      <template #append-inner>
        <v-chip color="text-secondary" label>Ctrl K</v-chip>
      </template>
    </campo-busqueda>

    <template #append>
      <btn-icono
        v-if="!mobile"
        :aria-label="app.esOscuro() ? 'Usar modo claro' : 'Usar modo oscuro'"
        :icon="app.esOscuro() ? '$modoClaro' : '$modoOscuro'"
        @click="app.alternarTema()"
      />

      <v-menu>
        <template #activator="{ props }">
          <v-btn
            v-bind="props"
            aria-label="Menú de usuario"

            color="default"
            :icon="mobile"
          >
            <v-avatar color="primary" variant="tonal">
              <span class="text-label-large font-weight-semibold">{{ auth.iniciales }}</span>
            </v-avatar>

            <template v-if="!mobile">
              <span class="d-flex flex-column align-start">
                <span class="text-label-large font-weight-semibold">{{ auth.usuario?.nombre }}</span>
                <span class="text-body-small text-text-muted">{{ auth.usuario ? NOMBRE_ROL[auth.usuario.rol] : "" }}</span>
              </span>

              <v-icon end icon="$dropdown" />
            </template>
          </v-btn>
        </template>

        <v-list>
          <v-list-item
            v-if="mobile"
            :prepend-icon="app.esOscuro() ? '$modoClaro' : '$modoOscuro'"
            :title="app.esOscuro() ? 'Modo claro' : 'Modo oscuro'"
            @click="app.alternarTema()"
          />

          <v-list-item prepend-icon="$salir" title="Cerrar sesión" @click="salir" />
        </v-list>
      </v-menu>
    </template>
  </v-app-bar>
</template>

<script lang="ts" setup>
  import { useRouter } from 'vue-router'
  import { useDisplay } from 'vuetify'
  import AppMarca from '@/components/AppMarca.vue'
  import { useAppStore } from '@/stores/app'
  import { NOMBRE_ROL, useAuthStore } from '@/stores/auth'

  const app = useAppStore()
  const { mobile } = useDisplay()

  const auth = useAuthStore()
  const router = useRouter()

  async function salir () {
    await auth.logout()
    await router.replace('/login')
  }
</script>
