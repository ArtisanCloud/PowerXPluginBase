<template>
  <div class="min-h-screen bg-white dark:bg-gray-900">
    <main class="w-full p-6">
      <slot />
    </main>
  </div>
</template>

<script setup lang="ts">
useHead({
  bodyAttrs: {
    class: "bg-white dark:bg-gray-900",
  },
});
</script>
