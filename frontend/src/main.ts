// Punto de entrada: fuente, plugins y montaje de la app.
import { createApp } from 'vue'
import { registerPlugins } from '@/plugins'
import App from './App.vue'

// Poppins solo con los pesos que usa la guía de diseño.
import '@fontsource/poppins/400.css'
import '@fontsource/poppins/500.css'
import '@fontsource/poppins/600.css'
import '@fontsource/poppins/700.css'

const app = createApp(App)

registerPlugins(app)

app.mount('#app')
