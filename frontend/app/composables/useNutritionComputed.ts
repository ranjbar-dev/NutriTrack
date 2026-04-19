import type {
  MealOptionItemResponse,
  MealOptionResponse,
  MealResponse,
  PlanDayResponse,
  NutritionTotals,
} from '~/types/plan.types'

export function useNutritionComputed() {
  const zero: NutritionTotals = { calories: 0, protein_g: 0, carbs_g: 0, fat_g: 0, fiber_g: 0 }

  // D-14: serving ratio = item.quantity / item.food.measurement_amount
  // Guards division by zero (T-03-04-C)
  function itemTotals(item: MealOptionItemResponse): NutritionTotals {
    const ratio = item.food.measurement_amount > 0
      ? item.quantity / item.food.measurement_amount
      : 0
    return {
      calories: +(item.food.calories * ratio).toFixed(1),
      protein_g: +(item.food.protein_g * ratio).toFixed(1),
      carbs_g: +(item.food.carbs_g * ratio).toFixed(1),
      fat_g: +(item.food.fat_g * ratio).toFixed(1),
      fiber_g: +(item.food.fiber_g * ratio).toFixed(1),
    }
  }

  function sumTotals(totals: NutritionTotals[]): NutritionTotals {
    return totals.reduce((acc, t) => ({
      calories: +(acc.calories + t.calories).toFixed(1),
      protein_g: +(acc.protein_g + t.protein_g).toFixed(1),
      carbs_g: +(acc.carbs_g + t.carbs_g).toFixed(1),
      fat_g: +(acc.fat_g + t.fat_g).toFixed(1),
      fiber_g: +(acc.fiber_g + t.fiber_g).toFixed(1),
    }), { ...zero })
  }

  function optionTotals(option: MealOptionResponse): NutritionTotals {
    if (!option.items.length) return { ...zero }
    return sumTotals(option.items.map(itemTotals))
  }

  // Uses first option (option 1) as representative for meal totals (per UI-SPEC)
  function mealTotals(meal: MealResponse): NutritionTotals {
    if (!meal.options.length) return { ...zero }
    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
    return optionTotals(meal.options[0]!)
  }

  function dayTotals(day: PlanDayResponse): NutritionTotals {
    if (!day.meals.length) return { ...zero }
    return sumTotals(day.meals.map(mealTotals))
  }

  return { itemTotals, optionTotals, mealTotals, dayTotals }
}
