package model

import (
	"database/sql/driver"
	"fmt"
)

type FoodCategory string

const (
	FoodCategoryBreakfast  FoodCategory = "breakfast"
	FoodCategoryLunch      FoodCategory = "lunch"
	FoodCategoryDinner     FoodCategory = "dinner"
	FoodCategorySnack      FoodCategory = "snack"
	FoodCategoryFruit      FoodCategory = "fruit"
	FoodCategoryBeverage   FoodCategory = "beverage"
	FoodCategorySupplement FoodCategory = "supplement"
	FoodCategoryOther      FoodCategory = "other"
)

func (e *FoodCategory) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = FoodCategory(s)
	case string:
		*e = FoodCategory(s)
	default:
		return fmt.Errorf("unsupported scan type for FoodCategory: %T", src)
	}
	return nil
}

func (e FoodCategory) Value() (driver.Value, error) {
	return string(e), nil
}

type MeasurementUnit string

const (
	MeasurementUnitGram       MeasurementUnit = "gram"
	MeasurementUnitKg         MeasurementUnit = "kg"
	MeasurementUnitTablespoon MeasurementUnit = "tablespoon"
	MeasurementUnitTeaspoon   MeasurementUnit = "teaspoon"
	MeasurementUnitCup        MeasurementUnit = "cup"
	MeasurementUnitPiece      MeasurementUnit = "piece"
	MeasurementUnitSlice      MeasurementUnit = "slice"
	MeasurementUnitPalm       MeasurementUnit = "palm"
	MeasurementUnitMatchbox   MeasurementUnit = "matchbox"
	MeasurementUnitBowl       MeasurementUnit = "bowl"
	MeasurementUnitMl         MeasurementUnit = "ml"
	MeasurementUnitLiter      MeasurementUnit = "liter"
)

func (e *MeasurementUnit) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = MeasurementUnit(s)
	case string:
		*e = MeasurementUnit(s)
	default:
		return fmt.Errorf("unsupported scan type for MeasurementUnit: %T", src)
	}
	return nil
}

func (e MeasurementUnit) Value() (driver.Value, error) {
	return string(e), nil
}
