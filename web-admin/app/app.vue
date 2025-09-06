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
// PowerX Scrum Plugin Admin Interface
useSeoMeta({
  title: 'PowerX Scrum Plugin',
  description: 'Scrum task management for PowerX platform'
})
</script>
