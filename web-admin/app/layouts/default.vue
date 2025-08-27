<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <!-- 顶部导航栏 - 根据环境变量控制显示 -->
    <UContainer v-if="showNavigation" class="max-w-none">
      <AppNavbar />
    </UContainer>

    <div class="flex">
      <!-- 左侧边栏 - 根据环境变量控制显示 -->
      <AppSidebar v-if="showNavigation" />

      <!-- 主内容区 -->
      <main :class="mainContentClass">
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup>
// 获取运行时配置
const runtimeConfig = useRuntimeConfig();

// 控制导航显示的环境变量
const showNavigation = computed(() => {
  // 优先检查环境变量 NUXT_PUBLIC_SHOW_NAVIGATION
  const envShowNav = runtimeConfig.public.showNavigation;

  // 如果没有设置环境变量，开发环境默认显示，生产环境默认隐藏
  if (envShowNav !== undefined) {
    return envShowNav === "true" || envShowNav === true;
  }

  // 默认：开发环境显示，生产环境隐藏
  return process.dev;
});

// 主内容区样式
const mainContentClass = computed(() => {
  return showNavigation.value ? "flex-1 p-6" : "w-full p-6";
});
</script>
