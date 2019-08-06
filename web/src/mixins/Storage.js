let storage = window.localStorage

export default {
  data: () => ({
    localStorage: storage,
  }),
  computed: {
    storageSupported () {
      return typeof(this.localStorage) !== "undefined"
    },
  },
}
