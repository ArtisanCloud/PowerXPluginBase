<template>
  <aside
    class="w-64 min-w-64 max-w-64 bg-white dark:bg-gray-900 border-r border-gray-200 dark:border-gray-800 min-h-screen flex-shrink-0"
  >
    <nav class="p-4 space-y-6">
      <!-- 仪表盘入口 -->
      <div>
        <UButton
          to="/dashboard"
          variant="ghost"
          color="neutral"
          class="w-full justify-start"
          :class="{
            'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
              isExactActive('/dashboard'),
          }"
        >
          <UIcon name="i-heroicons-home" class="w-4 h-4 mr-3" />
          {{ t('navigation.dashboard') }}
        </UButton>
      </div>

      <!-- 笔记模块 -->
      <div>
        <div class="px-3 py-2 text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">
          {{ t('navigation.contentManagement') }}
        </div>
        <div class="space-y-1">
          <div>
            <UButton
              variant="ghost"
              color="neutral"
              class="w-full justify-start"
              @click="toggleNotesMenu"
              :class="{
                'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                  isGroupActive(['/notes']),
              }"
            >
              <UIcon name="i-heroicons-document-text" class="w-4 h-4 mr-3" />
              {{ t('navigation.notes') }}
              <UIcon
                :name="
                  showNotesMenu
                    ? 'i-heroicons-chevron-down'
                    : 'i-heroicons-chevron-right'
                "
                class="w-4 h-4 ml-auto"
              />
            </UButton>

            <div v-show="showNotesMenu" class="ml-6 mt-1 space-y-1">
              <UButton
                to="/notes"
                variant="ghost"
                color="neutral"
                size="sm"
                class="w-full justify-start text-sm"
                :class="{
                  'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                    isExactActive('/notes'),
                }"
              >
                <UIcon name="i-heroicons-document-text" class="w-3 h-3 mr-2" />
                {{ t('notes.overview') }}
              </UButton>

          <UButton
            to="/notes/active"
            variant="ghost"
            color="neutral"
            size="sm"
            class="w-full justify-start text-sm"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                isExactActive('/notes/active'),
            }"
          >
            <UIcon name="i-heroicons-play-circle" class="w-3 h-3 mr-2" />
            {{ t('navigation.activenote') }}
          </UButton>
            </div>
          </div>
        </div>
      </div>

      <!-- 团队协作 -->
      <div>
        <div class="px-3 py-2 text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">
          {{ t('navigation.teamCollaboration') }}
        </div>
        <div class="space-y-1">
          <div>
            <UButton
              variant="ghost"
              color="neutral"
              class="w-full justify-start"
              @click="toggleTeamMenu"
              :class="{
                'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                  isGroupActive(['/team']),
              }"
            >
              <UIcon name="i-heroicons-users" class="w-4 h-4 mr-3" />
              {{ t('navigation.teamCollaboration') }}
              <UIcon
                :name="
                  showTeamMenu
                    ? 'i-heroicons-chevron-down'
                    : 'i-heroicons-chevron-right'
                "
                class="w-4 h-4 ml-auto"
              />
            </UButton>

            <div v-show="showTeamMenu" class="ml-6 mt-1 space-y-1">
              <UButton
                to="/team/management"
                variant="ghost"
                color="neutral"
                size="sm"
                class="w-full justify-start text-sm"
                :class="{
                  'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                    isExactActive('/team/management'),
                }"
              >
                <UIcon name="i-heroicons-building-office" class="w-3 h-3 mr-2" />
                {{ t('navigation.teamManagement') }}
              </UButton>

              <UButton
                to="/team/members"
                variant="ghost"
                color="neutral"
                size="sm"
                class="w-full justify-start text-sm"
                :class="{
                  'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                    isExactActive('/team/members'),
                }"
              >
                <UIcon name="i-heroicons-user-group" class="w-3 h-3 mr-2" />
                {{ t('navigation.memberManagement') }}
              </UButton>
            </div>
          </div>
        </div>
      </div>

      <!-- 报表模块 -->
      <div>
        <div class="px-3 py-2 text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">
          {{ t('navigation.reports') }}
        </div>
        <div class="space-y-1">
          <div>
            <UButton
              variant="ghost"
              color="neutral"
              class="w-full justify-start"
              @click="toggleReportsMenu"
              :class="{
                'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                  isGroupActive(['/reports']),
              }"
            >
              <UIcon name="i-heroicons-chart-bar" class="w-4 h-4 mr-3" />
              {{ t('navigation.reports') }}
              <UIcon
                :name="
                  showReportsMenu
                    ? 'i-heroicons-chevron-down'
                    : 'i-heroicons-chevron-right'
                "
                class="w-4 h-4 ml-auto"
              />
            </UButton>

            <div v-show="showReportsMenu" class="ml-6 mt-1 space-y-1">
              <UButton
                to="/reports"
                variant="ghost"
                color="neutral"
                size="sm"
                class="w-full justify-start text-sm"
                :class="{
                  'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                    isExactActive('/reports'),
                }"
              >
                <UIcon name="i-heroicons-clipboard-document-list" class="w-3 h-3 mr-2" />
                {{ t('reports.overview') }}
              </UButton>

              <UButton
                to="/reports/daily"
                variant="ghost"
                color="neutral"
                size="sm"
                class="w-full justify-start text-sm"
                :class="{
                  'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                    isExactActive('/reports/daily'),
                }"
              >
                <UIcon name="i-heroicons-calendar-days" class="w-3 h-3 mr-2" />
                {{ t('navigation.reportsDaily') }}
              </UButton>

              <UButton
                to="/reports/monthly"
                variant="ghost"
                color="neutral"
                size="sm"
                class="w-full justify-start text-sm"
                :class="{
                  'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                    isExactActive('/reports/monthly'),
                }"
              >
                <UIcon name="i-heroicons-calendar" class="w-3 h-3 mr-2" />
                {{ t('navigation.reportsMonthly') }}
              </UButton>
            </div>
          </div>
        </div>
      </div>

      <!-- 系统设置 -->
      <div>
        <div class="px-3 py-2 text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">
          {{ t('navigation.systemSettings') }}
        </div>
        <div class="space-y-1">
          <UButton
            to="/settings/permissions"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                isExactActive('/settings/permissions'),
            }"
          >
            <UIcon name="i-heroicons-key" class="w-4 h-4 mr-3" />
            {{ t('navigation.permissions') }}
          </UButton>
        </div>
      </div>
    </nav>
  </aside>
