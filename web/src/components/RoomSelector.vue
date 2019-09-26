<template>
  <v-container
    class="room-selector"
    ma-0
    pa-0
  >
    <v-list
      :dense="rooms.length > 15"
      width="100%"
    >
      <v-list-item v-for="room in rooms" :key="room.provider+room.id" @click="book(room)">
        <v-list-item-content>
          <v-list-item-title>{{room.provider}} > {{room.id}}</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </v-list>
  </v-container>
</template>

<script>

export default {
  name: 'DateSelector',
  props: {
    rooms: Array,
  },
  methods: {
    book (room) {
      console.log(room)
      this.$emit('book', room)
    },
  },
}
</script>

<style scoped lang="scss">
  .room-selector {
    overflow-y: auto;
  }
</style>
