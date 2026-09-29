<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import EventFormDialog from '@/components/EventFormDialog.vue'
import { getUpcomingEvents } from '@/services/events'
import { useAuthStore } from '@/stores/auth'
import type { Event } from '@/types/event'

const router = useRouter()
const authStore = useAuthStore()

const events = ref<Event[]>([])
const loading = ref(true)
const error = ref<string | null>(null)
const showCreateDialog = ref(false)

async function loadEvents() {
  loading.value = true
  error.value = null
  try {
    events.value = await getUpcomingEvents()
  } catch {
    error.value = 'Failed to load upcoming events. Is the backend running?'
  } finally {
    loading.value = false
  }
}

onMounted(loadEvents)

function onEventSaved() {
  showCreateDialog.value = false
  loadEvents()
}

function logout() {
  authStore.logout()
  router.push({ name: 'login' })
}

const dateFormatter = new Intl.DateTimeFormat('en-GB', {
  dateStyle: 'full',
  timeStyle: 'short',
})

function formatDate(timestamp: string): string {
  return dateFormatter.format(new Date(timestamp))
}
</script>

<template>
  <section class="events">
    <div class="events__header">
      <h1>Upcoming Events</h1>
      <div class="events__header-actions">
        <button class="events__create" type="button" @click="showCreateDialog = true">
          + Create Event
        </button>
        <button class="events__logout" type="button" @click="logout">Log Out</button>
      </div>
    </div>

    <EventFormDialog
      v-if="showCreateDialog"
      @saved="onEventSaved"
      @close="showCreateDialog = false"
    />

    <p v-if="loading" class="status">Loading events…</p>

    <p v-else-if="error" class="status status--error">{{ error }}</p>

    <p v-else-if="events.length === 0" class="status">No upcoming events. Check back soon!</p>

    <ul v-else class="events__list">
      <li v-for="event in events" :key="event.id" class="event-card">
        <div class="event-card__header">
          <h2>
            <RouterLink class="event-card__link" :to="{ name: 'event-details', params: { id: event.id } }">
              {{ event.name }}
            </RouterLink>
          </h2>
          <span class="event-card__when">{{ formatDate(event.startTimestamp) }}</span>
        </div>
        <p class="event-card__description">{{ event.description }}</p>
        <div class="event-card__meta">
          <span class="event-card__tag">📍 {{ event.location }}</span>
          <span class="event-card__tag">👥 {{ event.capacity }} places</span>
          <span class="event-card__tag">⏱ {{ event.duration }}h</span>
        </div>
      </li>
    </ul>
  </section>
</template>

<style scoped>
.events {
  max-width: 46rem;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.events__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
  flex-wrap: wrap;
}

.events__header h1 {
  margin: 0;
}

.events__header-actions {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.events__create {
  border: none;
  border-radius: 0.4rem;
  padding: 0.55rem 1.1rem;
  background: #0f766e;
  color: #fff;
  font: inherit;
  font-weight: 600;
  cursor: pointer;
}

.events__create:hover {
  background: #115e59;
}

.events__logout {
  border: 1px solid #cbd5e1;
  border-radius: 0.4rem;
  padding: 0.55rem 1.1rem;
  background: #fff;
  color: #334155;
  font: inherit;
  font-weight: 600;
  cursor: pointer;
}

.events__logout:hover {
  background: #f1f5f9;
}

.events__list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.status {
  padding: 2rem;
  text-align: center;
  color: #64748b;
}

.status--error {
  color: #dc2626;
}

.event-card {
  border: 1px solid #e2e8f0;
  border-radius: 0.75rem;
  padding: 1.25rem;
  background: #fff;
  box-shadow: 0 1px 2px rgb(0 0 0 / 5%);
}

.event-card__header {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  gap: 1rem;
  flex-wrap: wrap;
}

.event-card__header h2 {
  margin: 0;
  font-size: 1.25rem;
}

.event-card__link {
  color: inherit;
  text-decoration: none;
}

.event-card__link:hover {
  color: #0f766e;
  text-decoration: underline;
}

.event-card__when {
  color: #0f766e;
  font-weight: 600;
  white-space: nowrap;
}

.event-card__description {
  margin: 0.5rem 0 0.75rem;
  color: #475569;
}

.event-card__meta {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.event-card__tag {
  background: #f1f5f9;
  border-radius: 999px;
  padding: 0.25rem 0.75rem;
  font-size: 0.85rem;
  color: #334155;
}
</style>
