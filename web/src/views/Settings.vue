<template>
  <v-container>
    <v-form>
      <v-text-field
        v-model="name"
        label="CID"
        prepend-icon="mdi-account"
        required
      ></v-text-field>
      <v-text-field
        v-model="password"
        label="Password"
        prepend-icon="mdi-lock"
        type="password"
        required
      ></v-text-field>
      <v-alert
        icon="mdi-alert"
        type="warning"
      >
        Your credentials will be sent to the BooGroCha server.
        The BooGroCha server will only have access to your credentials while performing a request.
      </v-alert>
      <v-alert
        icon="mdi-information"
        type="info"
        v-if="hasCredentials"
      >
        Your already have saved credentials.
      </v-alert>
      <v-layout
        justify-end
      >
        <v-btn
          :disabled="!hasCredentials"
          class="ma-1"
          color="error"
          @click="del"
        >
          Delete
        </v-btn>
        <v-btn
          :disabled="!(notEmpty(name) && notEmpty(password))"
          class="ma-1"
          color="success"
          @click="send"
        >
          Save
        </v-btn>
      </v-layout>
    </v-form>
  </v-container>
</template>

<script>
import { mapActions, mapGetters } from 'vuex'

export default {
  name: 'Settings',
  data: () => ({
    name: '',
    password: '',
  }),
  computed: {
    ...mapGetters([
      'hasCredentials',
    ]),
  },
  methods: {
    send () {
      this.axios.post('auth', {
        'cid': this.name,
        'password': this.password,
      }).then(response => {
        if (response.status !== 200) {
          console.log('Resolved with response: ', response)
        } else {
          this.$router.push('/')
        }
      }).catch(reason => {
        console.log('Failed with reason: ', reason)
      }).finally(() => {
        this.updateCredentialsStatus()
      })
      this.name = ''
      this.password = ''
    },
    del () {
      this.axios.delete('auth').then(response => {
        if (response.status !== 200) {
          console.log('Resolved with response: ', response)
        }
      }).catch(reason => {
        console.log('Failed with reason: ', reason)
      }).finally(() => {
        this.updateCredentialsStatus()
      })
    },
    notEmpty (input) {
      return input !== ''
    },
    ...mapActions([
      'updateCredentialsStatus',
    ]),
  },
}
</script>

<style scoped>

</style>
