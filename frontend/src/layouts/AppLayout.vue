<template>
  <app-drawer />
  <app-header />

  <v-main>
    <v-container class="fill-height align-start" fluid>
      <div class="w-100 h-100">
        <router-view />
      </div>
    </v-container>
  </v-main>
</template>

<script lang="ts" setup>
  import AppDrawer from '@/components/AppDrawer.vue'
  import AppHeader from '@/components/AppHeader.vue'
</script>
