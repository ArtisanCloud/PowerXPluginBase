<script setup lang="ts">
/**
 * ContractPricingWizard.vue
 * 用途：新建/编辑 合同价格（选择客户 → 渠道与有效期 → 价格条款 → 选择商品 → 审批与优先级 → 预览保存）
 * Props:
 *  - modelValue: boolean  控制弹窗显隐
 *  - initialData?: any    编辑时传入
 * Emits:
 *  - update:modelValue
 *  - saved(payload)       保存成功时触发
 */

const props = defineProps<{
  modelValue: boolean;
  initialData?: any;
}>();
const emit = defineEmits<{
  (e: "update:modelValue", v: boolean): void;
  (e: "saved", payload: any): void;
}>();

const toast = useToast();

/* ====== 向导状态 ====== */
const step = ref(1);
const maxStep = 6;
const submitting = ref(false);
const conflictLoading = ref(false);
const conflicts = ref<
  {
    type: "价目表" | "规则" | "合同";
    title: string;
    severity: "info" | "warning" | "error";
    detail: string;
  }[]
>([]);

/* ====== 表单数据（统一容器） ====== */
const form = reactive({
  // step1 客户
  customer: {
    id: "",
    name: "",
    code: "",
    type: "",
    creditRating: "AA",
    salesRep: "",
  },
  // step2 渠道 & 有效期
  channels: [] as string[],
  dateRange: [] as any[], // [startDate, endDate]
  currency: "CNY",
  // step3 价格条款
  pricing: {
    mode: "固定折扣" as "固定折扣" | "固定价" | "阶梯折扣" | "批发价",
    value: "", // 折扣% 或 固定价
    basePrice: "标准价", // 市场价/标准价/批发价
    tiers: [{ min: 10, discount: 10 }], // 阶梯例子
    rounding: "四舍五入到元",
  },
  // step4 商品选择
  products: [] as any[], // 已选择商品列表
  // step5 审批与优先级
  approval: {
    approver: "",
    priority: "中" as "高" | "中" | "低",
    autoActivate: true,
    notes: "",
  },
});

/* ====== 选项 ====== */
const customerTypeOptions = [
  { label: "企业客户", value: "企业客户" },
  { label: "经销商", value: "经销商" },
  { label: "代理商", value: "代理商" },
  { label: "零售商", value: "零售商" },
];
const channelOptions = [
  { label: "直销", value: "直销" },
  { label: "经销", value: "经销" },
  { label: "在线", value: "在线" },
  { label: "零售", value: "零售" },
];
const roundingOptions = [
  "不舍入",
  "四舍五入到角",
  "四舍五入到元",
  "向下取整到元",
].map((v) => ({ label: v, value: v }));
const approverOptions = [
  { label: "价格委员会", value: "价格委员会" },
  { label: "大区总监", value: "大区总监" },
  { label: "定价经理", value: "定价经理" },
];

/* ====== 商品池（示例） ====== */
type ProductRow = {
  id: string;
  sku: string;
  name: string;
  category: string;
  msrp: number;
  std: number;
  stock: number;
};
const productPool = ref<ProductRow[]>([
  {
    id: "P001",
    sku: "IP15PM-256-NT",
    name: "iPhone 15 Pro Max",
    category: "手机",
    msrp: 9999,
    std: 8999,
    stock: 120,
  },
  {
    id: "P002",
    sku: "MBA-M2-256-SG",
    name: "MacBook Air M2",
    category: "电脑",
    msrp: 10499,
    std: 9499,
    stock: 55,
  },
  {
    id: "P003",
    sku: "APP2-USB-C",
    name: "AirPods Pro 2",
    category: "配件",
    msrp: 1999,
    std: 1799,
    stock: 320,
  },
]);
const productColumns = [
  { accessorKey: "sku", header: "SKU" },
  { accessorKey: "name", header: "商品名称" },
  { accessorKey: "category", header: "分类" },
  { accessorKey: "msrp", header: "MSRP" },
  { accessorKey: "std", header: "标准价" },
  { id: "select", header: "选择" },
] as const;

