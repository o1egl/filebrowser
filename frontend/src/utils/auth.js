import store from '@/store'
import router from '@/router'
import { noAuth } from '@/utils/constants'
import { users } from '@/api'
import {fetchURL} from "@/api/utils";

export function updateUserData (userData) {
  store.commit('setUser', userData)
}

export async function validateLogin () {
  try {
    const body = await users.me()
    updateUserData(body)
  } catch (e) {
    console.warn('user is not authorized', e) // eslint-disable-line
    if (noAuth) {
      await login("anonymous", "", "")
    }
  }
}

export async function login (username, password, recaptcha) {
  const data = { "user": username, "passwd":password, "recaptcha": recaptcha }

  const res = await fetchURL(`/auth/local/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  })

  const body = await res.text()

  if (res.status === 200) {
    const userData = JSON.parse(body)
    updateUserData(userData)
  } else {
    throw new Error(body)
  }
}

export async function signup (username, password) {
  const data = { username, password }

  const res = await fetchURL(`/api/signup`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  })

  if (res.status !== 200) {
    throw new Error(res.status)
  }
}

export async function logout () {
  store.commit('setUser', null)
  await fetchURL(`/auth/logout`, {})
  if (noAuth) {
    await login("anonymous", "", "")
    return
  }
  await router.push({path: '/login'})
}
