<template>
  <UModal v-model="isOpen">
    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <h3 class="text-base font-semibold leading-6 text-gray-900 dark:text-white">
            {{ options.title || $t('common.confirmation') }}
          </h3>
          <UButton
            color="gray"
            variant="ghost"
            icon="i-heroicons-x-mark-20-solid"
            class="-my-1"
            @click="cancel"
          />
        </div>
      </template>

      <div class="py-2">
        <p class="text-sm text-gray-600 dark:text-gray-400">
          {{ options.message }}
        </p>
      </div>

      <template #footer>
        <div class="flex justify-end gap-3">
          <UButton
            variant="outline"
            @click="cancel"
          >
            {{ options.cancelText || $t('common.cancel') }}
          </UButton>
          <UButton
            :color="options.type === 'error' ? 'red' : 'primary'"
            @click="confirm"
          >
            {{ options.confirmText || $t('common.confirm') }}
          </UButton>
        </div>
      </template>
    </UCard>
  </UModal>
</template>

<script setup>
import useConfirm from '~/composables/useConfirm';

const { isOpen, options, confirm, cancel } = useConfirm();
</script>