import * as i18n from '@/i18n'
import moment from 'moment'

const mutations = {
  closeHovers: state => {
    state.show = null
    state.showMessage = null
  },
  toggleShell: (state) => {
    state.showShell = !state.showShell
  },
  showHover: (state, value) => {
    if (typeof value !== 'object') {
      state.show = value
      return
    }

    state.show = value.prompt
    state.showMessage = value.message
    state.showConfirm = value.confirm
  },
  showError: (state, value) => {
    state.show = 'error'
    state.showMessage = value
  },
  showSuccess: (state, value) => {
    state.show = 'success'
    state.showMessage = value
  },
  setLoading: (state, value) => { state.loading = value },
  setReload: (state, value) => { state.reload = value },
  setUser: (state, value) => {
    if (value === null) {
      state.user = null
      return
    }

    let locale = value.attrs.locale

    if (locale === '') {
      locale = i18n.detectLocale()
    }

    moment.locale(locale)
    i18n.default.locale = locale
    state.user = value
  },
  setSorting: (state, value) => {
    state.sorting = value
    localStorage.setItem('sorting', JSON.stringify(value))
  },
  setJWT: (state, value) => (state.jwt = value),
  multiple: (state, value) => (state.multiple = value),
  addSelected: (state, value) => (state.selected.push(value)),
  addPlugin: (state, value) => {
    state.plugins.push(value)
  },
  removeSelected: (state, value) => {
    let i = state.selected.indexOf(value)
    if (i === -1) return
    state.selected.splice(i, 1)
  },
  resetSelected: (state) => {
    state.selected = []
  },
  updateUser: (state, value) => {
    if (typeof value !== 'object') return

    for (let field in value) {
      if (field === 'attrs') {
        for (let attr in value[field]) {
          if (attr === 'locale') {
            moment.locale(value[field])
            i18n.default.locale = value[field]
          }
          state.user[field][attr] = value[field][attr]
        }
      } else {
        state.user[field] = value[field]
      }
    }
  },
  updateRequest: (state, value) => {
    state.oldReq = state.req
    state.req = value
    state.req.items && state.req.items.sort((a, b) => {
      const sortOrder = (state.sorting.asc === true)? -1 : 1;
      const result = (a[state.sorting.by] < b[state.sorting.by]) ? -1 : (a[state.sorting.by] > b[state.sorting.by]) ? 1 : 0;
      return sortOrder*result
    })
  },
  sortRequestItems: (state, sorting) => {
    state.req.items && state.req.items.sort((a, b) => {
      const sortOrder = (sorting.asc === true)? -1 : 1;
      const result = (a[sorting.by] < b[sorting.by]) ? -1 : (a[sorting.by] > b[sorting.by]) ? 1 : 0;
      return sortOrder*result
    })
  },
  updateClipboard: (state, value) => {
    state.clipboard.key = value.key
    state.clipboard.items = value.items
    state.clipboard.path = value.path
  },
  resetClipboard: (state) => {
    state.clipboard.key = ''
    state.clipboard.items = []
  },
  setPreviewMode(state, value) {
    state.previewMode = value
  }
}

export default mutations
