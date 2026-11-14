<!-- Sidebar (mock: Sidebar): marca arriba, secciones con título y opciones sola o expandible. -->
<template>
  <v-navigation-drawer v-model="app.drawer">
    <template #prepend>
      <div class="pa-4">
        <app-marca />
      </div>
    </template>

    <v-list>
      <template v-for="seccion in secciones" :key="seccion.titulo">
        <v-list-subheader>{{ seccion.titulo }}</v-list-subheader>

        <template v-for="op in seccion.opciones" :key="op.titulo">
          <!-- Opción expandible: módulo principal y sus altas -->
          <v-list-group v-if="op.hijos" :value="op.titulo">
            <template #activator="{ props }">
              <v-list-item v-bind="props" :prepend-icon="op.icono" :title="op.titulo" />
            </template>

            <v-list-item
              v-for="hijo in op.hijos"
              :key="hijo.to"
              exact
              :title="hijo.titulo"
              :to="hijo.to"
            />
          </v-list-group>

          <!-- Opción sola: lleva directo al módulo -->
          <v-list-item
            v-else
            exact
            :prepend-icon="op.icono"
            :title="op.titulo"
            :to="op.to"
          />
        </template>
      </template>
    </v-list>
  </v-navigation-drawer>
</template>

<script lang="ts" setup>
  import { computed } from 'vue'
  import AppMarca from '@/components/AppMarca.vue'
  import { menu } from '@/config/menu'
  import { useAppStore } from '@/stores/app'
  import { useAuthStore } from '@/stores/auth'

  const app = useAppStore()

  const auth = useAuthStore()

  const secciones = computed(() =>
    menu
      .map(s => ({ ...s, opciones: s.opciones.filter(o => auth.tieneRol(o.roles)) }))
      .filter(s => s.opciones.length > 0),
  )
</script>
