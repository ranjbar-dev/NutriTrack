<script setup lang="ts">
import { toLatinDigits } from '~/utils/persian-digits'
import type { CreateFoodPayload, FoodResponse } from '~/stores/food'

definePageMeta({
  layout: 'nutritionist',
  middleware: ['role-guard'],
  roles: ['nutritionist'],
})

const foodStore = useFoodStore()
const route = useRoute()
const isSubmitting = ref(false)
const isLoading = ref(true)
const showToast = ref(false)
const toastMessage = ref('')

const categoryOptions = [
  { value: 'breakfast', label: 'صبحانه' },
  { value: 'lunch', label: 'ناهار' },
  { value: 'dinner', label: 'شام' },
  { value: 'snack', label: 'میان‌وعده' },
  { value: 'fruit', label: 'میوه' },
  { value: 'beverage', label: 'نوشیدنی' },
  { value: 'supplement', label: 'مکمل' },
  { value: 'other', label: 'سایر' },
]

const measurementUnits = [
  { value: 'gram', label: 'گرم' },
  { value: 'kg', label: 'کیلوگرم' },
  { value: 'tablespoon', label: 'قاشق غذاخوری' },
  { value: 'teaspoon', label: 'قاشق چایخوری' },
  { value: 'cup', label: 'لیوان' },
  { value: 'piece', label: 'عدد' },
  { value: 'slice', label: 'برش' },
  { value: 'palm', label: 'کف دست' },
  { value: 'matchbox', label: 'قوطی کبریت' },
  { value: 'bowl', label: 'کاسه' },
  { value: 'ml', label: 'میلی‌لیتر' },
  { value: 'liter', label: 'لیتر' },
]

const form = reactive({
  name: '',
  description: '',
  categories: [] as string[],
  calories: '',
  protein_g: '',
  carbs_g: '',
  fat_g: '',
  fiber_g: '',
  sugar_g: '',
  sodium_mg: '',
  measurement_unit: 'gram',
  measurement_amount: '',
})

const errors = reactive({
  name: '',
  description: '',
  categories: '',
  calories: '',
  protein_g: '',
  carbs_g: '',
  fat_g: '',
  fiber_g: '',
  sugar_g: '',
  sodium_mg: '',
  measurement_unit: '',
  measurement_amount: '',
})

const nutritionLimits: Record<string, number> = {
  calories: 9999.99,
  protein_g: 999.99,
  carbs_g: 999.99,
  fat_g: 999.99,
  fiber_g: 999.99,
  sugar_g: 999.99,
  sodium_mg: 9999.99,
}

const foodId = computed(() => route.params.id as string)

function triggerToast(message: string) {
  toastMessage.value = message
  showToast.value = true
  setTimeout(() => {
    showToast.value = false
  }, 3000)
}

function resetErrors() {
  Object.keys(errors).forEach((key) => {
    errors[key as keyof typeof errors] = ''
  })
}

function parseNumber(value: string) {
  const normalized = toLatinDigits(value).trim()
  if (!normalized) return NaN
  if (!/^\d+(\.\d{1,2})?$/.test(normalized)) return NaN
  return Number.parseFloat(normalized)
}

function validateNumberField(
  key: keyof typeof form,
  label: string,
  max: number,
  allowZero = true,
) {
  const value = parseNumber(form[key] as string)
  if (Number.isNaN(value)) {
    errors[key as keyof typeof errors] = `${label} باید عدد معتبر باشد`
    return false
  }
  if (value < 0 || (!allowZero && value <= 0)) {
    errors[key as keyof typeof errors] = allowZero
      ? `${label} نمی‌تواند منفی باشد`
      : `${label} باید بزرگ‌تر از صفر باشد`
    return false
  }
  if (value > max) {
    errors[key as keyof typeof errors] = `${label} نباید بیشتر از ${max} باشد`
    return false
  }
  return true
}

