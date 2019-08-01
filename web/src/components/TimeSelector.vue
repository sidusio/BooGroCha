<template>
  <v-container
    d-flex
    align-content-center
    justify-center
    fill-height
  >
    <v-date-picker
      no-title
      first-day-of-week="1"
      :min="today"
      class="elevation-10"
      @click:date="select"
    ></v-date-picker>
  </v-container>
</template>

<script>
import moment from 'moment'

export default {
  name: 'TimeSelector',
  data: () => ({
    today: moment().format('YYYY-MM-DD'),
  }),
  methods: {
    select(date) {
      this.$emit('picked', date)
    },
  }
}
</script>

<style scoped lang="scss">
</style>
