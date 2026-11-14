<template>
  <app-drawer />
  <app-header />
  <app-avisos />

  <v-main>
    <v-container class="fill-height align-start" fluid>
      <div class="w-100 h-100">
        <router-view />
      </div>
    </v-container>
  </v-main>
</template>

<script lang="ts" setup>
  import { watch } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import AppAvisos from '@/components/AppAvisos.vue'
  import AppDrawer from '@/components/AppDrawer.vue'
  import AppHeader from '@/components/AppHeader.vue'
  import { useAuthStore } from '@/stores/auth'

  const auth = useAuthStore()
  const router = useRouter()
  const route = useRoute()

  // Si la sesión se pierde mientras se trabaja (refresh vencido o revocado),
  // vuelve al login recordando dónde estaba.
  watch(() => auth.autenticado, activo => {
    if (!activo) router.replace({ path: '/login', query: { volver: route.fullPath } })
  })
</script>