function validateForm() {
  resetErrors()
  let isValid = true

  if (!form.name.trim()) {
    errors.name = 'نام غذا الزامی است'
    isValid = false
  }
  else if (form.name.trim().length > 200) {
    errors.name = 'نام غذا نباید بیشتر از ۲۰۰ کاراکتر باشد'
    isValid = false
  }

  if (form.description && form.description.length > 1000) {
    errors.description = 'توضیحات نباید بیشتر از ۱۰۰۰ کاراکتر باشد'
    isValid = false
  }

  if (form.categories.length === 0) {
    errors.categories = 'حداقل یک دسته را انتخاب کنید'
    isValid = false
  }

  Object.entries(nutritionLimits).forEach(([key, max]) => {
    const labelMap: Record<string, string> = {
      calories: 'کالری',
      protein_g: 'پروتئین',
      carbs_g: 'کربوهیدرات',
      fat_g: 'چربی',
      fiber_g: 'فیبر',
      sugar_g: 'قند',
      sodium_mg: 'سدیم',
    }
    if (
      !validateNumberField(key as keyof typeof form, labelMap[key] ?? key, max, true)
    ) {
      isValid = false
    }
  })

  if (!form.measurement_unit) {
    errors.measurement_unit = 'واحد اندازه‌گیری الزامی است'
    isValid = false
  }

  if (!validateNumberField('measurement_amount', 'مقدار', 999999, false)) {
    isValid = false
  }

  return isValid
}

function buildPayload(): CreateFoodPayload {
  return {
    name: form.name.trim(),
    description: form.description.trim() || undefined,
    categories: [...form.categories],
    calories: parseNumber(form.calories),
    protein_g: parseNumber(form.protein_g),
    carbs_g: parseNumber(form.carbs_g),
    fat_g: parseNumber(form.fat_g),
    fiber_g: parseNumber(form.fiber_g),
    sugar_g: parseNumber(form.sugar_g),
    sodium_mg: parseNumber(form.sodium_mg),
    measurement_unit: form.measurement_unit,
    measurement_amount: parseNumber(form.measurement_amount),
  }
}

function populateForm(food: FoodResponse) {
  form.name = food.name
  form.description = food.description ?? ''
  form.categories = [...food.categories]
  form.calories = String(food.calories)
  form.protein_g = String(food.protein_g)
  form.carbs_g = String(food.carbs_g)
  form.fat_g = String(food.fat_g)
  form.fiber_g = String(food.fiber_g)
  form.sugar_g = String(food.sugar_g)
  form.sodium_mg = String(food.sodium_mg)
  form.measurement_unit = food.measurement_unit
  form.measurement_amount = String(food.measurement_amount)
}

async function loadFood() {
  isLoading.value = true
  try {
    const data = await foodStore.fetchFood(foodId.value)
    populateForm(data)
  }
  catch {
    triggerToast('بارگذاری غذا با خطا مواجه شد')
    navigateTo('/nutritionist/foods')
  }
  finally {
    isLoading.value = false
  }
}

async function handleSubmit() {
  if (!validateForm()) {
    triggerToast('لطفاً خطاهای فرم را اصلاح کنید')
    return
  }

  isSubmitting.value = true
  try {
    await foodStore.updateFood(foodId.value, buildPayload())
    navigateTo('/nutritionist/foods')
  }
  catch (error) {
    triggerToast((error as Error).message || 'ویرایش غذا با خطا مواجه شد')
  }
  finally {
    isSubmitting.value = false
  }
}

onMounted(() => {
  loadFood()
})
</script>

