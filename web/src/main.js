import Vue from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify'
import router from './router'
import store from './store'
import axios from 'axios'
import VueAxios from 'vue-axios'

axios.defaults.withCredentials = true
axios.defaults.baseURL = 'http://localhost:8081/api/'
axios.defaults.responseType = 'json'
axios.defaults.headers['Accept'] = 'application/json'

Vue.use(VueAxios, axios)

Vue.config.productionTip = false

new Vue({
  vuetify,
  router,
  store,
  render: h => h(App),
}).$mount('#app')
