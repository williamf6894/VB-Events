<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const route = useRoute()
const router = useRouter()

const mode = ref<'login' | 'register'>('login')
const isRegister = computed(() => mode.value === 'register')

const name = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const submitting = ref(false)
const error = ref<string | null>(null)

function validate(): string | null {
  if (isRegister.value && !name.value.trim()) {
    return 'Please enter your name.'
  }
  if (!email.value.trim()) {
    return 'Please enter your email address.'
  }
  if (password.value.length < 8) {
    return 'Password must be at least 8 characters.'
  }
  if (isRegister.value && password.value !== confirmPassword.value) {
    return 'Passwords do not match.'
  }
  return null
}

async function submit() {
  error.value = validate()
  if (error.value) return

  submitting.value = true
  try {
    if (isRegister.value) {
      await authStore.register(name.value.trim(), email.value.trim(), password.value)
    } else {
      await authStore.login(email.value.trim(), password.value)
    }
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : null
    await router.push(redirect ?? { name: 'home' })
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Something went wrong. Please try again.'
  } finally {
    submitting.value = false
  }
}

function switchMode() {
  mode.value = isRegister.value ? 'login' : 'register'
  error.value = null
}
</script>

<template>
  <section class="auth">
    <div class="auth__card">
      <h1 class="auth__title">{{ isRegister ? 'Create an account' : 'Welcome back' }}</h1>
      <p class="auth__subtitle">
        {{ isRegister ? 'Register to create and manage events.' : 'Log in to see upcoming events.' }}
      </p>

      <form class="auth__form" @submit.prevent="submit">
        <label v-if="isRegister" class="field">
          <span class="field__label">Name</span>
          <input v-model="name" type="text" class="field__input" autocomplete="name" />
        </label>

        <label class="field">
          <span class="field__label">Email</span>
          <input
            v-model="email"
            type="email"
            class="field__input"
            autocomplete="email"
            required
          />
        </label>

        <label class="field">
          <span class="field__label">Password</span>
          <input
            v-model="password"
            type="password"
            class="field__input"
            :autocomplete="isRegister ? 'new-password' : 'current-password'"
            minlength="8"
            required
          />
        </label>

        <label v-if="isRegister" class="field">
          <span class="field__label">Confirm Password</span>
          <input
            v-model="confirmPassword"
            type="password"
            class="field__input"
            autocomplete="new-password"
            minlength="8"
            required
          />
        </label>

        <p v-if="error" class="auth__error">{{ error }}</p>

        <button type="submit" class="auth__submit" :disabled="submitting">
          {{ submitting ? 'Please wait…' : isRegister ? 'Register' : 'Log In' }}
        </button>
      </form>

      <p class="auth__switch">
        {{ isRegister ? 'Already have an account?' : "Don't have an account?" }}
        <button type="button" class="auth__switch-button" @click="switchMode">
          {{ isRegister ? 'Log in' : 'Register' }}
        </button>
      </p>
    </div>
  </section>
</template>

<style scoped>
.auth {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 70vh;
  padding: 2rem 1rem;
}

.auth__card {
  width: 100%;
  max-width: 24rem;
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 0.75rem;
  padding: 2rem;
  box-shadow: 0 1px 2px rgb(0 0 0 / 5%);
}

.auth__title {
  margin: 0 0 0.25rem;
}

.auth__subtitle {
  margin: 0 0 1.5rem;
  color: #64748b;
  font-size: 0.95rem;
}

.auth__form {
  display: flex;
  flex-direction: column;
  gap: 0.9rem;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 0.3rem;
}

.field__label {
  font-size: 0.85rem;
  font-weight: 600;
  color: #334155;
}

.field__input {
  border: 1px solid #cbd5e1;
  border-radius: 0.4rem;
  padding: 0.5rem 0.65rem;
  font: inherit;
}

.field__input:focus {
  outline: 2px solid #0f766e;
  outline-offset: 1px;
  border-color: #0f766e;
}

.auth__error {
  margin: 0;
  color: #dc2626;
  font-size: 0.9rem;
}

.auth__submit {
  border: none;
  border-radius: 0.4rem;
  padding: 0.6rem;
  background: #0f766e;
  color: #fff;
  font: inherit;
  font-weight: 600;
  cursor: pointer;
}

.auth__submit:hover:not(:disabled) {
  background: #115e59;
}

.auth__submit:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.auth__switch {
  margin: 1.25rem 0 0;
  text-align: center;
  color: #64748b;
  font-size: 0.9rem;
}

.auth__switch-button {
  border: none;
  background: none;
  color: #0f766e;
  font: inherit;
  font-weight: 600;
  cursor: pointer;
  padding: 0;
}

.auth__switch-button:hover {
  text-decoration: underline;
}
</style>
