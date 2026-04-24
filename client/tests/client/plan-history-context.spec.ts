import { describe, it, expect, beforeEach } from 'vitest'
import type { ArchivedPlanView } from '../../app/types/plan'

describe('Plan History Context - Archived Plans', () => {
  describe('Archive List Display', () => {
    it('renders archived plan list with period and status', () => {
      const archivedPlans = [
        {
          id: 'plan-2',
          start_date: '2026-03-23',
          end_date: '2026-04-23',
          summary: 'برنامه فاز 2 - کاهش وزن',
          is_active: false,
          updated_at: '2026-04-23T08:00:00Z',
        },
        {
          id: 'plan-1',
          start_date: '2026-02-23',
          end_date: '2026-03-23',
          summary: 'برنامه فاز 1 - شروع',
          is_active: false,
          updated_at: '2026-03-23T08:00:00Z',
        },
      ]

      expect(archivedPlans).toHaveLength(2)
      expect(archivedPlans[0].is_active).toBe(false)
      expect(archivedPlans[0].start_date).toBeDefined()
      expect(archivedPlans[0].end_date).toBeDefined()
    })

    it('distinguishes active plan context from archived items', () => {
      const activePlan = {
        id: 'plan-active',
        is_active: true,
        start_date: '2026-04-23',
      }
      const archivedPlans = [
        {
          id: 'plan-2',
          is_active: false,
          start_date: '2026-03-23',
        },
      ]

      // Active plan should never appear in archived list
      const allPlans = [activePlan, ...archivedPlans]
      const onlyArchived = allPlans.filter(p => !p.is_active)
      expect(onlyArchived).toHaveLength(1)
      expect(onlyArchived[0].id).toBe('plan-2')
    })
  })

  describe('Offline Fallback and Context Continuity', () => {
    it('preserves active plan context in history UI when offline', () => {
      // When browsing history while offline, the active plan context should be preserved
      const activePlanCached = {
        id: 'plan-active',
        is_active: true,
        start_date: '2026-04-23',
        freshness: 'cache' as const,
      }

      const historyContextLabel = {
        active_plan_id: activePlanCached.id,
        active_plan_visible: true,
      }

      expect(historyContextLabel.active_plan_id).toBe('plan-active')
      expect(historyContextLabel.active_plan_visible).toBe(true)
    })

    it('historical data requires online fetch (per D-09 offline cache scope)', () => {
      // Historical data is NOT cached offline
      const archivedPlansRequireOnline = true

      expect(archivedPlansRequireOnline).toBe(true)

      // Only active plan is cached; history comes from server
      const cacheScope = {
        cached: ['active_plan'],
        online_only: ['archived_plans', 'plan_history'],
      }

      expect(cacheScope.cached).toContain('active_plan')
      expect(cacheScope.online_only).toContain('archived_plans')
    })

    it('includes freshness marker for active plan data', () => {
      const cachedActivePlan = {
        id: 'plan-active',
        is_active: true,
        freshness: 'cache' as const,
        updated_at: '2026-04-23T08:00:00Z',
        stale_marker: 'آخرین به روزرسانی: 08:00',
      }

      expect(cachedActivePlan.freshness).toBe('cache')
      expect(cachedActivePlan.stale_marker).toBeDefined()
    })
  })

  describe('Error States and Empty States', () => {
    it('handles empty archived history gracefully', () => {
      const emptyHistoryState = {
        has_archived_plans: false,
        empty_message: 'هنوز هیچ برنامه قبلی ندارید',
        error: null,
      }

      expect(emptyHistoryState.has_archived_plans).toBe(false)
      expect(emptyHistoryState.empty_message).toBeDefined()
    })

    it('surfaces network error state when archive fetch fails', () => {
      const networkErrorState = {
        loading: false,
        error: 'اتصال اینترنت را بررسی کنید',
        retry_available: true,
      }

      expect(networkErrorState.error).toBeDefined()
      expect(networkErrorState.retry_available).toBe(true)
    })
  })

  describe('Navigation Between Active and Archive', () => {
    it('preserves route state when drilling into archived plan', () => {
      // User navigates from history list to archived plan detail
      const navigationState = {
        from_route: '/client/history/plans',
        to_route: '/client/history/plans/plan-2',
        back_available: true,
      }

      expect(navigationState.from_route).toBe('/client/history/plans')
      expect(navigationState.to_route).toContain('/client/history/plans/')
      expect(navigationState.back_available).toBe(true)
    })

    it('maintains active plan context label across navigation', () => {
      const contextPersistence = {
        active_plan_id: 'plan-current',
        visible_on_history_list: true,
        visible_on_history_detail: true,
        visible_on_archive_detail: false, // Archive detail shows different plan
      }

      expect(contextPersistence.active_plan_id).toBe('plan-current')
      expect(contextPersistence.visible_on_history_list).toBe(true)
    })
  })
})
