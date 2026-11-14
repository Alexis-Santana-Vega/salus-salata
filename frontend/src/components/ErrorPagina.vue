<!--
  Página de error (mock: Errores / ErrorMovil). Se pinta dentro del área de
  contenido del layout, centrada; nunca reemplaza el sidebar ni el header.
-->
<template>
  <v-row align-content="center" class="fill-height text-center" justify="center">
    <v-col aria-hidden="true" class="codigo-error font-weight-bold text-brand" cols="12">{{ codigo }}</v-col>

    <v-col cols="12" md="6" sm="8">
      <h1 class="font-weight-semibold" :class="mobile ? 'text-title-large' : 'text-headline-large'">{{ info.titulo }}</h1>
      <p class="text-body-large text-text-muted">{{ info.texto }}</p>
    </v-col>

    <v-col cols="12" />

    <v-col v-if="info.secundaria" :cols="mobile ? 12 : 'auto'">
      <v-btn :block="mobile" size="large" @click="secundaria">{{ info.secundaria }}</v-btn>
    </v-col>

    <v-col :cols="mobile ? 12 : 'auto'">
      <btn-primario :block="mobile" size="large" to="/">Volver al inicio</btn-primario>
    </v-col>

    <v-col class="text-body-small text-text-muted" cols="12">
      Código de error {{ codigo }} · Si el problema continúa, comuníquese con el administrador del sistema.
    </v-col>
  </v-row>
</template>

<script lang="ts" setup>
  import { computed } from 'vue'
  import { useRouter } from 'vue-router'
  import { useDisplay } from 'vuetify'

  const props = defineProps<{ codigo: '404' | '403' | '500' | '503' }>()

  const router = useRouter()
  const { mobile } = useDisplay()

  // Textos del mock.
  const textos = {
    404: { titulo: 'Página no encontrada', texto: 'La página que busca no existe o cambió de dirección. Revise el enlace o vuelva al inicio.', secundaria: '' },
    403: { titulo: 'Sin permiso', texto: 'Su rol no tiene acceso a esta sección. Si cree que es un error, pida acceso al administrador.', secundaria: 'Cambiar de cuenta' },
    500: { titulo: 'Algo salió mal', texto: 'Ocurrió un error en el servidor. Ya lo registramos; intente de nuevo en unos minutos.', secundaria: 'Reintentar' },
    503: { titulo: 'Servicio en mantenimiento', texto: 'Estamos actualizando la plataforma. Volverá a estar disponible en breve.', secundaria: 'Reintentar' },
  }
  const info = computed(() => textos[props.codigo])

  function secundaria () {
    if (props.codigo === '403') router.push('/login')
    else router.go(0)
  }
</script>

<style scoped>
/* Caso especial: el código gigante del mock (96 px móvil → 168 px web) no existe en la escala de Vuetify. */
.codigo-error {
  font-size: clamp(6rem, 14vw, 10.5rem);
  line-height: 1;
  letter-spacing: -0.03em;
}
</style>
