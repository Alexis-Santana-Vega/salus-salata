<!--
  Animación Lottie (caso especial: Vuetify no tiene componente para esto).
  Usa el build "light" de lottie-web: sin eval, compatible con la CSP.
  Si el usuario pidió reducir movimiento, se muestra el primer cuadro fijo.
-->
<template>
  <div ref="contenedor" aria-hidden="true" class="w-100" />
</template>

<script lang="ts" setup>
  import type { AnimationItem } from 'lottie-web'
  import lottie from 'lottie-web/build/player/lottie_light'
  import { onBeforeUnmount, onMounted, ref } from 'vue'

  const props = defineProps<{ src: string }>()

  const contenedor = ref<HTMLElement>()
  let animacion: AnimationItem | undefined

  onMounted(async () => {
    try {
      const datos = await fetch(props.src).then(r => r.json())
      if (!contenedor.value) return // se desmontó mientras cargaba
      const reducir = window.matchMedia('(prefers-reduced-motion: reduce)').matches
      animacion = lottie.loadAnimation({
        container: contenedor.value,
        renderer: 'svg',
        loop: !reducir,
        autoplay: !reducir,
        animationData: datos,
      })
    } catch {
      // Sin animación no pasa nada: es decorativa.
    }
  })

  onBeforeUnmount(() => animacion?.destroy())
</script>