<template>
  <div class="p-4 space-y-4">
    <div
      v-if="showToast"
      class="fixed top-4 start-1/2 -translate-x-1/2 bg-red-600 text-white px-4 py-2 rounded-lg shadow-lg z-50 text-sm"
    >
      {{ toastMessage }}
    </div>

    <header class="flex items-center justify-between gap-3">
      <h1 class="text-xl font-bold text-gray-800">ویرایش غذا</h1>
      <UiAppButton
        class="w-auto"
        size="sm"
        variant="secondary"
        @click="navigateTo('/nutritionist/foods')"
      >
        بازگشت
      </UiAppButton>
    </header>

    <div v-if="isLoading" class="py-16">
      <UiLoadingSpinner size="lg" />
    </div>

    <form v-else class="space-y-4" @submit.prevent="handleSubmit">
      <section class="bg-white rounded-2xl p-4 shadow-sm space-y-4">
        <h2 class="text-sm font-semibold text-gray-700">اطلاعات پایه</h2>
        <UiAppInput
          v-model="form.name"
          label="نام غذا"
          :disabled="isSubmitting"
          :error="errors.name"
        />
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">
            توضیحات
          </label>
          <textarea
            v-model="form.description"
            class="w-full rounded-lg border px-3 py-2.5 text-base transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-offset-0"
            :class="[
              errors.description
                ? 'border-red-400 focus:border-red-500 focus:ring-red-200'
                : 'border-gray-300 focus:border-emerald-500 focus:ring-emerald-200',
              isSubmitting ? 'bg-gray-100 cursor-not-allowed' : 'bg-white',
            ]"
            rows="3"
            :disabled="isSubmitting"
            maxlength="1000"
          />
          <p v-if="errors.description" class="mt-1 text-sm text-red-600">
            {{ errors.description }}
          </p>
        </div>
      </section>

      <section class="bg-white rounded-2xl p-4 shadow-sm space-y-4">
        <h2 class="text-sm font-semibold text-gray-700">دسته‌بندی</h2>
        <div class="grid grid-cols-2 gap-2">
          <label
            v-for="category in categoryOptions"
            :key="category.value"
            class="flex items-center gap-2 text-sm text-gray-700 bg-gray-50 rounded-lg px-3 py-2"
          >
            <input
              v-model="form.categories"
              type="checkbox"
              :value="category.value"
              :disabled="isSubmitting"
              class="accent-emerald-600"
            />
            <span>{{ category.label }}</span>
          </label>
        </div>
        <p v-if="errors.categories" class="text-sm text-red-600">
          {{ errors.categories }}
        </p>
      </section>

      <section class="bg-white rounded-2xl p-4 shadow-sm space-y-4">
        <h2 class="text-sm font-semibold text-gray-700">اطلاعات تغذیه‌ای</h2>
        <div class="grid grid-cols-2 gap-3">
          <UiAppInput
            v-model="form.calories"
            label="کالری"
            type="number"
            inputDir="ltr"
            :disabled="isSubmitting"
            :error="errors.calories"
          />
          <UiAppInput
            v-model="form.protein_g"
            label="پروتئین"
            type="number"
            inputDir="ltr"
            :disabled="isSubmitting"
            :error="errors.protein_g"
          />
          <UiAppInput
            v-model="form.carbs_g"
            label="کربوهیدرات"
            type="number"
            inputDir="ltr"
            :disabled="isSubmitting"
            :error="errors.carbs_g"
          />
          <UiAppInput
            v-model="form.fat_g"
            label="چربی"
            type="number"
            inputDir="ltr"
            :disabled="isSubmitting"
            :error="errors.fat_g"
          />
          <UiAppInput
            v-model="form.fiber_g"
            label="فیبر"
            type="number"
            inputDir="ltr"
            :disabled="isSubmitting"
            :error="errors.fiber_g"
          />
          <UiAppInput
            v-model="form.sugar_g"
            label="قند"
            type="number"
            inputDir="ltr"
            :disabled="isSubmitting"
            :error="errors.sugar_g"
          />
        </div>
        <UiAppInput
          v-model="form.sodium_mg"
          label="سدیم"
          type="number"
          inputDir="ltr"
          :disabled="isSubmitting"
          :error="errors.sodium_mg"
        />
      </section>

      <section class="bg-white rounded-2xl p-4 shadow-sm space-y-4">
        <h2 class="text-sm font-semibold text-gray-700">واحد اندازه‌گیری</h2>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">
            واحد
          </label>
          <select
            v-model="form.measurement_unit"
            class="w-full rounded-lg border px-3 py-2.5 text-base transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-offset-0"
            :class="[
              errors.measurement_unit
                ? 'border-red-400 focus:border-red-500 focus:ring-red-200'
                : 'border-gray-300 focus:border-emerald-500 focus:ring-emerald-200',
              isSubmitting ? 'bg-gray-100 cursor-not-allowed' : 'bg-white',
            ]"
            :disabled="isSubmitting"
          >
            <option v-for="unit in measurementUnits" :key="unit.value" :value="unit.value">
              {{ unit.label }}
            </option>
          </select>
          <p v-if="errors.measurement_unit" class="mt-1 text-sm text-red-600">
            {{ errors.measurement_unit }}
          </p>
        </div>
        <UiAppInput
          v-model="form.measurement_amount"
          label="مقدار"
          type="number"
          inputDir="ltr"
          :disabled="isSubmitting"
          :error="errors.measurement_amount"
        />
      </section>

      <UiAppButton type="submit" :loading="isSubmitting" :disabled="isSubmitting">
        {{ isSubmitting ? 'در حال ذخیره...' : 'ذخیره تغییرات' }}
      </UiAppButton>
    </form>
  </div>
</template>
