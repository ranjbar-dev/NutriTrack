package dietplan

import "github.com/ranjbar-dev/nutritrack/internal/domain/dietplan/entity"

// computeOptionTotals sums item computed values for one option.
func computeOptionTotals(items []*entity.MealOptionItem) *entity.NutritionalSummary {
	s := &entity.NutritionalSummary{}
	for _, it := range items {
		if it.Computed != nil {
			s.Calories += it.Computed.Calories
			s.Protein += it.Computed.Protein
			s.Carbs += it.Computed.Carbs
			s.Fat += it.Computed.Fat
			s.Fiber += it.Computed.Fiber
		}
	}
	return s
}

// computeMealRange derives min/max NutritionalSummary across a meal's options.
func computeMealRange(options []*entity.MealOption) *entity.NutritionalRange {
	if len(options) == 0 {
		return nil
	}
	first := true
	r := &entity.NutritionalRange{}
	for _, opt := range options {
		if opt.Totals == nil {
			continue
		}
		t := opt.Totals
		if first {
			r.Min = *t
			r.Max = *t
			first = false
			continue
		}
		if t.Calories < r.Min.Calories {
			r.Min.Calories = t.Calories
		}
		if t.Calories > r.Max.Calories {
			r.Max.Calories = t.Calories
		}
		if t.Protein < r.Min.Protein {
			r.Min.Protein = t.Protein
		}
		if t.Protein > r.Max.Protein {
			r.Max.Protein = t.Protein
		}
		if t.Carbs < r.Min.Carbs {
			r.Min.Carbs = t.Carbs
		}
		if t.Carbs > r.Max.Carbs {
			r.Max.Carbs = t.Carbs
		}
		if t.Fat < r.Min.Fat {
			r.Min.Fat = t.Fat
		}
		if t.Fat > r.Max.Fat {
			r.Max.Fat = t.Fat
		}
		if t.Fiber < r.Min.Fiber {
			r.Min.Fiber = t.Fiber
		}
		if t.Fiber > r.Max.Fiber {
			r.Max.Fiber = t.Fiber
		}
	}
	if first {
		return nil
	}
	return r
}

// computeDayRange sums meal min/max ranges for a day.
func computeDayRange(meals []*entity.DietMeal) *entity.NutritionalRange {
	day := &entity.NutritionalRange{}
	for _, m := range meals {
		if m.TotalRange == nil {
			continue
		}
		day.Min.Calories += m.TotalRange.Min.Calories
		day.Min.Protein += m.TotalRange.Min.Protein
		day.Min.Carbs += m.TotalRange.Min.Carbs
		day.Min.Fat += m.TotalRange.Min.Fat
		day.Min.Fiber += m.TotalRange.Min.Fiber
		day.Max.Calories += m.TotalRange.Max.Calories
		day.Max.Protein += m.TotalRange.Max.Protein
		day.Max.Carbs += m.TotalRange.Max.Carbs
		day.Max.Fat += m.TotalRange.Max.Fat
		day.Max.Fiber += m.TotalRange.Max.Fiber
	}
	return day
}
