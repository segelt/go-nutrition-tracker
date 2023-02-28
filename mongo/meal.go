package mongo

import (
	"context"
	"errors"
	"fmt"
	"nutritiontracker/resource"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ resource.MealService = (*MealService)(nil)

const (
	EXACT = iota
	LTE
	GTE
)

type MealService struct {
	db *DB
}

func NewMealService(db *DB) *MealService {
	return &MealService{db: db}
}

func (s *MealService) FindMealById(id *string) (*resource.Meal, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.db.operationTimeout)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(*id)
	if err != nil {
		return nil, errors.New("invalid id")
	}

	coll := s.db.client.Database("nutrition-tracker").Collection("items")
	filter := bson.D{{"_id", objectId}}
	meal := resource.Meal{}
	if err = coll.FindOne(ctx, filter).Decode(&meal); err != nil {
		return nil, err
	}

	return &meal, nil
}

func (s *MealService) ListMeals(filter resource.MealFilter) ([]resource.Meal, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.db.operationTimeout)
	defer cancel()

	coll := s.db.client.Database("nutrition-tracker").Collection("items")

	sortBSON, filterBSON := toBson(&filter)
	opts := options.Find().SetSort(sortBSON)

	if filter.Offset != nil {
		opts.SetSkip(int64(*filter.Offset))
	}
	if filter.Limit != nil {
		opts.SetLimit(int64(*filter.Limit))
	}

	cursor, err := coll.Find(ctx, filterBSON, opts)
	if err != nil {
		return nil, err
	}

	var results []resource.Meal
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
func (s *MealService) CreateMeal(create *resource.MealCreate) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.db.operationTimeout)
	defer cancel()

	coll := s.db.client.Database("nutrition-tracker").Collection("items")

	newMeal := resource.Meal{
		ID:   primitive.NewObjectID(),
		Name: *create.Name,
		NutritionalValue: resource.NutritionalValue{
			BaseWeight:     *create.NutritionalValue.BaseWeight,
			Calories:       *create.NutritionalValue.Calories,
			Protein:        *create.NutritionalValue.Protein,
			Fat:            *create.NutritionalValue.Fat,
			Carbonhydrates: *create.NutritionalValue.Carbonhydrates,
			Fiber:          *create.NutritionalValue.Fiber,
			Sugar:          *create.NutritionalValue.Sugar,
		},
	}
	_, err := coll.InsertOne(ctx, newMeal)
	return err
}
func (s *MealService) UpdateMeal(upd *resource.MealUpdate) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.db.operationTimeout)
	defer cancel()

	coll := s.db.client.Database("nutrition-tracker").Collection("items")

	var update []bson.E
	if upd.Name != nil && len(*upd.Name) > 0 {
		update = append(update, bson.E{"name", *upd.Name})
	}

	if upd.NutritionalValue != nil {
		if upd.NutritionalValue.BaseWeight != nil {
			update = append(update, bson.E{"nutritionalValue.baseWeight", *upd.NutritionalValue.BaseWeight})
		}

		if upd.NutritionalValue.Calories != nil {
			update = append(update, bson.E{"nutritionalValue.calories", *upd.NutritionalValue.Calories})
		}

		if upd.NutritionalValue.Protein != nil {
			update = append(update, bson.E{"nutritionalValue.protein", *upd.NutritionalValue.Protein})
		}

		if upd.NutritionalValue.Fat != nil {
			update = append(update, bson.E{"nutritionalValue.fat", *upd.NutritionalValue.Fat})
		}

		if upd.NutritionalValue.Carbonhydrates != nil {
			update = append(update, bson.E{"nutritionalValue.carbonhydrates", *upd.NutritionalValue.Carbonhydrates})
		}

		if upd.NutritionalValue.Fiber != nil {
			update = append(update, bson.E{"nutritionalValue.fiber", *upd.NutritionalValue.Fiber})
		}

		if upd.NutritionalValue.Sugar != nil {
			update = append(update, bson.E{"nutritionalValue.sugar", *upd.NutritionalValue.Sugar})
		}
	}

	if len(update) == 0 {
		return errors.New("no field to update.")
	}

	updateCommand := bson.D{{"$set", update}}
	id, _ := primitive.ObjectIDFromHex(upd.ID)
	filter := bson.D{{"_id", id}}

	result, err := coll.UpdateOne(ctx, filter, updateCommand)

	if err != nil {
		return err
	}
	if result.ModifiedCount != 1 {
		return fmt.Errorf("invalid amount of documents updated %d", result.ModifiedCount)
	}

	return nil
}
func (s *MealService) DeleteMeal(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.db.operationTimeout)
	defer cancel()

	coll := s.db.client.Database("nutrition-tracker").Collection("items")

	itemId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", itemId}}

	result, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount != 1 {
		return fmt.Errorf("mismatching amount of deleted items. %d items were deleted. Expected 1", result.DeletedCount)
	}
	return nil
}

func toBson(f *resource.MealFilter) (*bson.D, *bson.M) {
	sorts := bson.D{}
	if f.Sorts == nil {
		sorts = nil
	} else {
		for _, sortItem := range *f.Sorts {
			sorts = append(sorts, bson.E{*sortItem.Column, sortItem.Direction})
		}
	}

	filters := bson.M{}
	if f.Filters == nil {
		filters = nil
	} else {
		for _, filterItem := range *f.Filters {
			// TODO: Finding a better way to capture values that can either be float or string.
			floatValue, err := strconv.ParseFloat(*filterItem.Value, 32)

			if err == nil {
				switch filterItem.FilterType {
				case EXACT:
					filters[*filterItem.Column] = floatValue
				case LTE:
					filters[*filterItem.Column] = bson.M{"$lte": floatValue}
				case GTE:
					filters[*filterItem.Column] = bson.M{"$gte": floatValue}
				}
			} else {
				switch filterItem.FilterType {
				case EXACT:
					filters[*filterItem.Column] = filterItem.Value
				case LTE:
					filters[*filterItem.Column] = bson.M{"$lte": filterItem.Value}
				case GTE:
					filters[*filterItem.Column] = bson.M{"$gte": filterItem.Value}
				}
			}
		}
	}

	return &sorts, &filters
}
