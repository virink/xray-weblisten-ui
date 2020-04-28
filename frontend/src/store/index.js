import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
    state: {
        count: 0
    },

    mutations: {
        increment(state, payload) {
            state.count += payload.step
        },
        decrement: state => state.count--,
    }
})