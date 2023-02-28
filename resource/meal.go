package resource

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Meal struct {
	ID               primitive.ObjectID `json:"id" bson:"_id"`
	Name             string             `json:"name" bson:"name"`
	NutritionalValue NutritionalValue   `json:"nutritionalValue" bson:"nutritionalValue"`
}

type NutritionalValue struct {
	BaseWeight     float32 `json:"baseWeight" bson:"baseWeight"`
	Calories       float32 `json:"calories" bson:"calories"`
	Protein        float32 `json:"protein" bson:"protein"`
	Fat            float32 `json:"fat" bson:"fat"`
	Carbonhydrates float32 `json:"carbonhydrates" bson:"carbonhydrates"`
	Fiber          float32 `json:"fiber" bson:"fiber"`
	Sugar          float32 `json:"sugar" bson:"sugar"`
}

type MealService interface {
	FindMealById(id *string) (*Meal, error)
	ListMeals(filter MealFilter) ([]Meal, error)
	CreateMeal(create *MealCreate) error
	UpdateMeal(upd *MealUpdate) error
	DeleteMeal(id string) error
}

type MealCreate struct {
	Name             *string                `json:"name" binding:"required" bson:"name"`
	NutritionalValue NutritionalValueCreate `json:"nutritionalValue" binding:"required" bson:"nutritionalValue"`
}

type NutritionalValueCreate struct {
	BaseWeight     *float32 `json:"baseWeight" binding:"required"`
	Calories       *float32 `json:"calories" binding:"required"`
	Protein        *float32 `json:"protein" binding:"required"`
	Fat            *float32 `json:"fat" binding:"required"`
	Carbonhydrates *float32 `json:"carbonhydrates" binding:"required"`
	Fiber          *float32 `json:"fiber" binding:"required"`
	Sugar          *float32 `json:"sugar" binding:"required"`
}

type MealFilter struct {
	Filters *[]FilterParam `json:"filterParams"`
	Sorts   *[]SortParam   `json:"sortParams"`

	// Restrict to subset of results.
	Offset *int `json:"offset"`
	Limit  *int `json:"limit"`
}

type FilterParam struct {
	Column     *string `json:"filter" binding:"required"`
	Value      *string `json:"value" default:"asc"`
	FilterType int     `json:"filterType"`
}

type SortParam struct {
	Column    *string `json:"column" binding:"required"`
	Direction int     `json:"direction" default:"asc"`
}

type MealUpdate struct {
	ID               string                        `json:"id"`
	Name             *string                       `json:"name"`
	NutritionalValue *NutritionalValueUpdateParams `json:"nutritionalValue"`
}

type NutritionalValueUpdateParams struct {
	BaseWeight     *float32 `json:"baseWeight"`
	Calories       *float32 `json:"calories"`
	Protein        *float32 `json:"protein"`
	Fat            *float32 `json:"fat"`
	Carbonhydrates *float32 `json:"carbonhydrates"`
	Fiber          *float32 `json:"fiber"`
	Sugar          *float32 `json:"sugar"`
}
