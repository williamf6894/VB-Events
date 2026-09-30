<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import EventFormDialog from '@/components/EventFormDialog.vue'
import { deleteEvent, getEvent, joinEvent, leaveEvent } from '@/services/events'
import { useAuthStore } from '@/stores/auth'
import type { Event } from '@/types/event'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const event = ref<Event | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)
const notFound = ref(false)
const showEditDialog = ref(false)
const deleting = ref(false)
const joining = ref(false)
const leaving = ref(false)
const joinError = ref<string | null>(null)

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

async function onDelete() {
  if (!event.value || deleting.value) return
  if (!window.confirm(`Delete "${event.value.name}"? This cannot be undone.`)) return

  deleting.value = true
  error.value = null

  try {
    await deleteEvent(event.value.id)
    router.push({ name: 'home' })
  } catch {
    error.value = 'Failed to delete event. Please try again.'
    deleting.value = false
  }
}

const spotsAvailable = computed(() =>
  event.value ? event.value.capacity - (event.value.participants?.length ?? 0) : 0,
)

const isParticipating = computed(
  () => event.value?.participants?.some((p) => p.id === authStore.user?.id) ?? false,
)

const isFull = computed(() => spotsAvailable.value <= 0)

const hasStarted = computed(() =>
  event.value ? new Date(event.value.startTimestamp).getTime() <= Date.now() : false,
)

const joinDisabledReason = computed(() => {
  if (hasStarted.value) return 'Event Started'
  if (isFull.value) return 'Event Full'
  return null
})

async function onParticipate() {
  if (!event.value || joining.value || isParticipating.value || joinDisabledReason.value) return

  joining.value = true
  joinError.value = null

  try {
    await joinEvent(event.value.id)
    await loadEvent()
  } catch (err) {
    joinError.value = err instanceof Error ? err.message : 'Failed to join event.'
  } finally {
    joining.value = false
  }
}

async function onLeave() {
  if (!event.value || leaving.value || !isParticipating.value) return

  leaving.value = true
  joinError.value = null

  try {
    await leaveEvent(event.value.id)
    await loadEvent()
  } catch (err) {
    joinError.value = err instanceof Error ? err.message : 'Failed to leave event.'
  } finally {
    leaving.value = false
  }
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
          <button
            v-if="!isParticipating"
            class="event-details__participate"
            type="button"
            :disabled="joinDisabledReason !== null || joining"
            :title="joinDisabledReason ? `Cannot join: ${joinDisabledReason}` : undefined"
            @click="onParticipate"
          >
            {{ joining ? 'Joining…' : joinDisabledReason ?? 'Participate' }}
          </button>
          <template v-else>
            <span class="event-details__going">✓ You're going</span>
            <button
              class="event-details__leave"
              type="button"
              :disabled="leaving"
              @click="onLeave"
            >
              {{ leaving ? 'Leaving…' : 'Leave' }}
            </button>
          </template>
          <button class="event-details__edit" type="button" @click="showEditDialog = true">
            Edit
          </button>
          <button class="event-details__delete" type="button" :disabled="deleting" @click="onDelete">
            {{ deleting ? 'Deleting…' : 'Delete' }}
          </button>
        </div>
      </header>

      <p class="event-details__description">{{ event.description }}</p>

      <p v-if="joinError" class="event-details__join-error">{{ joinError }}</p>

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

.event-details__participate {
  border: none;
  border-radius: 0.4rem;
  padding: 0.45rem 1rem;
  background: #0f766e;
  color: #fff;
  font: inherit;
  font-weight: 600;
  cursor: pointer;
}

.event-details__participate:hover:not(:disabled) {
  background: #115e59;
}

.event-details__participate:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.event-details__going {
  color: #0f766e;
  font-weight: 600;
  white-space: nowrap;
}

.event-details__leave {
  border: 1px solid #fecaca;
  border-radius: 0.4rem;
  padding: 0.45rem 1rem;
  background: #fff;
  color: #dc2626;
  font: inherit;
  font-weight: 600;
  cursor: pointer;
}

.event-details__leave:hover:not(:disabled) {
  background: #fef2f2;
}

.event-details__leave:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.event-details__join-error {
  margin: 0.75rem 0 0;
  padding: 0.6rem 0.9rem;
  border: 1px solid #fecaca;
  border-radius: 0.4rem;
  background: #fef2f2;
  color: #dc2626;
}

.event-details__edit:hover {
  background: #f1f5f9;
}

.event-details__delete {
  border: 1px solid #fecaca;
  border-radius: 0.4rem;
  padding: 0.45rem 1rem;
  background: #fff;
  color: #dc2626;
  font: inherit;
  font-weight: 600;
  cursor: pointer;
}

.event-details__delete:hover:not(:disabled) {
  background: #fef2f2;
}

.event-details__delete:disabled {
  opacity: 0.6;
  cursor: not-allowed;
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