</template>

<script setup lang="ts">
const { t } = useI18n();
const route = useRoute();

const showNotesMenu = ref(true);
const showTeamMenu = ref(false);
const showReportsMenu = ref(true);

const normalizePath = (value: string) => {
  if (!value) {
    return '/';
  }
  if (value !== '/' && value.endsWith('/')) {
    return value.replace(/\/+$/, '');
  }
  return value;
};

const isExactActive = (target: string) => {
  const current = normalizePath(route.path);
  const normalizedTarget = normalizePath(target);

  if (normalizedTarget === '/dashboard') {
    return current === '/dashboard' || current === '/';
  }

  return current === normalizedTarget;
};

const isGroupActive = (prefixes: string | string[]) => {
  const current = normalizePath(route.path);
  const list = Array.isArray(prefixes) ? prefixes : [prefixes];

  return list.some((prefix) => {
    const normalizedPrefix = normalizePath(prefix);

    if (normalizedPrefix === '/dashboard') {
      return current === '/dashboard' || current === '/';
    }

    return (
      current === normalizedPrefix || current.startsWith(`${normalizedPrefix}/`)
    );
  });
};

const toggleNotesMenu = () => {
  showNotesMenu.value = !showNotesMenu.value;
};

const toggleTeamMenu = () => {
  showTeamMenu.value = !showTeamMenu.value;
};

const toggleReportsMenu = () => {
  showReportsMenu.value = !showReportsMenu.value;
};

watch(
  () => route.path,
  (newPath) => {
    if (newPath.startsWith('/notes')) {
      showNotesMenu.value = true;
    }
    if (newPath.startsWith('/team')) {
      showTeamMenu.value = true;
    }
    if (newPath.startsWith('/reports')) {
      showReportsMenu.value = true;
    }
  },
  { immediate: true }
);
</script>
