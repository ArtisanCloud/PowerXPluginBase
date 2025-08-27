<template>
  <aside
    class="w-64 min-w-64 max-w-64 bg-white dark:bg-gray-900 border-r border-gray-200 dark:border-gray-800 min-h-screen flex-shrink-0"
  >
    <nav class="p-4 space-y-6">
      <!-- 仪表盘 -->
      <div>
        <UButton
          to="/dashboard"
          variant="ghost"
          color="neutral"
          class="w-full justify-start"
          :class="{
            'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
              $route.path === '/dashboard',
          }"
        >
          <UIcon name="i-heroicons-chart-bar" class="w-4 h-4 mr-3" />
          {{ $t("nav.dashboard") }}
        </UButton>
      </div>

      <!-- 客户运营 -->
      <div>
        <div
          class="px-3 py-2 text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider"
        >
          {{ $t("nav.customerOps") }}
        </div>
        <div class="space-y-1">
          <UButton
            to="/customer"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                $route.path === '/customer',
            }"
          >
            <UIcon name="i-heroicons-users" class="w-4 h-4 mr-3" />
            {{ $t("nav.customers") }}
          </UButton>

          <UButton
            to="/customer/members"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                $route.path === '/customer/members',
            }"
          >
            <UIcon name="i-heroicons-user-group" class="w-4 h-4 mr-3" />
            {{ $t("nav.members") }}
          </UButton>

          <!-- 用户分析 - 带子菜单 -->
          <div>
            <UButton
              variant="ghost"
              color="neutral"
              class="w-full justify-start"
              @click="toggleUserAnalytics"
            >
              <UIcon name="i-heroicons-chart-pie" class="w-4 h-4 mr-3" />
              {{ $t("nav.userAnalytics") }}
              <UIcon
                :name="
                  showUserAnalytics
                    ? 'i-heroicons-chevron-down'
                    : 'i-heroicons-chevron-right'
                "
                class="w-4 h-4 ml-auto"
              />
            </UButton>

            <div v-show="showUserAnalytics" class="ml-6 mt-1 space-y-1">
              <UButton
                to="/customer/analytics/segmentation"
                variant="ghost"
                color="neutral"
                size="sm"
                class="w-full justify-start text-sm"
                :class="{
                  'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                    $route.path === '/customer/analytics/segmentation',
                }"
              >
                <UIcon name="i-heroicons-user-group" class="w-3 h-3 mr-2" />
                {{ $t("nav.userSegmentation") }}
              </UButton>

              <UButton
                to="/customer/analytics/cohort"
                variant="ghost"
                color="neutral"
                size="sm"
                class="w-full justify-start text-sm"
                :class="{
                  'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                    $route.path === '/customer/analytics/cohort',
                }"
              >
                <UIcon name="i-heroicons-user-group" class="w-3 h-3 mr-2" />
                {{ $t("nav.cohortAnalysis") }}
              </UButton>

              <UButton
                to="/customer/analytics/retention"
                variant="ghost"
                color="neutral"
                size="sm"
                class="w-full justify-start text-sm"
                :class="{
                  'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                    $route.path === '/customer/analytics/retention',
                }"
              >
                <UIcon name="i-heroicons-arrow-path" class="w-3 h-3 mr-2" />
                {{ $t("nav.retentionAnalysis") }}
              </UButton>

              <UButton
                to="/customer/analytics/behavior"
                variant="ghost"
                color="neutral"
                size="sm"
                class="w-full justify-start text-sm"
                :class="{
                  'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                    $route.path === '/customer/analytics/behavior',
                }"
              >
                <UIcon
                  name="i-heroicons-cursor-arrow-rays"
                  class="w-3 h-3 mr-2"
                />
                {{ $t("nav.behaviorAnalysis") }}
              </UButton>

              <UButton
                to="/customer/analytics/funnel"
                variant="ghost"
                color="neutral"
                size="sm"
                class="w-full justify-start text-sm"
                :class="{
                  'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                    $route.path === '/customer/analytics/funnel',
                }"
              >
                <UIcon name="i-heroicons-funnel" class="w-3 h-3 mr-2" />
                {{ $t("nav.funnelAnalysis") }}
              </UButton>
            </div>
          </div>
        </div>
      </div>

      <!-- 商品中心 -->
      <div>
        <div
          class="px-3 py-2 text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider"
        >
          {{ $t("nav.productCenter") }}
        </div>
        <div class="space-y-1">
          <UButton
            to="/product"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                $route.path === '/product',
            }"
          >
            <UIcon name="i-heroicons-cube" class="w-4 h-4 mr-3" />
            {{ $t("nav.products") }}
          </UButton>

          <UButton
            to="/product/categories"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                $route.path === '/product/categories',
            }"
          >
            <UIcon name="i-heroicons-tag" class="w-4 h-4 mr-3" />
            {{ $t("nav.categories") }}
          </UButton>

          <UButton
            to="/product/inventory"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                $route.path === '/product/inventory',
            }"
          >
            <UIcon name="i-heroicons-archive-box" class="w-4 h-4 mr-3" />
            {{ $t("nav.inventory") }}
          </UButton>

          <UButton
            to="/product/suppliers"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                $route.path === '/product/suppliers',
            }"
          >
            <UIcon name="i-heroicons-building-office" class="w-4 h-4 mr-3" />
            {{ $t("nav.suppliers") }}
          </UButton>
        </div>
      </div>

      <!-- 定价中心 -->
      <div>
        <div
          class="px-3 py-2 text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider"
        >
          {{ $t("pricing.center") }}
        </div>
        <div class="space-y-1">
          <UButton
            to="/pricing"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                $route.path === '/pricing',
            }"
          >
            <UIcon name="i-heroicons-currency-dollar" class="w-4 h-4 mr-3" />
            {{ $t("pricing.overview") }}
          </UButton>

          <UButton
            to="/pricing/pricebooks"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                $route.path === '/pricing/pricebooks',
            }"
          >
            <UIcon name="i-heroicons-book-open" class="w-4 h-4 mr-3" />
            {{ $t("pricing.pricebooks") }}
          </UButton>

          <UButton
            to="/pricing/rules"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                $route.path === '/pricing/rules',
            }"
          >
            <UIcon name="i-heroicons-cog-6-tooth" class="w-4 h-4 mr-3" />
            {{ $t("pricing.rules") }}
          </UButton>

          <UButton
            to="/pricing/contracts"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                $route.path === '/pricing/contracts',
            }"
          >
            <UIcon name="i-heroicons-document-text" class="w-4 h-4 mr-3" />
            {{ $t("pricing.contracts") }}
          </UButton>
        </div>
      </div>

      <!-- 营销增长 -->
      <div>
        <div
          class="px-3 py-2 text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider"
        >
          {{ $t("nav.marketingGrowth") }}
        </div>
        <div class="space-y-1">
          <UButton
            to="/market/orders"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                $route.path === '/market/orders',
            }"
          >
            <UIcon
              name="i-heroicons-clipboard-document-list"
              class="w-4 h-4 mr-3"
            />
            {{ $t("nav.orders") }}
          </UButton>

          <UButton
            to="/market/payment"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                $route.path === '/market/payment',
            }"
          >
            <UIcon
              name="i-heroicons-clipboard-document-list"
              class="w-4 h-4 mr-3"
            />
            {{ $t("nav.payment") }}
          </UButton>

          <UButton
            to="/market/marketing"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                $route.path === '/market/marketing',
            }"
          >
            <UIcon name="i-heroicons-megaphone" class="w-4 h-4 mr-3" />
            {{ $t("nav.marketing") }}
          </UButton>

          <UButton
            to="/market/promotions"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                $route.path === '/market/promotions',
            }"
          >
            <UIcon name="i-heroicons-gift" class="w-4 h-4 mr-3" />
            {{ $t("nav.promotions") }}
          </UButton>

          <UButton
            to="/market/channels"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                $route.path === '/market/channels',
            }"
          >
            <UIcon name="i-heroicons-globe-alt" class="w-4 h-4 mr-3" />
            {{ $t("nav.channels") }}
          </UButton>

          <UButton
            to="/market/analytics"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                $route.path === '/market/analytics',
            }"
          >
            <UIcon name="i-heroicons-chart-bar-square" class="w-4 h-4 mr-3" />
            {{ $t("nav.analytics") }}
          </UButton>
        </div>
      </div>
    </nav>
  </aside>
</template>

<script setup lang="ts">
const { t } = useI18n();

// 控制用户分析子菜单的展开状态
const showUserAnalytics = ref(false);

// 切换用户分析子菜单
const toggleUserAnalytics = () => {
  showUserAnalytics.value = !showUserAnalytics.value;
};

// 监听路由变化，如果当前路由在用户分析模块下，自动展开子菜单
const route = useRoute();
watch(
  () => route.path,
  (newPath) => {
    if (newPath.startsWith("/customer/analytics/")) {
      showUserAnalytics.value = true;
    }
  },
  { immediate: true }
);
</script>
