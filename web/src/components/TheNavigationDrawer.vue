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
        :key="item.title"
        link
        :to="item.route"
      >
        <v-list-item-icon>
          <v-icon>{{ item.icon }}</v-icon>
        </v-list-item-icon>
        <v-list-item-content>
          <v-list-item-title>{{ item.title }}</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </v-list>
  </v-navigation-drawer>
</template>

<script>
export default {
  name: 'TheNavigationDrawer',
  data: () => ({
    items: [
      {
        title: 'Book',
        route: '/',
        icon: 'mdi-plus',
      },
      {
        title: 'Settings',
        route: 'Settings',
        icon: 'mdi-settings',
      },
      {
        title: 'About',
        route: 'About',
        icon: 'mdi-information',
      },
    ],
  }),
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
  },
}
</script>

<style scoped>

</style>
