<template>
  <v-container
    d-flex
    align-content-center
    justify-center
    fill-height
  >
    <v-time-picker
      no-title
      :allowed-minutes="[0, 15, 30, 45]"
      v-model="value"
      format="24hr"
      :min="min"
      :max="max"
      class="elevation-10"
      @change="select"
    ></v-time-picker>
  </v-container>
</template>

<script>
import moment from 'moment'

export default {
  name: 'TimeSelector',
  data: () => ({
    today: moment().format('YYYY-MM-DD'),
    value: new Date()
  }),
  props: {
    min: String,
    max: String,
  },
  methods: {
    select(time) {
      this.$emit('picked', time)
      this.value = null
    },
  }
}
</script>

<style scoped lang="scss">
</style>
