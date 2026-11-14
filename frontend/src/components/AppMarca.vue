<!-- Logo + nombre. El logo cambia con el tema (claro/oscuro). -->
<template>
  <div class="d-flex align-center">
    <v-img alt="" :src="logo" width="40" />
    <span class="text-title-large font-weight-semibold">Salus Salata</span>
  </div>
</template>

<script lang="ts" setup>
  import { computed } from 'vue'
  import { useTheme } from 'vuetify'
  import logoClaro from '@/assets/logo-claro.webp'
  import logoOscuro from '@/assets/logo-oscuro.webp'

  const theme = useTheme()
  const logo = computed(() => theme.global.current.value.dark ? logoOscuro : logoClaro)
</script>