/* ====== 计算 & 校验 ====== */
const canNext = computed(() => {
  if (step.value === 1) {
    return !!form.customer.name && !!form.customer.type;
  }
  if (step.value === 2) {
    return form.channels.length > 0 && form.dateRange?.length === 2;
  }
  if (step.value === 3) {
    if (form.pricing.mode === "固定折扣")
      return !!form.pricing.value && Number(form.pricing.value) > 0;
    if (form.pricing.mode === "固定价")
      return !!form.pricing.value && Number(form.pricing.value) > 0;
    if (form.pricing.mode === "阶梯折扣") return form.pricing.tiers.length > 0;
    return true;
  }
  if (step.value === 4) return form.products.length > 0;
  if (step.value === 5) return !!form.approval.approver;
  return true;
});

function next() {
  if (!canNext.value) {
    toast.add({ title: "请完善当前步骤必填项", color: "red" });
    return;
  }
  if (step.value < maxStep) step.value++;
  if (step.value === 5) runConflictCheck(); // 进入审批与优先级时自动做冲突检测
}
function prev() {
  if (step.value > 1) step.value--;
}

/* ====== 冲突检测（示例逻辑） ====== */
async function runConflictCheck() {
  conflictLoading.value = true;
  await new Promise((r) => setTimeout(r, 500));
  // 伪造冲突结果：如果选了“在线”渠道，提示与价目表冲突；如果折扣>15%，提示与规则冲突
  const out: typeof conflicts.value = [];
  if (form.channels.includes("在线")) {
    out.push({
      type: "价目表",
      title: "与「在线-全国-春节特惠」价目表重叠",
      severity: "warning",
      detail:
        "同一渠道+日期范围存在重叠；建议设置合同优先级更高，或调整有效期。",
    });
  }
  const isHighDiscount =
    form.pricing.mode === "固定折扣" && Number(form.pricing.value) > 15;
  if (isHighDiscount) {
    out.push({
      type: "规则",
      title: "超出“批量折扣上限”规则",
      severity: "error",
      detail: "合同折扣高于规则上限 15%。需审批人明确豁免或下调折扣。",
    });
  }
  // 如果选择了“经销商”客户类型，并且模式为固定价，提醒与B2B合同模板不一致
  if (form.customer.type === "经销商" && form.pricing.mode === "固定价") {
    out.push({
      type: "合同",
      title: "与经销模板定价模式不一致",
      severity: "info",
      detail: "经销模板推荐使用“批发价/阶梯折扣”。",
    });
  }
  conflicts.value = out;
  conflictLoading.value = false;
}

/* ====== 商品选择（加入/移除） ====== */
function toggleProduct(p: ProductRow) {
  const idx = form.products.findIndex((x: any) => x.id === p.id);
  if (idx >= 0) form.products.splice(idx, 1);
  else {
    // 初始合同价按模式给出建议
    const suggested = suggestContractPrice(p);
    form.products.push({
      id: p.id,
      sku: p.sku,
      name: p.name,
      msrp: p.msrp,
      std: p.std,
      contractPrice: suggested.price,
      discount: suggested.discount,
    });
  }
}
function suggestContractPrice(p: ProductRow) {
  if (form.pricing.mode === "固定折扣") {
    const d = Number(form.pricing.value) || 0;
    const price = Math.round(p.std * (1 - d / 100));
    return { price, discount: d };
  }
  if (form.pricing.mode === "固定价") {
    const price = Number(form.pricing.value) || p.std;
    const disc = Math.round((1 - price / p.std) * 100);
    return { price, discount: disc };
  }
  if (form.pricing.mode === "批发价") {
    const price = Math.round(p.std * 0.95); // demo: 批发=标准价*95%
    const disc = Math.round((1 - price / p.std) * 100);
    return { price, discount: disc };
  }
  // 阶梯折扣：用第一档做预估
  const t0 = form.pricing.tiers[0];
  const d = t0?.discount ?? 0;
  const price = Math.round(p.std * (1 - d / 100));
  return { price, discount: d };
}

/* ====== 预览汇总 ====== */
const reviewRows = computed(() => {
  const prod = form.products;
  const total = prod.reduce(
    (acc, it) => acc + Number(it.contractPrice || 0),
    0
  );
  const avgDisc = prod.length
    ? Math.round(
        prod.reduce((acc, it) => acc + Number(it.discount || 0), 0) /
          prod.length
      )
    : 0;
  return { total, avgDisc };
});

