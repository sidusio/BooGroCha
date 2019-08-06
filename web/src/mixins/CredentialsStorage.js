import Storage from './Storage'

let credentialsKey = 'credentials'

export default {
  computed: {
    credentials: {
      get () {
        this.localStorage.getItem(credentialsKey)
      },
      set (credentials) {
        if (credentials === null) {
          this.localStorage.removeItem(credentialsKey)
        } else {
          this.localStorage.setItem(credentialsKey, credentials)
        }
      },
    },
  },
  mixins: [
    Storage,
  ],
}
