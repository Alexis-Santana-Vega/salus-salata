<template>
  <app-drawer />
  <app-header />

  <v-main>
    <v-container class="pa-4 pa-md-6" fluid>
      <router-view />
    </v-container>
  </v-main>
</template>

<script lang="ts" setup>
  import AppDrawer from '@/components/AppDrawer.vue'
  import AppHeader from '@/components/AppHeader.vue'
</script>
