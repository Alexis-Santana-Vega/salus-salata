<!-- Logo de la marca. Cambia con el tema (claro/oscuro). Solo el ícono, sin el nombre. -->
<template>
  <v-img alt="Salus Salata" :src="logo" :width="tamano" />
</template>

<script lang="ts" setup>
  import { computed } from 'vue'
  import { useTheme } from 'vuetify'
  import logoClaro from '@/assets/logo-claro.webp'
  import logoOscuro from '@/assets/logo-oscuro.webp'

  withDefaults(defineProps<{ tamano?: number | string }>(), { tamano: 40 })

  const theme = useTheme()
  const logo = computed(() => theme.global.current.value.dark ? logoOscuro : logoClaro)
</script>
