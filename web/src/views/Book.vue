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
    <TimeSelector
      key="from-selector"
      v-if="state === states.selectFromTime"
      :max="toTime"
      @picked="timeSelected"
    >
    </TimeSelector>
    <TimeSelector
      key="to-selector"
      v-if="state === states.selectToTime"
      :min="fromTime"
      @picked="timeSelected"
    >
    </TimeSelector>
  </v-layout>
</template>

<script>
import moment from 'moment'
import DateSelector from '../components/DateSelector'
import TimeSelector from '../components/TimeSelector'

let states = Object.freeze({
  selectDate: 1,
  selectFromTime: 2,
  selectToTime: 3,
  selectRoom: 4,
  loading: 5,
})

export default {
  name: 'Book',
  components: {
    TimeSelector,
    DateSelector,
  },
  data: () => ({
    states: states,
    date: null,
    fromTime: null,
    toTime: null,
  }),
  computed: {
    state() {
      if (this.date === null) return states.selectDate
      else if (this.fromTime === null) return states.selectFromTime
      else if (this.toTime === null) return states.selectToTime
      else return states.loading
    }
  },
  methods: {
    dateSelected(date) {
      this.date = date
    },
    timeSelected(time) {
      if (this.state === states.selectFromTime) {
        this.fromTime = time
      } else if (this.state === states.selectToTime) {
        this.toTime = time
      }
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
