<template>
  <teleport to="body">
    <dialog
      v-bind="$attrs"
      ref="dialogRef"
      class="modal common-modal"
      :aria-label="title || 'Dialog'"
      :aria-labelledby="title ? titleId : undefined"
      @cancel.prevent="requestClose"
    >
      <section class="modal-box common-modal__box" :style="panelStyle">
        <header
          class="common-modal__header"
          :class="{ 'common-modal__header--draggable': draggable }"
          @mousedown="$emit('header-mousedown', $event)"
        >
          <slot name="header">
            <h2 v-if="title" :id="titleId" class="common-modal__title">{{ title }}</h2>
          </slot>
          <button
            v-if="showClose"
            type="button"
            class="btn btn-sm btn-circle btn-ghost common-modal__close"
            aria-label="Close"
            @mousedown.stop
            @click="requestClose"
          >
            &times;
          </button>
        </header>

        <div class="common-modal__body">
          <slot />
        </div>

        <slot name="panel-extra" />
      </section>

      <form
        v-if="closeOnBackdrop"
        method="dialog"
        class="modal-backdrop common-modal__backdrop"
        @submit.prevent="requestClose"
      >
        <button type="submit" aria-label="Close dialog">close</button>
      </form>
      <div v-else class="modal-backdrop common-modal__backdrop" aria-hidden="true"></div>
    </dialog>
  </teleport>
</template>

<script>
export default {
  inheritAttrs: false
};
</script>

<script setup>
import { onMounted, ref, watch } from 'vue';

let modalId = 0;

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: ''
  },
  closeOnBackdrop: {
    type: Boolean,
    default: true
  },
  showClose: {
    type: Boolean,
    default: true
  },
  draggable: {
    type: Boolean,
    default: false
  },
  panelStyle: {
    type: [Object, Array, String],
    default: undefined
  }
});

const emit = defineEmits(['update:visible', 'close', 'header-mousedown']);
const dialogRef = ref(null);
const titleId = `app-modal-title-${++modalId}`;

const syncVisibility = () => {
  if (!dialogRef.value) return;

  if (props.visible && !dialogRef.value.open) {
    dialogRef.value.showModal();
  } else if (!props.visible && dialogRef.value.open) {
    dialogRef.value.close();
  }
};

const requestClose = () => {
  if (dialogRef.value?.open) {
    dialogRef.value.close();
  }

  emit('update:visible', false);
  emit('close');
};

watch(() => props.visible, syncVisibility);
onMounted(syncVisibility);
</script>

<style scoped>
.common-modal {
  width: 100vw;
  height: 100vh;
  max-width: none;
  max-height: none;
  margin: 0;
  padding: 0;
  border: 0;
  background: transparent;
  overflow: hidden;
  z-index: 10000;
}

.common-modal[open] {
  display: grid;
  place-items: center;
}

.common-modal::backdrop {
  background-color: rgba(15, 23, 42, 0.58);
  backdrop-filter: blur(3px);
}

.common-modal__box {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  width: min(32rem, calc(100vw - 2rem));
  max-width: calc(100vw - 2rem);
  max-height: calc(100vh - 2rem);
  padding: 0;
  overflow: hidden;
  border-radius: 1rem;
  background-color: #fff;
  box-shadow: 0 24px 64px rgba(15, 23, 42, 0.3);
}

.common-modal__header {
  display: flex;
  flex: none;
  align-items: center;
  justify-content: space-between;
  min-height: 3rem;
  padding: var(--modal-header-padding, 0.65rem) 0.7rem;
  padding-left: calc(var(--modal-header-padding, 0.65rem) + var(--modal-title-offset, 0.4rem));
  border-bottom: 1px solid rgba(148, 163, 184, 0.24);
  background-color: #f8fafc;
  user-select: none;
}

.common-modal__header--draggable {
  cursor: move;
}

.common-modal__title {
  margin: 0;
  color: #0f172a;
  font-size: var(--modal-title-size, 1rem);
  font-weight: 600;
  line-height: 1;
}

.common-modal__close {
  display: inline-flex;
  width: var(--modal-close-size, 2rem);
  height: var(--modal-close-size, 2rem);
  min-height: 0;
  align-items: center;
  justify-content: center;
  padding: 0;
  border: 0;
  border-radius: 999px;
  background-color: transparent;
  color: #475569;
  font-size: 1.3rem;
  line-height: 1;
}

.common-modal__close:hover {
  background-color: #e2e8f0;
  color: #0f172a;
}

.common-modal__body {
  flex: 1;
  min-height: 0;
  overflow: auto;
}

.common-modal__backdrop {
  position: absolute;
  inset: 0;
  z-index: 0;
}

.common-modal__backdrop button {
  width: 100%;
  height: 100%;
  border: 0;
  background-color: transparent;
  color: transparent;
  cursor: default;
}
</style>
