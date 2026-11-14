<template>
  <v-app-bar>
    <template #prepend>
      <btn-icono aria-label="Abrir menú" icon="$menu" @click="app.drawer = !app.drawer" />
    </template>

    <!-- Móvil: logo centrado. Web: buscador. -->
    <v-app-bar-title v-if="mobile" class="text-center">
      <span class="font-weight-bold text-primary">SALUS</span>
    </v-app-bar-title>

    <campo-busqueda
      v-else
      aria-label="Buscar"
      class="mx-2"
      max-width="360"
      placeholder="Buscar  (Ctrl K)"
    />

    <template #append>
      <btn-icono
        :aria-label="app.esOscuro() ? 'Usar modo claro' : 'Usar modo oscuro'"
        :icon="app.esOscuro() ? '$modoClaro' : '$modoOscuro'"
        @click="app.alternarTema()"
      />

      <v-menu>
        <template #activator="{ props }">
          <v-btn v-bind="props" class="ml-1" color="default" :icon="mobile">
            <v-avatar color="primary" size="32" variant="tonal">
              <v-icon icon="$usuario" size="18" />
            </v-avatar>

            <div v-if="!mobile" class="text-start ml-2">
              <div class="text-body-2 font-weight-medium">{{ usuario.nombre }}</div>
              <div class="text-caption text-text-muted">{{ usuario.rol }}</div>
            </div>

            <v-icon v-if="!mobile" class="ml-1" icon="$dropdown" size="16" />
          </v-btn>
        </template>

        <v-list density="compact">
          <v-list-item prepend-icon="$salir" title="Cerrar sesión" to="/login" />
        </v-list>
      </v-menu>
    </template>
  </v-app-bar>
</template>

<script lang="ts" setup>
  import { useDisplay } from 'vuetify'
  import { useAppStore } from '@/stores/app'

  const app = useAppStore()
  const { mobile } = useDisplay()

  // Datos fijos hasta la integración con la API (paso 0.6).
  const usuario = { nombre: 'Administrador', rol: 'Sysadmin' }
</script>