/* ====== 保存 ====== */
async function handleSave() {
  // 硬性校验：有 error 冲突不允许直接激活
  const hasBlocker = conflicts.value.some((c) => c.severity === "error");
  if (hasBlocker && form.approval.autoActivate) {
    toast.add({
      title: "存在阻断冲突，需取消“自动生效”或处理冲突",
      color: "red",
    });
    return;
  }
  submitting.value = true;
  await new Promise((r) => setTimeout(r, 800));
  submitting.value = false;
  toast.add({ title: "合同价格已保存", color: "green" });
  emit("saved", {
    ...form,
    summary: {
      productCount: form.products.length,
      total: reviewRows.value.total,
    },
  });
  emit("update:modelValue", false);
}

/* ====== 编辑模式：载入initialData ====== */
onMounted(() => {
  if (props.initialData) {
    Object.assign(form, props.initialData);
  }
});
</script>

<template>
  <UModal
    :model-value="modelValue"
    @update:model-value="emit('update:modelValue', $event)"
    :ui="{ width: 'max-w-6xl' }"
  >
    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <div>
            <h3 class="text-lg font-semibold">合同价格向导</h3>
            <p class="text-xs text-gray-500">步骤 {{ step }} / {{ maxStep }}</p>
          </div>
          <UButton
            color="gray"
            variant="ghost"
            icon="i-heroicons-x-mark"
            @click="emit('update:modelValue', false)"
          />
        </div>
      </template>

      <!-- 步骤导航（简版） -->
      <div class="grid grid-cols-2 md:grid-cols-6 gap-2 mb-4">
        <UBadge
          v-for="i in maxStep"
          :key="i"
          :color="i === step ? 'primary' : 'gray'"
          variant="soft"
          class="justify-center"
        >
          {{
            [
              "选择客户",
              "渠道与有效期",
              "价格条款",
              "选择商品",
              "审批与优先级",
              "预览与保存",
            ][i - 1]
          }}
        </UBadge>
      </div>

      <!-- Step 1 选择客户 -->
      <div v-if="step === 1" class="space-y-4">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <UInput v-model="form.customer.name" placeholder="客户名称 *" />
          <UInput v-model="form.customer.code" placeholder="客户代码" />
          <USelect
            v-model="form.customer.type"
            :options="customerTypeOptions"
            placeholder="客户类型 *"
          />
          <USelect
            v-model="form.customer.creditRating"
            :options="
              ['AAA', 'AA+', 'AA', 'A+', 'A'].map((v) => ({
                label: v,
                value: v,
              }))
            "
            placeholder="信用等级"
          />
          <UInput v-model="form.customer.salesRep" placeholder="销售代表" />
        </div>
        <p class="text-xs text-gray-500">
          提示：客户类型会影响定价模板与冲突检测建议。
        </p>
      </div>

      <!-- Step 2 渠道与有效期 -->
      <div v-if="step === 2" class="space-y-4">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <USelect
            v-model="form.channels"
            multiple
            :options="channelOptions"
            placeholder="适用渠道 *"
          />
          <USelect
            v-model="form.currency"
            :options="
              ['CNY', 'USD', 'EUR'].map((v) => ({ label: v, value: v }))
            "
            placeholder="币种"
          />
          <UDatePicker v-model="form.dateRange" range placeholder="有效期 *" />
        </div>
        <p class="text-xs text-gray-500">
          建议：合同价格通常优先级高于公共价目表；请避免与活动价重叠的大范围有效期。
        </p>
      </div>

      <!-- Step 3 价格条款 -->
      <div v-if="step === 3" class="space-y-4">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
          <USelect
            v-model="form.pricing.mode"
            :options="
              ['固定折扣', '固定价', '阶梯折扣', '批发价'].map((v) => ({
                label: v,
                value: v,
              }))
            "
            placeholder="定价模式 *"
          />
          <UInput
            v-if="form.pricing.mode === '固定折扣'"
            v-model="form.pricing.value"
            type="number"
            placeholder="折扣%（如 10） *"
          />
          <UInput
            v-if="form.pricing.mode === '固定价'"
            v-model="form.pricing.value"
            type="number"
            placeholder="固定价格 *"
          />
          <USelect
            v-model="form.pricing.basePrice"
            :options="
              ['标准价', '市场价', '批发价'].map((v) => ({
                label: v,
                value: v,
              }))
            "
            placeholder="基准价格"
          />
          <USelect
            v-model="form.pricing.rounding"
            :options="roundingOptions"
            placeholder="舍入规则"
          />
        </div>

        <div v-if="form.pricing.mode === '阶梯折扣'" class="space-y-2">
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium">阶梯设置</span>
            <UButton
              size="xs"
              variant="outline"
              icon="i-heroicons-plus"
              @click="
                form.pricing.tiers.push({
                  min: (form.pricing.tiers.at(-1)?.min || 0) + 10,
                  discount: 10,
                })
              "
              >新增阶梯</UButton
            >
          </div>
          <div class="space-y-2">
            <div
              v-for="(t, idx) in form.pricing.tiers"
              :key="idx"
              class="grid grid-cols-3 gap-2"
            >
              <UInput
                v-model.number="t.min"
                type="number"
                placeholder="起订量"
              />
              <UInput
                v-model.number="t.discount"
                type="number"
                placeholder="折扣%"
              />
              <UButton
                size="xs"
                color="red"
                variant="soft"
                @click="form.pricing.tiers.splice(idx, 1)"
                >删除</UButton
              >
            </div>
          </div>
        </div>

        <p class="text-xs text-gray-500">
          提示：固定折扣是最常用模式；阶梯折扣适合大单/长期协议。
        </p>
      </div>

      <!-- Step 4 选择商品 -->
      <div v-if="step === 4" class="space-y-4">
        <UTable :data="productPool" :columns="productColumns">
          <template #msrp-cell="{ row }"
            >¥{{ row.original.msrp.toLocaleString() }}</template
          >
          <template #std-cell="{ row }"
            >¥{{ row.original.std.toLocaleString() }}</template
          >
          <template #select-cell="{ row }">
            <UToggle
              :model-value="
                !!form.products.find((x: any) => x.id === row.original.id)
              "
              @update:model-value="toggleProduct(row.original)"
            />
          </template>
        </UTable>

        <div v-if="form.products.length" class="mt-4">
          <h4 class="font-medium mb-2">
            已选商品（{{ form.products.length }}）
          </h4>
          <div class="space-y-2">
            <div
              v-for="(p, idx) in form.products"
              :key="p.id"
              class="grid grid-cols-5 gap-2 items-center"
            >
              <div class="col-span-2 text-sm">
                <p class="font-medium">{{ p.name }}</p>
                <p class="text-xs text-gray-500">{{ p.sku }}</p>
              </div>
              <div class="text-right text-sm">MSRP ¥{{ p.msrp }}</div>
              <div class="text-right text-sm">标准价 ¥{{ p.std }}</div>
              <div class="flex items-center gap-2">
                <UInput
                  v-model.number="p.contractPrice"
                  type="number"
                  placeholder="合同价"
                />
                <span class="text-xs text-gray-500"
                  >折扣 {{ p.discount }}%</span
                >
                <UButton
                  size="xs"
                  variant="soft"
                  color="gray"
                  @click="form.products.splice(idx, 1)"
                  >移除</UButton
                >
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Step 5 审批与优先级 & 冲突检测 -->
      <div v-if="step === 5" class="space-y-4">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
          <USelect
            v-model="form.approval.approver"
            :options="approverOptions"
            placeholder="审批人 *"
          />
          <USelect
            v-model="form.approval.priority"
            :options="['高', '中', '低'].map((v) => ({ label: v, value: v }))"
            placeholder="优先级"
          />
          <UToggle v-model="form.approval.autoActivate" />
          <UInput
            v-model="form.approval.notes"
            placeholder="审批备注（可选）"
          />
        </div>

        <UCard>
          <template #header>
            <div class="flex justify-between items-center">
              <h4 class="font-medium">冲突检测</h4>
              <UButton
                size="xs"
                variant="outline"
                icon="i-heroicons-arrow-path"
                :loading="conflictLoading"
                @click="runConflictCheck"
                >重新检测</UButton
              >
            </div>
          </template>

          <div v-if="conflicts && conflicts.length > 0">
            <div
              v-for="(c, i) in conflicts"
              :key="i"
              class="flex items-start gap-3 p-3 rounded border"
              :class="
                c.severity === 'error'
                  ? 'border-red-200 bg-red-50'
                  : c.severity === 'warning'
                    ? 'border-yellow-200 bg-yellow-50'
                    : 'border-blue-200 bg-blue-50'
              "
            >
              <UBadge
                :color="
                  c.severity === 'error'
                    ? 'error'
                    : c.severity === 'warning'
                      ? 'warning'
                      : 'info'
                "
                variant="soft"
                >{{ c.type }}</UBadge
              >
              <div class="flex-1">
                <p class="font-medium">{{ c.title }}</p>
                <p class="text-xs text-gray-600">{{ c.detail }}</p>
              </div>
              <UButton size="xs" variant="ghost">查看详情</UButton>
            </div>
            <p class="text-xs text-gray-500 mt-2">
              建议：遇到
              <span class="text-red-600 font-medium">阻断冲突</span>
              时，需审批豁免或调整参数（有效期/渠道/折扣/优先级）。
            </p>
          </div>
          <template #noConflict>
            <div class="text-center text-gray-500">暂无冲突</div>
          </template>
        </UCard>
      </div>

      <!-- Step 6 预览与保存 -->
      <div v-if="step === 6" class="space-y-4">
        <UCard>
          <template #header><h4 class="font-medium">概要</h4></template>
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
            <div>
              <span class="text-gray-500">客户：</span
              >{{ form.customer.name }}（{{ form.customer.type }}）
            </div>
            <div>
              <span class="text-gray-500">渠道：</span
              >{{ form.channels.join("、") }}
            </div>
            <div>
              <span class="text-gray-500">有效期：</span
              >{{ form.dateRange?.[0] }} ~ {{ form.dateRange?.[1] }}
            </div>
            <div>
              <span class="text-gray-500">定价模式：</span
              >{{ form.pricing.mode }}（基准：{{ form.pricing.basePrice }}）
            </div>
            <div>
              <span class="text-gray-500">舍入：</span
              >{{ form.pricing.rounding }}
            </div>
            <div>
              <span class="text-gray-500">审批：</span
              >{{ form.approval.approver }}（优先级：{{
                form.approval.priority
              }}）
            </div>
          </div>
        </UCard>

        <UCard>
          <template #header><h4 class="font-medium">商品与价格</h4></template>
          <div class="space-y-2 text-sm">
            <div
              v-for="p in form.products"
              :key="p.id"
              class="grid grid-cols-5 gap-2"
            >
              <div class="col-span-2">
                {{ p.name }}
                <span class="text-xs text-gray-500">({{ p.sku }})</span>
              </div>
              <div class="text-right">标准价 ¥{{ p.std }}</div>
              <div class="text-right">合同价 ¥{{ p.contractPrice }}</div>
              <div class="text-right text-green-600">
                折扣 {{ p.discount }}%
              </div>
            </div>
            <div class="border-t pt-2 text-right text-sm">
              平均折扣：<span class="font-medium"
                >{{ reviewRows.avgDisc }}%</span
              >， 合同小计：<span class="font-bold"
                >¥{{ reviewRows.total.toLocaleString() }}</span
              >
            </div>
          </div>
        </UCard>
      </div>

      <template #footer>
        <div class="flex justify-between">
          <UButton variant="ghost" @click="emit('update:modelValue', false)"
            >取消</UButton
          >
          <div class="flex gap-2">
            <UButton v-if="step > 1" variant="outline" @click="prev"
              >上一步</UButton
            >
            <UButton v-if="step < maxStep" :disabled="!canNext" @click="next"
              >下一步</UButton
            >
            <UButton
              v-else
              color="primary"
              :loading="submitting"
              @click="handleSave"
              >保存</UButton
            >
          </div>
        </div>
      </template>
    </UCard>
  </UModal>
</template>

<style scoped>
/* 可按需微调 */
</style>
