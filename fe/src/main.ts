import { createApp } from 'vue'
import { createPinia } from 'pinia'
import './assets/styles/main.css'
import './plugins/leaflet'
import App from './App.vue'
import router from './router'
import vClickOutside from './components/directives/v-click-outside'
import { vPermission } from './components/directives/v-permission'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.directive('click-outside', vClickOutside)
app.directive('permission', vPermission)

app.mount('#app')
