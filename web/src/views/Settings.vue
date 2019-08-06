<template>
  <v-container>
    <v-form
      v-model="valid"
    >
      <v-text-field
        v-model="name"
        label="CID"
        prepend-icon="mdi-account"
        :rules="nameRules"
      ></v-text-field>
      <v-text-field
        v-model="password"
        label="Password"
        prepend-icon="mdi-lock"
        :rules="passwordRules"
        type="password"
      ></v-text-field>
      <v-alert
        icon="mdi-alert"
        type="warning"
      >
        Your credentials will be sent to the BooGroCha server.
        The BooGroCha server will only have access to your credentials while performing a request.
      </v-alert>
      <v-layout
        justify-end
      >
        <v-btn
          :disabled="!valid"
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

import CredentialsStorage from '../mixins/CredentialsStorage'

function notEmpty (input) {
  return input !== '';
}
export default {
  name: 'Settings',
  data: () => ({
    valid: false,
    name: "",
    nameRules : [notEmpty],
    password: "",
    passwordRules: [notEmpty],
  }),
  methods: {
    send () {

    },
  },
  mixins: [
    CredentialsStorage
  ]
}
</script>

<style scoped>

</style>
