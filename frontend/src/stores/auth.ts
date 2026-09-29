import { defineStore } from 'pinia'
import { loginRequest, meRequest, registerRequest } from '@/services/auth'
import { readToken } from '@/services/api'
import type { Participant } from '@/types/event'

const TOKEN_KEY = 'vb-events-token'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: readToken() ?? '',
    user: null as Participant | null,
  }),
  getters: {
    isAuthenticated: (state): boolean => state.token !== '',
  },
  actions: {
    async fetchUser() {
      if (!this.token) return
      try {
        this.user = await meRequest()
      } catch {
        this.logout()
      }
    },
    async login(email: string, password: string) {
      this.token = await loginRequest(email, password)
      localStorage.setItem(TOKEN_KEY, this.token)
      await this.fetchUser()
    },
    async register(name: string, email: string, password: string) {
      await registerRequest(name, email, password)
      await this.login(email, password)
    },
    logout() {
      this.token = ''
      this.user = null
      localStorage.removeItem(TOKEN_KEY)
    },
  },
})
