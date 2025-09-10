<template>
  <UModal
    :close="{ onClick: () => emit('close', false) }"
    :title="options.title || $t('common.confirmation')"
  >
    <div class="py-2">
      <p class="text-sm text-gray-600 dark:text-gray-400">
        {{ options.message }}
      </p>
    </div>

    <template #footer>
      <div class="flex justify-end gap-3">
        <UButton variant="outline" @click="emit('close', false)">
          {{ options.cancelText || $t("common.cancel") }}
        </UButton>
        <UButton
          :color="options.type === 'error' ? 'red' : 'primary'"
          @click="emit('close', true)"
        >
          {{ options.confirmText || $t("common.confirm") }}
        </UButton>
      </div>
    </template>
  </UModal>
</template>

<script setup>
const { t } = useI18n();

// 定义props
const props = defineProps({
  options: {
    type: Object,
    default: () => ({}),
  },
});

// 定义emits
const emit = defineEmits(["close"]);
</script>
