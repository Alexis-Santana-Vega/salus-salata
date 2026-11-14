<!--
  Catálogos (mock: Catálogos / CatalogosMovil): encabezado con el botón
  "Nuevo" de la pestaña activa, pestañas por catálogo y su listado debajo.
-->
<template>
  <page-header subtitulo="Hospitales, doctores, técnicos y procedimientos" titulo="Catálogos">
    <template v-if="activa" #acciones>
      <btn-primario :block="mobile" prepend-icon="$agregar" size="large" :to="activa.nuevo.to">
        {{ activa.nuevo.texto }}
      </btn-primario>
    </template>
  </page-header>

  <v-row>
    <v-col cols="12">
      <v-tabs color="primary" show-arrows>
        <v-tab v-for="p in pestanasCatalogos" :key="p.to" :text="p.titulo" :to="p.to" />
      </v-tabs>
    </v-col>

    <v-col cols="12">
      <router-view />
    </v-col>
  </v-row>
</template>

<script lang="ts" setup>
  import { computed } from 'vue'
  import { useRoute } from 'vue-router'
  import { useDisplay } from 'vuetify'
  import PageHeader from '@/components/PageHeader.vue'
  import { pestanasCatalogos } from '@/config/catalogos'

  const route = useRoute()
  const { mobile } = useDisplay()

  const activa = computed(() => pestanasCatalogos.find(p => route.path.startsWith(p.to)))
</script>
