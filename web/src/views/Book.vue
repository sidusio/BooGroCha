<template>
  <v-layout
    column
    fill-height
  >
    <v-container>
      <span v-if="date !== null">Date: {{formatDate(date)}} </span>
      <span v-if="fromTime !== null">| From: {{fromTime}} </span>
      <span v-if="toTime !== null">| To: {{toTime}}</span>
    </v-container>
    <DateSelector
      v-if="state === states.selectDate"
      @picked="dateSelected"
    ></DateSelector>
  </v-layout>
</template>

<script>
import moment from 'moment'
import DateSelector from '../components/DateSelector'

var states = Object.freeze({
  selectDate: 1,
  selectFromTime: 2,
  selectToTime: 3,
  selectRoom: 4,
  loading: 5,
})

export default {
  name: 'Book',
  components: {DateSelector},
  data: () => ({
    states: states,
    date: null,
    fromTime: null,
    toTime: null,
  }),
  computed: {
    state() {
      if (this.date === null) return states.selectDate
      else return states.loading
    }
  },
  methods: {
    dateSelected(date) {
      this.date = date
    },
    formatDate(date) {
      if (date === null) {
        return ''
      }
      return moment(date).format("dd D/M")
    }
  }
}
</script>

<style scoped>
</style>
