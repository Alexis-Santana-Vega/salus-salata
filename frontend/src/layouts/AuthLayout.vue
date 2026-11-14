<!--
  Layout de acceso (login y cambio de contraseña). Web: panel de marca con
  frase y animación a la izquierda; a la derecha la marca centrada, el
  contenido y el pie legal. Móvil: solo la parte derecha.
-->
<template>
  <v-main>
    <v-row class="fill-height" no-gutters>
      <!-- Solo desde md: en móvil no se pinta ni ocupa espacio. -->
      <v-col class="d-none d-md-flex" md="6">
        <v-sheet class="d-flex flex-column justify-center flex-grow-1 pa-4" color="primary">
          <v-row justify="center">
            <v-col cols="10" lg="8">
              <h2 class="text-display-small font-weight-semibold">Todo tu equipo médico, bajo control</h2>
              <p class="text-body-large">Programa cirugías, sabe dónde está cada equipo y cotiza sin errores.</p>
            </v-col>

            <v-col cols="10" lg="8">
              <animacion-lottie src="/hospital-lottie.json" />
            </v-col>
          </v-row>
        </v-sheet>
      </v-col>

      <v-col class="d-flex flex-column" cols="12" md="6">
        <v-container class="d-flex flex-column flex-grow-1">
          <v-row class="flex-grow-0" justify="center">
            <v-col cols="auto">
              <app-marca />
            </v-col>
          </v-row>

          <v-row align="center" class="flex-grow-1" justify="center">
            <v-col cols="12" lg="8" sm="10">
              <slot />
            </v-col>
          </v-row>

          <v-row class="flex-grow-0" justify="center">
            <v-col class="text-center text-body-small text-text-muted" cols="12">
              © {{ anio }} Salus Salata. Todos los derechos reservados.
            </v-col>
          </v-row>
        </v-container>
      </v-col>
    </v-row>
  </v-main>
</template>

<script lang="ts" setup>
  import AnimacionLottie from '@/components/AnimacionLottie.vue'
  import AppMarca from '@/components/AppMarca.vue'

  const anio = new Date().getFullYear()
</script>
