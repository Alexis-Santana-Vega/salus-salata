<!--
  Confirmación de baja (guía 05 › Eliminar): icono en pastel rojo, texto que
  explica la consecuencia, Cancelar (texto) y Eliminar sólido en rojo.
-->
<template>
  <v-dialog v-model="abierto" :persistent="eliminando">
    <v-card>
      <v-card-item>
        <template #prepend>
          <v-avatar color="error" icon="$eliminar" variant="tonal" />
        </template>

        <v-card-title>{{ titulo }}</v-card-title>
      </v-card-item>

      <v-card-text>
        <p class="text-body-medium text-text-secondary">{{ texto }}</p>
      </v-card-text>

      <v-card-actions>
        <v-spacer />
        <v-btn :disabled="eliminando" @click="abierto = false">Cancelar</v-btn>

        <btn-peligro :disabled="eliminando" :loading="eliminando" @click="emit('confirmar')">
          {{ eliminando ? 'Cargando' : 'Eliminar' }}
        </btn-peligro>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  defineProps<{ titulo: string, texto: string, eliminando: boolean }>()
  const emit = defineEmits<{ confirmar: [] }>()
  const abierto = defineModel<boolean>({ required: true })
</script>
