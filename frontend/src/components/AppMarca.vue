<!-- Logo de la marca. Cambia con el tema (claro/oscuro). Solo el ícono, sin el nombre. -->
<template>
  <!-- flex-grow-0: dentro de un d-flex, v-img crece por defecto e ignora el ancho. -->
  <v-img alt="Salus Salata" class="flex-grow-0" :src="logo" :width="tamano" />
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
