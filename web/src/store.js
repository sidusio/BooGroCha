import Vuex from 'vuex'
import Vue from 'vue'

Vue.use(Vuex)

const store = new Vuex.Store({
  state: {
    credentialsStatus: {},
  },
  getters: {
    hasCredentials (state) {
      return state.credentialsStatus.HasCookie === true
    },
  },
  mutations: {
    setCredentialsStatus (state, status) {
      state.credentialsStatus = status
    },
  },
  actions: {
    updateCredentialsStatus (context) {
      return new Promise((resolve, reject) => {
        Vue.axios.get('auth/test').then(response => {
          context.commit('setCredentialsStatus', response.data)
          resolve()
        }).catch(reason => {
          console.log('Failed with reason: ', reason)
          reject(reason)
        })
      })
    },
  },
})
export default store
