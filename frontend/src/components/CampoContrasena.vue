<!--
  Campo de contraseña con botón para mostrar/ocultar. Recibe las mismas
  props que v-text-field (label, rules, autocomplete, error-messages…).
-->
<template>
  <v-text-field
    v-model="valor"
    :append-inner-icon="visible ? '$ocultar' : '$ver'"
    maxlength="128"
    :type="visible ? 'text' : 'password'"
    @click:append-inner="visible = !visible"
  />
</template>

<script lang="ts" setup>
  import { ref } from 'vue'

  const valor = defineModel<string>({ default: '' })
  const visible = ref(false)
</script>
