// Punto de entrada: fuente, plugins y montaje de la app.
import { createApp } from 'vue'
import { registerPlugins } from '@/plugins'
import App from './App.vue'

// Títulos: Poppins 600/700. Texto: Inter 400/500.
import '@fontsource/inter/400.css'
import '@fontsource/inter/500.css'
import '@fontsource/poppins/600.css'
import '@fontsource/poppins/700.css'

const app = createApp(App)

registerPlugins(app)

app.mount('#app')
