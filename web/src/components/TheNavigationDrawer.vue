<template>
  <v-navigation-drawer
    v-model="open"
    app
    right
  >
    <v-list
      dense
      nav
    >
      <v-list-item
        v-for="item in items"
        :key="item.name"
        link
        :to="item"
        :disabled="item.meta.requiresCredentials && !hasCredentials"
      >
        <v-list-item-icon>
          <v-icon>{{ item.meta.icon }}</v-icon>
        </v-list-item-icon>
        <v-list-item-content>
          <v-list-item-title>{{ item.name }}</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </v-list>
  </v-navigation-drawer>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  name: 'TheNavigationDrawer',
  data: () => ({
    items: [],
  }),
  created () {
    this.$router.options.routes.forEach(route => {
      this.items.push(route)
    })
    this.items.sort((a, b) => {
      return a.meta.priority - b.meta.priority
    })
  },
  props: {
    drawer: Boolean,
  },
  computed: {
    open: {
      get () {
        return this.drawer
      },
      set (v) {
        this.$emit('update:drawer', v)
      },
    },
    ...mapGetters([
      'hasCredentials',
    ]),
  },
}
</script>

<style scoped>

</style>
