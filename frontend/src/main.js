import { sync } from 'vuex-router-sync'
import store from '@/store'
import router from '@/router'
import i18n from '@/i18n'
import Vue from '@/utils/vue'
import { loginPage } from '@/utils/constants'
//import { recaptcha, loginPage } from '@/utils/constants'
import { login, validateLogin } from '@/utils/auth'
import App from '@/App'

sync(store, router)

async function start () {
  if (loginPage) {
    await validateLogin()
  } else {
    await login('anonymous', '', '')
  }

/*
  if (recaptcha) {
    await new Promise (resolve => {
      const check = () => {
        if (typeof window.grecaptcha === 'undefined') {
          setTimeout(check, 100)
        } else {
          resolve()
        }
      }

      check()
    })
  }*/

  const sort = JSON.parse(localStorage.getItem('sorting'))
  if (sort) {
    store.commit('setSorting', sort)
  }

  new Vue({
    el: '#app',
    store,
    router,
    i18n,
    template: '<App/>',
    components: { App }
  })
}

start()
