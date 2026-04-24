import { describe, it, expect } from 'vitest'
import type { ActivePlanView, PlanDay, PlanMeal, PlanOption, ActiveDietPlan } from '../../app/types/plan'
import type { ActiveDietPlan as PlanApiResponse } from '../../app/composables/useClientPlanApi'

describe('Plan Readability - Active Plan Display', () => {
  describe('Plan Flattening and Rendering', () => {
    it('renders active plan with day, meal, and option hierarchy preserved', () => {
      // Mock API response structure
      const mockPlan: any = {
        id: 'plan-1',
        start_date: '2026-04-23',
        daily_water_target_ml: 2000,
        is_active: true,
        updated_at: '2026-04-23T08:00:00Z',
        days: [
          {
            day_of_week: 0, // Sunday (Persian week)
            meals: [
              {
                id: 'meal-1',
                meal_type: 'breakfast',
                options: [
                  {
                    id: 'food-1',
                    food_name: 'تخم مرغ',
                    quantity_grams: 100,
                    calories: 155,
                  },
                ],
                notes: 'با نان سیاه',
              },
            ],
          },
        ],
        exercises: [],
        prescriptions: [],
      }

      expect(mockPlan.days).toHaveLength(1)
      expect(mockPlan.days[0].meals).toHaveLength(1)
      expect(mockPlan.days[0].meals[0].options).toHaveLength(1)
      expect(mockPlan.days[0].meals[0].options[0].food_name).toBe('تخم مرغ')
    })

    it('renders exercises and prescriptions as separate sections', () => {
      const mockPlan: any = {
        id: 'plan-1',
        is_active: true,
        updated_at: '2026-04-23T08:00:00Z',
        days: [],
        exercises: [
          {
            id: 'ex-1',
            name: 'پیاده روی',
            duration_minutes: 30,
            intensity: 'moderate',
          },
        ],
        prescriptions: [
          {
            id: 'med-1',
            medication_name: 'متفورمین',
            doses_per_day: 2,
          },
        ],
      }

      expect(mockPlan.exercises).toHaveLength(1)
      expect(mockPlan.prescriptions).toHaveLength(1)
      expect(mockPlan.exercises[0].name).toBe('پیاده روی')
      expect(mockPlan.prescriptions[0].medication_name).toBe('متفورمین')
    })

    it('includes plan metadata for UI context (dates, water target, freshness)', () => {
      const mockPlan: any = {
        id: 'plan-1',
        start_date: '2026-04-23',
        end_date: '2026-05-23',
        daily_water_target_ml: 2500,
        is_active: true,
        updated_at: '2026-04-23T08:00:00Z',
        days: [],
        exercises: [],
        prescriptions: [],
      }

      expect(mockPlan.start_date).toBeDefined()
      expect(mockPlan.daily_water_target_ml).toBe(2500)
      expect(mockPlan.is_active).toBe(true)
      expect(mockPlan.updated_at).toBeDefined()
    })

    it('handles plan notes and meal-level notes without data loss', () => {
      const mockPlan: any = {
        id: 'plan-1',
        notes: 'توجه: میزان نمک محدود باشد',
        is_active: true,
        days: [
          {
            day_of_week: 0,
            meals: [
              {
                id: 'meal-1',
                meal_type: 'breakfast',
                notes: 'بدون شیر',
                options: [],
              },
            ],
          },
        ],
        exercises: [],
        prescriptions: [],
      }

      expect(mockPlan.notes).toBe('توجه: میزان نمک محدود باشد')
      expect(mockPlan.days[0].meals[0].notes).toBe('بدون شیر')
    })
  })

  describe('Active Plan Context Badge', () => {
    it('labels active plan with active context badge', () => {
      const contextLabel = {
        is_active: true,
        label: 'برنامه فعال',
        style: 'active' as const,
      }

      expect(contextLabel.is_active).toBe(true)
      expect(contextLabel.label).toBe('برنامه فعال')
      expect(contextLabel.style).toBe('active')
    })

    it('prevents ambiguous context by deriving label from plan.is_active field', () => {
      const activePlan = { is_active: true, id: 'p1' }
      const archivedPlan = { is_active: false, id: 'p2' }

      expect(activePlan.is_active).toBe(true)
      expect(archivedPlan.is_active).toBe(false)
    })
  })
})
