<template>
  <v-app class="app">
    <TheNavigationDrawer
      :drawer.sync="drawer"
    ></TheNavigationDrawer>
    <v-app-bar app>
      <v-toolbar-title class="headline">
        <span>BooGroCha</span>
      </v-toolbar-title>
      <v-spacer></v-spacer>
      <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
    </v-app-bar>
    <v-content>
      <router-view />
    </v-content>
  </v-app>
</template>

<script>
import TheNavigationDrawer from './components/TheNavigationDrawer'
export default {
  name: 'App',
  components: {TheNavigationDrawer},
  data: () => ({
    drawer: false,
  }),
}
</script>

<style scoped lang="scss">
  .app {
    height: 100vh;
    .v-content {
      height: 100%;
    }
  }
</style>
