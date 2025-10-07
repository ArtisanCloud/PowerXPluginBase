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
// PowerX Base Template Plugin Admin Interface
useSeoMeta({
  title: "PowerX Base Template Plugin",
  description: "Template management showcase for the PowerX platform",
});
</script>
