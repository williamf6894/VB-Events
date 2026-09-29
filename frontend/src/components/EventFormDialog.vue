<script setup lang="ts">
import { ref } from 'vue'
import { createEvent, updateEvent } from '@/services/events'
import type { Event } from '@/types/event'

const props = defineProps<{ event?: Event }>()

const emit = defineEmits<{ saved: []; close: [] }>()

const isEditing = props.event !== undefined

function toDatetimeLocal(iso: string): string {
  const d = new Date(iso)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`
}

const form = ref({
  name: props.event?.name ?? '',
  description: props.event?.description ?? '',
  location: props.event?.location ?? '',
  capacity: props.event?.capacity ?? 20,
  duration: props.event?.duration ?? 2,
  startTimestamp: props.event ? toDatetimeLocal(props.event.startTimestamp) : '',
})

const submitting = ref(false)
const error = ref<string | null>(null)

async function submit() {
  if (!form.value.name.trim()) {
    error.value = 'Please give the event a name.'
    return
  }
  if (!form.value.startTimestamp) {
    error.value = 'Please pick a start date and time.'
    return
  }

  submitting.value = true
  error.value = null

  const payload = {
    name: form.value.name.trim(),
    description: form.value.description.trim(),
    location: form.value.location.trim(),
    capacity: Number(form.value.capacity),
    duration: Number(form.value.duration),
    startTimestamp: new Date(form.value.startTimestamp).toISOString(),
  }

  try {
    if (isEditing && props.event) {
      await updateEvent(props.event.id, payload)
    } else {
      await createEvent(payload)
    }
    emit('saved')
  } catch {
    error.value = isEditing ? 'Failed to update event. Please try again.' : 'Failed to create event. Please try again.'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="dialog__backdrop" @click.self="emit('close')">
    <div class="dialog" role="dialog" aria-modal="true" aria-labelledby="dialog-title">
      <h2 id="dialog-title" class="dialog__title">{{ isEditing ? 'Edit Event' : 'Create Event' }}</h2>

      <form class="dialog__form" @submit.prevent="submit">
        <label class="field">
          <span class="field__label">Name *</span>
          <input v-model="form.name" type="text" class="field__input" required />
        </label>

        <label class="field">
          <span class="field__label">Description</span>
          <textarea v-model="form.description" rows="3" class="field__input"></textarea>
        </label>

        <label class="field">
          <span class="field__label">Location</span>
          <input v-model="form.location" type="text" class="field__input" />
        </label>

        <div class="dialog__row">
          <label class="field">
            <span class="field__label">Capacity</span>
            <input
              v-model.number="form.capacity"
              type="number"
              min="1"
              step="1"
              class="field__input"
            />
          </label>

          <label class="field">
            <span class="field__label">Duration (hours)</span>
            <input
              v-model.number="form.duration"
              type="number"
              min="0.5"
              step="0.5"
              class="field__input"
            />
          </label>
        </div>

        <label class="field">
          <span class="field__label">Starts *</span>
          <input v-model="form.startTimestamp" type="datetime-local" class="field__input" required />
        </label>

        <p v-if="error" class="dialog__error">{{ error }}</p>

        <div class="dialog__actions">
          <button type="button" class="btn btn--secondary" @click="emit('close')">Cancel</button>
          <button type="submit" class="btn btn--primary" :disabled="submitting">
            {{ submitting ? (isEditing ? 'Saving…' : 'Creating…') : isEditing ? 'Save Changes' : 'Create Event' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>
.dialog__backdrop {
  position: fixed;
  inset: 0;
  background: rgb(15 23 42 / 50%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem;
  z-index: 10;
}

.dialog {
  background: #fff;
  border-radius: 0.75rem;
  padding: 1.5rem;
  width: 100%;
  max-width: 28rem;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 10px 30px rgb(0 0 0 / 20%);
}

.dialog__title {
  margin: 0 0 1.25rem;
}

.dialog__form {
  display: flex;
  flex-direction: column;
  gap: 0.9rem;
}

.dialog__row {
  display: grid;
  grid-template-columns: 1fr 1fr;
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

.dialog__error {
  margin: 0;
  color: #dc2626;
  font-size: 0.9rem;
}

.dialog__actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
  margin-top: 0.5rem;
}

.btn {
  border: none;
  border-radius: 0.4rem;
  padding: 0.55rem 1.1rem;
  font: inherit;
  font-weight: 600;
  cursor: pointer;
}

.btn--primary {
  background: #0f766e;
  color: #fff;
}

.btn--primary:hover:not(:disabled) {
  background: #115e59;
}

.btn--primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn--secondary {
  background: #f1f5f9;
  color: #334155;
}

.btn--secondary:hover {
  background: #e2e8f0;
}
</style>
