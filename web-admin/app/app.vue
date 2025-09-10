<template>
  <UApp>
    <NuxtLoadingIndicator :throttle="150" />
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
    <ConfirmDialog />
  </UApp>
</template>

<script setup>
// PowerX Note Plugin Admin Interface
useSeoMeta({
  title: "PowerX Note Plugin",
  description: "Note note management for PowerX platform",
});
</script>
