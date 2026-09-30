<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { getParticipants, inviteParticipant } from '@/services/events'
import type { Event, Participant } from '@/types/event'

const props = defineProps<{ event: Event }>()

const emit = defineEmits<{ invited: []; close: [] }>()

const participants = ref<Participant[]>([])
const loading = ref(true)
const error = ref<string | null>(null)
const invitingId = ref<string | null>(null)

const notGoing = computed(() =>
  participants.value.filter(
    (p) => !props.event.participants?.some((going) => going.id === p.id),
  ),
)

onMounted(async () => {
  try {
    participants.value = await getParticipants()
  } catch {
    error.value = 'Failed to load participants. Is the backend running?'
  } finally {
    loading.value = false
  }
})

async function invite(participant: Participant) {
  if (invitingId.value) return

  invitingId.value = participant.id
  error.value = null

  try {
    await inviteParticipant(props.event.id, participant.id)
    participants.value = participants.value.filter((p) => p.id !== participant.id)
    emit('invited')
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to invite participant.'
  } finally {
    invitingId.value = null
  }
}
</script>

<template>
  <div class="dialog__backdrop" @click.self="emit('close')">
    <div class="dialog" role="dialog" aria-modal="true" aria-labelledby="invite-dialog-title">
      <h2 id="invite-dialog-title" class="dialog__title">Invite Participants</h2>
      <p class="dialog__subtitle">
        Click someone to add them to <strong>{{ event.name }}</strong>
      </p>

      <p v-if="loading" class="dialog__status">Loading participants…</p>

      <template v-else>
        <p v-if="error" class="dialog__error">{{ error }}</p>

        <p v-if="notGoing.length === 0" class="dialog__status">
          Everyone is already registered for this event.
        </p>

        <ul v-else class="invite-list">
          <li v-for="participant in notGoing" :key="participant.id">
            <button
              type="button"
              class="invite-list__item"
              :disabled="invitingId !== null"
              @click="invite(participant)"
            >
              <span class="invite-list__avatar">{{ participant.name.charAt(0) }}</span>
              <span class="invite-list__name">{{ participant.name }}</span>
              <span class="invite-list__email">{{ participant.email }}</span>
              <span class="invite-list__action">
                {{ invitingId === participant.id ? 'Adding…' : '＋ Invite' }}
              </span>
            </button>
          </li>
        </ul>
      </template>

      <div class="dialog__actions">
        <button type="button" class="btn btn--secondary" @click="emit('close')">Close</button>
      </div>
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
  max-width: 26rem;
  max-height: 85vh;
  overflow-y: auto;
  box-shadow: 0 10px 30px rgb(0 0 0 / 20%);
}

.dialog__title {
  margin: 0 0 0.25rem;
}

.dialog__subtitle {
  margin: 0 0 1rem;
  color: #64748b;
  font-size: 0.9rem;
}

.dialog__status {
  padding: 1.5rem 0;
  text-align: center;
  color: #64748b;
}

.dialog__error {
  margin: 0 0 0.75rem;
  padding: 0.6rem 0.9rem;
  border: 1px solid #fecaca;
  border-radius: 0.4rem;
  background: #fef2f2;
  color: #dc2626;
  font-size: 0.9rem;
}

.invite-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}

.invite-list__item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  width: 100%;
  border: 1px solid #e2e8f0;
  border-radius: 0.5rem;
  padding: 0.55rem 0.75rem;
  background: #fff;
  font: inherit;
  cursor: pointer;
  text-align: left;
}

.invite-list__item:hover:not(:disabled) {
  border-color: #0f766e;
  background: #f0fdfa;
}

.invite-list__item:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.invite-list__avatar {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  border-radius: 999px;
  background: #0f766e;
  color: #fff;
  font-weight: 700;
  flex-shrink: 0;
}

.invite-list__name {
  font-weight: 600;
  color: #1e293b;
  white-space: nowrap;
}

.invite-list__email {
  color: #64748b;
  font-size: 0.85rem;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
}

.invite-list__action {
  color: #0f766e;
  font-weight: 600;
  font-size: 0.85rem;
  white-space: nowrap;
}

.dialog__actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 1rem;
}

.btn {
  border: none;
  border-radius: 0.4rem;
  padding: 0.55rem 1.1rem;
  font: inherit;
  font-weight: 600;
  cursor: pointer;
}

.btn--secondary {
  background: #f1f5f9;
  color: #334155;
}

.btn--secondary:hover {
  background: #e2e8f0;
}
</style>
