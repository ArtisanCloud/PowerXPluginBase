<template>
  <div class="space-y-6">
    <!-- 页面头部 -->
    <div class="border-b border-gray-200 dark:border-gray-700 pb-4">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
        {{ $t("navigation.velocityChart") }}
      </h1>
      <p class="text-gray-600 dark:text-gray-400 mt-1">
        {{ $t("reports.velocityDescription") }}
      </p>
    </div>

    <!-- 过滤器 -->
    <UCard>
      <div class="flex flex-wrap gap-4">
        <USelectMenu
          v-model="selectedTeam"
          :options="teamOptions"
          placeholder="选择团队"
          class="w-48"
        />
        <USelectMenu
          v-model="selectedPeriod"
          :options="periodOptions"
          placeholder="时间周期"
          class="w-48"
        />
        <UButton variant="outline" @click="exportChart">
          {{ $t("reports.exportChart") }}
        </UButton>
      </div>
    </UCard>

    <!-- 速度图表 -->
    <UCard>
      <template #header>
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
          {{ $t("reports.teamVelocityTrend") }}
        </h2>
      </template>

      <div class="h-96 flex items-center justify-center bg-gray-50 dark:bg-gray-800 rounded-lg">
        <div class="text-center">
          <UIcon name="i-heroicons-chart-bar" class="w-16 h-16 text-gray-400 mx-auto mb-4" />
          <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">
            {{ $t("reports.velocityChartPlaceholder") }}
          </h3>
          <p class="text-gray-600 dark:text-gray-400">
            {{ $t("reports.velocityChartDescription") }}
          </p>
        </div>
      </div>
    </UCard>

    <!-- 速度统计 -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <UCard>
        <div class="text-center">
          <div class="text-2xl font-bold text-blue-600 dark:text-blue-400">{{ averageVelocity }}</div>
          <div class="text-sm text-gray-600 dark:text-gray-400">{{ $t("reports.averageVelocity") }}</div>
        </div>
      </UCard>
      <UCard>
        <div class="text-center">
          <div class="text-2xl font-bold text-green-600 dark:text-green-400">{{ currentVelocity }}</div>
          <div class="text-sm text-gray-600 dark:text-gray-400">{{ $t("reports.currentVelocity") }}</div>
        </div>
      </UCard>
      <UCard>
        <div class="text-center">
          <div class="text-2xl font-bold text-yellow-600 dark:text-yellow-400">{{ velocityTrend }}%</div>
          <div class="text-sm text-gray-600 dark:text-gray-400">{{ $t("reports.velocityTrend") }}</div>
        </div>
      </UCard>
      <UCard>
        <div class="text-center">
          <div class="text-2xl font-bold text-purple-600 dark:text-purple-400">{{ predictedVelocity }}</div>
          <div class="text-sm text-gray-600 dark:text-gray-400">{{ $t("reports.predictedVelocity") }}</div>
        </div>
      </UCard>
    </div>

    <!-- 历史速度数据 -->
    <UCard>
      <template #header>
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
          {{ $t("reports.velocityHistory") }}
        </h2>
      </template>

      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
          <thead class="bg-gray-50 dark:bg-gray-800">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                {{ $t("sprint.name") }}
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                {{ $t("reports.plannedPoints") }}
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                {{ $t("reports.completedPoints") }}
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                {{ $t("reports.velocity") }}
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                {{ $t("sprint.endDate") }}
              </th>
            </tr>
          </thead>
          <tbody class="bg-white dark:bg-gray-900 divide-y divide-gray-200 dark:divide-gray-700">
            <tr v-for="sprint in velocityHistory" :key="sprint.id">
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-white">
                {{ sprint.name }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600 dark:text-gray-400">
                {{ sprint.plannedPoints }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600 dark:text-gray-400">
                {{ sprint.completedPoints }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600 dark:text-gray-400">
                <UBadge :color="getVelocityColor(sprint.velocity)" variant="soft">
                  {{ sprint.velocity }}
                </UBadge>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600 dark:text-gray-400">
                {{ formatDate(sprint.endDate) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </UCard>
  </div>
</template>

<script setup>
// 页面元数据
definePageMeta({
  title: 'velocityChart'
});

// 响应式数据
const selectedTeam = ref(null);
const selectedPeriod = ref('last6sprints');

// 选项数据
const teamOptions = ref([
  { label: '开发团队A', value: 'team-a' },
  { label: '开发团队B', value: 'team-b' },
  { label: '全部团队', value: 'all' }
]);

const periodOptions = ref([
  { label: '最近6个冲刺', value: 'last6sprints' },
  { label: '最近12个冲刺', value: 'last12sprints' },
  { label: '本季度', value: 'thisQuarter' }
]);

// 统计数据
const averageVelocity = ref(45);
const currentVelocity = ref(52);
const velocityTrend = ref(15);
const predictedVelocity = ref(48);

// 历史数据
const velocityHistory = ref([
  {
    id: 1,
    name: 'Sprint 6',
    plannedPoints: 50,
    completedPoints: 52,
    velocity: 52,
    endDate: '2024-02-15'
  },
  {
    id: 2,
    name: 'Sprint 5',
    plannedPoints: 45,
    completedPoints: 43,
    velocity: 43,
    endDate: '2024-02-01'
  },
  {
    id: 3,
    name: 'Sprint 4',
    plannedPoints: 48,
    completedPoints: 47,
    velocity: 47,
    endDate: '2024-01-18'
  },
  {
    id: 4,
    name: 'Sprint 3',
    plannedPoints: 40,
    completedPoints: 38,
    velocity: 38,
    endDate: '2024-01-04'
  }
]);

// 方法
const exportChart = () => {
  console.log('导出速度图表');
};

const getVelocityColor = (velocity) => {
  if (velocity >= 50) return 'green';
  if (velocity >= 40) return 'yellow';
  return 'red';
};

const formatDate = (date) => {
  return new Date(date).toLocaleDateString();
};

// 初始化
onMounted(() => {
  selectedTeam.value = teamOptions.value[0].value;
});
</script>