<template>
  <div
    class="border-b border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900"
  >
    <div class="flex h-16 items-center justify-between px-6">
      <!-- 左侧品牌 -->
      <div class="flex items-center">
        <h1 class="text-xl font-semibold text-gray-900 dark:text-white">
          {{ $t("common.welcome") }}
        </h1>
      </div>

      <!-- 右侧控制区 -->
      <div class="flex items-center space-x-4">
        <ThemeSelector />
        <LanguageSelector />
        <UAvatar
          src="https://avatars.githubusercontent.com/u/739984?v=4"
          alt="管理员"
          size="sm"
        />
      </div>
    </div>
  </div>
</template>
