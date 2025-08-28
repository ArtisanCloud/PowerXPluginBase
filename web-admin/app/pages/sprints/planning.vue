<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-semibold text-gray-900 dark:text-white">
          {{ $t("navigation.sprintPlanning") }}
        </h1>
        <p class="mt-2 text-sm text-gray-600 dark:text-gray-400">
          {{ $t("planning.description") }}
        </p>
      </div>
      <UButton
        color="primary"
        icon="i-heroicons-plus"
        @click="createSprint"
      >
        {{ $t("sprint.createSprint") }}
      </UButton>
    </div>

    <!-- 计划步骤 -->
    <UCard>
      <div class="flex items-center justify-between mb-6">
        <h3 class="text-lg font-medium text-gray-900 dark:text-white">
          {{ $t("planning.planningProcess") }}
        </h3>
        <UBadge variant="soft" color="blue">
          {{ $t("planning.step") }} {{ currentStep }}/4
        </UBadge>
      </div>

      <!-- 步骤指示器 -->
      <div class="flex items-center justify-between mb-8">
        <div
          v-for="(step, index) in planningSteps"
          :key="index"
          class="flex items-center"
          :class="{ 'flex-1': index < planningSteps.length - 1 }"
        >
          <div class="flex items-center">
            <div
              :class="[
                'w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium',
                index + 1 <= currentStep
                  ? 'bg-primary-600 text-white'
                  : 'bg-gray-200 dark:bg-gray-700 text-gray-500 dark:text-gray-400'
              ]"
            >
              {{ index + 1 }}
            </div>
            <div class="ml-3">
              <div class="text-sm font-medium text-gray-900 dark:text-white">
                {{ step.title }}
              </div>
              <div class="text-xs text-gray-500 dark:text-gray-400">
                {{ step.description }}
              </div>
            </div>
          </div>
          <div
            v-if="index < planningSteps.length - 1"
            :class="[
              'h-0.5 flex-1 mx-4',
              index + 1 < currentStep
                ? 'bg-primary-600'
                : 'bg-gray-200 dark:bg-gray-700'
            ]"
          />
        </div>
      </div>

      <!-- 步骤内容 -->
      <div v-if="currentStep === 1" class="space-y-4">
        <h4 class="font-medium text-gray-900 dark:text-white">
          {{ $t("planning.defineGoal") }}
        </h4>
        <UTextarea
          v-model="sprintGoal"
          :placeholder="$t('planning.goalPlaceholder')"
          rows="3"
        />
      </div>

      <div v-if="currentStep === 2" class="space-y-4">
        <h4 class="font-medium text-gray-900 dark:text-white">
          {{ $t("planning.selectItems") }}
        </h4>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <!-- 产品待办列表 -->
          <div>
            <h5 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">
              {{ $t("backlog.productBacklog") }}
            </h5>
            <div class="space-y-2 max-h-64 overflow-y-auto">
              <div
                v-for="item in backlogItems"
                :key="item.id"
                class="p-3 border border-gray-200 dark:border-gray-700 rounded-lg cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-800"
                @click="addToSprint(item)"
              >
                <div class="flex items-center justify-between">
                  <div class="flex-1">
                    <div class="text-sm font-medium text-gray-900 dark:text-white">
                      {{ item.title }}
                    </div>
                    <div class="text-xs text-gray-500 dark:text-gray-400">
                      {{ $t("task.storyPoints") }}: {{ item.storyPoints }}
                    </div>
                  </div>
                  <UButton
                    icon="i-heroicons-plus"
                    size="xs"
                    variant="outline"
                  />
                </div>
              </div>
            </div>
          </div>

          <!-- 冲刺待办列表 -->
          <div>
            <h5 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">
              {{ $t("sprint.sprintBacklog") }}
            </h5>
            <div class="space-y-2 max-h-64 overflow-y-auto">
              <div
                v-for="item in sprintItems"
                :key="item.id"
                class="p-3 border border-blue-200 dark:border-blue-700 bg-blue-50 dark:bg-blue-900/20 rounded-lg"
              >
                <div class="flex items-center justify-between">
                  <div class="flex-1">
                    <div class="text-sm font-medium text-gray-900 dark:text-white">
                      {{ item.title }}
                    </div>
                    <div class="text-xs text-gray-500 dark:text-gray-400">
                      {{ $t("task.storyPoints") }}: {{ item.storyPoints }}
                    </div>
                  </div>
                  <UButton
                    icon="i-heroicons-minus"
                    size="xs"
                    variant="outline"
                    color="red"
                    @click="removeFromSprint(item)"
                  />
                </div>
              </div>
            </div>
            <div class="mt-3 text-sm text-gray-600 dark:text-gray-400">
              {{ $t("planning.totalPoints") }}: {{ totalSprintPoints }}
            </div>
          </div>
        </div>
      </div>

      <div v-if="currentStep === 3" class="space-y-4">
        <h4 class="font-medium text-gray-900 dark:text-white">
          {{ $t("planning.setDuration") }}
        </h4>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              {{ $t("sprint.startDate") }}
            </label>
            <UInput
              v-model="sprintStartDate"
              type="date"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              {{ $t("sprint.endDate") }}
            </label>
            <UInput
              v-model="sprintEndDate"
              type="date"
            />
          </div>
        </div>
      </div>

      <div v-if="currentStep === 4" class="space-y-4">
        <h4 class="font-medium text-gray-900 dark:text-white">
          {{ $t("planning.reviewAndStart") }}
        </h4>
        <div class="bg-gray-50 dark:bg-gray-800 p-4 rounded-lg">
          <div class="space-y-3">
            <div>
              <span class="text-sm font-medium text-gray-700 dark:text-gray-300">
                {{ $t("sprint.goal") }}:
              </span>
              <span class="text-sm text-gray-900 dark:text-white ml-2">
                {{ sprintGoal || $t("planning.noGoal") }}
              </span>
            </div>
            <div>
              <span class="text-sm font-medium text-gray-700 dark:text-gray-300">
                {{ $t("planning.selectedItems") }}:
              </span>
              <span class="text-sm text-gray-900 dark:text-white ml-2">
                {{ sprintItems.length }} {{ $t("planning.items") }}
              </span>
            </div>
            <div>
              <span class="text-sm font-medium text-gray-700 dark:text-gray-300">
                {{ $t("planning.totalPoints") }}:
              </span>
              <span class="text-sm text-gray-900 dark:text-white ml-2">
                {{ totalSprintPoints }}
              </span>
            </div>
            <div>
              <span class="text-sm font-medium text-gray-700 dark:text-gray-300">
                {{ $t("planning.duration") }}:
              </span>
              <span class="text-sm text-gray-900 dark:text-white ml-2">
                {{ sprintStartDate }} - {{ sprintEndDate }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- 导航按钮 -->
      <div class="flex justify-between mt-8">
        <UButton
          v-if="currentStep > 1"
          variant="outline"
          @click="previousStep"
        >
          {{ $t("common.previous") }}
        </UButton>
        <div v-else/>

        <div class="flex space-x-2">
          <UButton
            v-if="currentStep < 4"
            @click="nextStep"
          >
            {{ $t("common.next") }}
          </UButton>
          <UButton
            v-if="currentStep === 4"
            color="green"
            @click="startSprint"
          >
            {{ $t("planning.startSprint") }}
          </UButton>
        </div>
      </div>
    </UCard>
  </div>
</template>

<script setup>
// 国际化
const { t } = useI18n();

// 页面元数据
useSeoMeta({
  title: () => `${t('common.appName')} - ${t('navigation.sprintPlanning')}`,
  description: () => t('planning.description')
});

// 当前步骤
const currentStep = ref(1);

// 计划步骤
const planningSteps = computed(() => [
  {
    title: t('planning.defineGoal'),
    description: t('planning.defineGoalDesc')
  },
  {
    title: t('planning.selectItems'),
    description: t('planning.selectItemsDesc')
  },
  {
    title: t('planning.setDuration'),
    description: t('planning.setDurationDesc')
  },
  {
    title: t('planning.reviewAndStart'),
    description: t('planning.reviewAndStartDesc')
  }
]);

// 表单数据
const sprintGoal = ref('');
const sprintStartDate = ref('');
const sprintEndDate = ref('');
const sprintItems = ref([]);

// Mock 产品待办列表
const backlogItems = ref([
  {
    id: 'US-001',
    title: '用户注册功能',
    storyPoints: 5
  },
  {
    id: 'US-002',
    title: '用户登录功能',
    storyPoints: 3
  },
  {
    id: 'US-003',
    title: '密码重置功能',
    storyPoints: 2
  },
  {
    id: 'US-004',
    title: '个人资料管理',
    storyPoints: 8
  }
]);

// 计算属性
const totalSprintPoints = computed(() => {
  return sprintItems.value.reduce((sum, item) => sum + item.storyPoints, 0);
});

// 方法
const nextStep = () => {
  if (currentStep.value < 4) {
    currentStep.value++;
  }
};

const previousStep = () => {
  if (currentStep.value > 1) {
    currentStep.value--;
  }
};

const addToSprint = (item) => {
  if (!sprintItems.value.find(i => i.id === item.id)) {
    sprintItems.value.push(item);
  }
};

const removeFromSprint = (item) => {
  const index = sprintItems.value.findIndex(i => i.id === item.id);
  if (index > -1) {
    sprintItems.value.splice(index, 1);
  }
};

const createSprint = () => {
  // TODO: 实现创建冲刺功能
  console.log('创建冲刺');
};

const startSprint = () => {
  // TODO: 实现启动冲刺功能
  console.log('启动冲刺');
  navigateTo('/sprints/active');
};
</script>