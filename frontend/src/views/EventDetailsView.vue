<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import EventFormDialog from '@/components/EventFormDialog.vue'
import { getEvent } from '@/services/events'
import type { Event } from '@/types/event'

const route = useRoute()

const event = ref<Event | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)
const notFound = ref(false)
const showEditDialog = ref(false)

async function loadEvent() {
  loading.value = true
  error.value = null
  try {
    event.value = await getEvent(String(route.params.id))
  } catch (err) {
    if (err instanceof Error && err.message.includes('404')) {
      notFound.value = true
    } else {
      error.value = 'Failed to load event. Is the backend running?'
    }
  } finally {
    loading.value = false
  }
}

onMounted(loadEvent)

function onEventSaved() {
  showEditDialog.value = false
  loadEvent()
}

const spotsAvailable = computed(() =>
  event.value ? event.value.capacity - (event.value.participants?.length ?? 0) : 0,
)

const dateFormatter = new Intl.DateTimeFormat('en-GB', {
  dateStyle: 'full',
  timeStyle: 'short',
})

function formatDate(timestamp: string): string {
  return dateFormatter.format(new Date(timestamp))
}
</script>

<template>
  <section class="event-details">
    <RouterLink class="event-details__back" to="/">← Back to upcoming events</RouterLink>

    <p v-if="loading" class="status">Loading event…</p>

    <p v-else-if="notFound" class="status status--error">Event not found.</p>

    <p v-else-if="error" class="status status--error">{{ error }}</p>

    <template v-else-if="event">
      <EventFormDialog
        v-if="showEditDialog"
        :event="event"
        @saved="onEventSaved"
        @close="showEditDialog = false"
      />

      <header class="event-details__header">
        <h1>{{ event.name }}</h1>
        <div class="event-details__header-actions">
          <span class="event-details__when">{{ formatDate(event.startTimestamp) }}</span>
          <button class="event-details__edit" type="button" @click="showEditDialog = true">
            Edit
          </button>
        </div>
      </header>

      <p class="event-details__description">{{ event.description }}</p>

      <dl class="event-details__meta">
        <div class="event-details__row">
          <dt>Location</dt>
          <dd>📍 {{ event.location }}</dd>
        </div>
        <div class="event-details__row">
          <dt>Spots Available</dt>
          <dd>👥 {{ spotsAvailable }}</dd>
        </div>
        <div class="event-details__row">
          <dt>Duration</dt>
          <dd>⏱ {{ event.duration }} hours</dd>
        </div>
      </dl>

      <h2 class="event-details__participants-title">
        Participants
        <span class="event-details__count">({{ event.participants?.length ?? 0 }})</span>
      </h2>

      <ul v-if="event.participants?.length" class="event-details__participants">
        <li v-for="participant in event.participants" :key="participant.id" class="participant">
          <span class="participant__avatar">{{ participant.name.charAt(0) }}</span>
          <div>
            <span class="participant__name">{{ participant.name }}</span>
          </div>
        </li>
      </ul>

      <p v-else class="status">No participants registered yet.</p>
    </template>
  </section>
</template>

<style scoped>
.event-details {
  max-width: 40rem;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.event-details__back {
  display: inline-block;
  margin-bottom: 1.5rem;
  color: #0f766e;
  text-decoration: none;
  font-weight: 600;
}

.event-details__back:hover {
  text-decoration: underline;
}

.event-details__header {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  gap: 1rem;
  flex-wrap: wrap;
}

.event-details__header h1 {
  margin: 0;
}

.event-details__header-actions {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.event-details__when {
  color: #0f766e;
  font-weight: 600;
  white-space: nowrap;
}

.event-details__edit {
  border: 1px solid #cbd5e1;
  border-radius: 0.4rem;
  padding: 0.45rem 1rem;
  background: #fff;
  color: #334155;
  font: inherit;
  font-weight: 600;
  cursor: pointer;
}

.event-details__edit:hover {
  background: #f1f5f9;
}

.event-details__description {
  margin: 0.75rem 0 1.5rem;
  color: #475569;
  font-size: 1.05rem;
}

.event-details__meta {
  margin: 0 0 2rem;
  border: 1px solid #e2e8f0;
  border-radius: 0.75rem;
  background: #fff;
  overflow: hidden;
}

.event-details__row {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
  padding: 0.75rem 1rem;
}

.event-details__row + .event-details__row {
  border-top: 1px solid #e2e8f0;
}

.event-details__row dt {
  color: #64748b;
  font-weight: 600;
}

.event-details__row dd {
  margin: 0;
}

.event-details__participants-title {
  font-size: 1.1rem;
  margin-bottom: 0.75rem;
}

.event-details__count {
  color: #64748b;
  font-weight: 400;
}

.event-details__participants {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.participant {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  border: 1px solid #e2e8f0;
  border-radius: 0.5rem;
  padding: 0.6rem 0.9rem;
  background: #fff;
}

.participant__avatar {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 2.25rem;
  height: 2.25rem;
  border-radius: 999px;
  background: #0f766e;
  color: #fff;
  font-weight: 700;
}

.participant__name {
  font-weight: 600;
}

.status {
  padding: 2rem;
  text-align: center;
  color: #64748b;
}

.status--error {
  color: #dc2626;
}
</style>
