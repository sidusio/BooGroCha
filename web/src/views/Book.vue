<template>
  <div
    class="book"
  >
    <v-container
      text-center
      class="pb-0 mb-0"
    >
      <span class="underlined" v-if="date !== null" @click="date = null">
        {{formatDate(date)}} <v-icon>mdi-menu-down</v-icon>
      </span>
    </v-container>
    <v-container
      text-center
      v-if="fromTime !== null || toTime !== null"
    >
      <span class="underlined" v-if="fromTime !== null" @click="fromTime = null">
        {{fromTime}} <v-icon>mdi-menu-down</v-icon>
      </span>
      -
      <span class="underlined ml-1" v-if="toTime !== null" @click="toTime = null">
        {{toTime}} <v-icon>mdi-menu-down</v-icon>
      </span>
    </v-container>
    <DateSelector
      v-if="state === states.selectDate"
      @picked="dateSelected"
    />
    <TimeSelector
      key="from-selector"
      v-if="state === states.selectFromTime"
      :max="toTime"
      @picked="timeSelected"
    />
    <TimeSelector
      key="to-selector"
      v-if="state === states.selectToTime"
      :min="fromTime"
      @picked="timeSelected"
    />
    <RoomSelector
      v-if="state === states.selectRoom"
      :rooms="availableRooms"
      @book="book"
    />
  </div>
</template>

<script>
import moment from 'moment'
import DateSelector from '../components/DateSelector'
import TimeSelector from '../components/TimeSelector'
import RoomSelector from '../components/RoomSelector'

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
    RoomSelector,
  },
  data: () => ({
    states: states,
    date: null,
    fromTime: null,
    toTime: null,
    availableRooms: null,
  }),
  computed: {
    state () {
      if (this.date === null) return states.selectDate
      else if (this.fromTime === null) return states.selectFromTime
      else if (this.toTime === null) return states.selectToTime
      else if (this.availableRooms !== null) return states.selectRoom
      else return states.loading
    },
  },
  methods: {
    dateSelected (date) {
      this.date = date
    },
    timeSelected (time) {
      if (this.state === states.selectFromTime) {
        this.fromTime = time
      } else if (this.state === states.selectToTime) {
        this.toTime = time
      }
    },
    formatDate (date) {
      if (date === null) {
        return ''
      }
      return moment(date).format('dddd D/M')
    },
    book (room) {
      console.log(room)
    },
    getAvailableRooms () {
      this.availableRooms = null
      if (this.date === null) return
      if (this.fromTime === null) return
      if (this.toTime === null) return
      this.axios.get('booking/available', {
        params: {
          from: moment(moment(this.date).format('YYYY-M-DT') + this.fromTime).format('YYYY-M-DTH:mm'),
          to: moment(moment(this.date).format('YYYY-M-DT') + this.toTime).format('YYYY-M-DTH:mm'),
        },
      }).then((response) => {
        if (response.status === 200) {
          this.availableRooms = response.data.Rooms
        }
      })
    },
  },
  watch: {
    date () {
      this.getAvailableRooms()
    },
    fromTime () {
      this.getAvailableRooms()
    },
    toTime () {
      this.getAvailableRooms()
    },
  },
}
</script>

<style scoped>
  .book {
    height: 100%;
    display: flex;
    flex-direction: column;
  }
  .underlined {
    border-bottom: 1px solid rgba(0, 0, 0, 0.4);
  }
</style>
