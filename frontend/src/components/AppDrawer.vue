<template>
  <v-navigation-drawer v-model="app.drawer">
    <template #prepend>
      <v-list-item class="py-4" title="SALUS" to="/">
        <template #title>
          <span class="text-h6 font-weight-bold text-primary">SALUS</span>
        </template>
      </v-list-item>
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
  import { menu } from '@/config/menu'
  import { useAppStore } from '@/stores/app'

  const app = useAppStore()

  // Rol fijo hasta la integración con la API (paso 0.6).
  const rol = 'SYSADMIN'

  const secciones = computed(() =>
    menu
      .map(s => ({ ...s, opciones: s.opciones.filter(o => !o.roles || o.roles.includes(rol)) }))
      .filter(s => s.opciones.length > 0),
  )
</script>
