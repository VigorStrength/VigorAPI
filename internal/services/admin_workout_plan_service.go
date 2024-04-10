package services

import (
	"context"
	"fmt"

	"github.com/GhostDrew11/vigor-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrWorkoutPlanAlreadyExists = fmt.Errorf("workout plan already exists")
	ErrWorkoutPlanNotFound = fmt.Errorf("workout plan not found")
)

func (as *AdminService) GetWorkoutPlanByID(ctx context.Context, workoutPlanID primitive.ObjectID) (models.WorkoutPlan, error) {
	// Get the workout plan collection
	workoutPlanCollection := as.database.Collection("workoutPlans")

	// Find the workout plan by ID
	filter := bson.M{"_id": workoutPlanID}
	var workoutPlan models.WorkoutPlan
	err := workoutPlanCollection.FindOne(ctx, filter).Decode(&workoutPlan)
	if err != nil {
		return models.WorkoutPlan{}, fmt.Errorf("error finding workout plan: %w", err)
	}

	return workoutPlan, nil
}

func (as *AdminService) GetWorkoutPlans(ctx context.Context) ([]models.WorkoutPlan, error) {
	workoutCollection := as.database.Collection("workoutPlans")

	cursor, err := workoutCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error finding workout plans: %w", err)
	}
	defer cursor.Close(ctx)

	var workoutPlans []models.WorkoutPlan
	if err := cursor.All(ctx, &workoutPlans); err != nil {
		return nil, fmt.Errorf("error decoding workout plans: %w", err)
	}

	return workoutPlans, nil
}

func(as *AdminService) SearchWorkoutPlansByName(ctx context.Context, name string) ([]models.WorkoutPlan, error) {
    workoutPlanCollection := as.database.Collection("workoutPlans")

    filter := bson.M{"name": primitive.Regex{Pattern: name, Options: "i"}}
    cursor, err := workoutPlanCollection.Find(ctx, filter)
    if err != nil {
        return nil, fmt.Errorf("error finding workout plans by name: %w", err)
    }
    defer cursor.Close(ctx)

    var workoutPlans []models.WorkoutPlan
    if err := cursor.All(ctx, &workoutPlans); err != nil {
        return nil, fmt.Errorf("error decoding workout plans: %w", err)
    }

    return workoutPlans, nil
}

func (as *AdminService) CreateWorkoutPlan(ctx context.Context, workoutPlanInput models.WorkoutPlan) error {
	// Get the workout plan collection
	workoutPlanCollection := as.database.Collection("workoutPlans")

	// check if the new workout plan already exists
	filter := bson.M{"name": workoutPlanInput.Name}
	count, err := workoutPlanCollection.CountDocuments(ctx, filter)
	if err != nil {
		return fmt.Errorf("error checking if workout plan already exists: %w", err)
	}

	if count > 0 {
		return ErrWorkoutPlanAlreadyExists
	}

	workoutPlanInput.ID = primitive.NewObjectID()
	for i := range workoutPlanInput.Weeks {
		workoutPlanInput.Weeks[i].ID = primitive.NewObjectID()
		for j := range workoutPlanInput.Weeks[i].Days {
			workoutPlanInput.Weeks[i].Days[j].ID = primitive.NewObjectID()
			for k := range workoutPlanInput.Weeks[i].Days[j].WarmUps {
				workoutPlanInput.Weeks[i].Days[j].WarmUps[k].ID = primitive.NewObjectID()
			}
			for l := range workoutPlanInput.Weeks[i].Days[j].Workouts {
				workoutPlanInput.Weeks[i].Days[j].Workouts[l].ID = primitive.NewObjectID()
			}
			for m := range workoutPlanInput.Weeks[i].Days[j].CoolDowns {
				workoutPlanInput.Weeks[i].Days[j].CoolDowns[m].ID = primitive.NewObjectID()
			}
		}
	}

	// Insert the new workout plan
	_, err = workoutPlanCollection.InsertOne(ctx, workoutPlanInput)
	if err != nil {
		return fmt.Errorf("error inserting workout plan: %w", err)
	}

	return nil
}

func (as *AdminService) UpdateWorkoutPlan(ctx context.Context, workoutPlanID primitive.ObjectID, updateInput models.WorkoutPlanInput) error {
	// Get the workout plan collection 
	workoutPlanCollection := as.database.Collection("workoutPlans")

	var existingPlan models.WorkoutPlan
    if err := workoutPlanCollection.FindOne(ctx, bson.M{"_id": workoutPlanID}).Decode(&existingPlan); err != nil {
        return fmt.Errorf("error fetching existing workout plan: %w", err)
    }

	updatedDoc := mergeUpdatesIntoExistingPlan(existingPlan, updateInput)

	// Use $set to only update the provided fields
	update := bson.M{"$set": updatedDoc}
	filter := bson.M{"_id": workoutPlanID}
	result, err := workoutPlanCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error updating workout plan: %w", err)
	}

	if result.MatchedCount == 0 {
		return ErrWorkoutPlanNotFound
	}

	return nil
}

func (as *AdminService) DeleteWorkoutPlan(ctx context.Context, workoutPlanID primitive.ObjectID) error {
	// Get the workout plan collection
	workoutPlanCollection := as.database.Collection("workoutPlans")
	filter := bson.M{"_id": workoutPlanID}

	// Delete the workout plan
	result, err := workoutPlanCollection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error deleting workout plan: %w", err)
	}

	if result.DeletedCount == 0 {
		return ErrWorkoutPlanNotFound
	}

	return nil
}

func mergeUpdatesIntoExistingPlan(existingPlan models.WorkoutPlan, updateInput models.WorkoutPlanInput) models.WorkoutPlan {
    // First, apply top-level updates from updateInput to existingPlan
    if updateInput.Name != nil {
        existingPlan.Name = *updateInput.Name
    }
    if updateInput.Duration != nil {
        existingPlan.Duration = *updateInput.Duration
    }

    // Initialize a map for quick lookup of existing weeks by weekNumber
    existingWeeksMap := make(map[int]*models.WorkoutWeek)
    for i := range existingPlan.Weeks {
        weekNumber := existingPlan.Weeks[i].WeekNumber
        existingWeeksMap[weekNumber] = &existingPlan.Weeks[i]
    }

    // Iterate through updateInput.Weeks, merging changes into existing weeks or creating new weeks
    updatedWeeks := []models.WorkoutWeek{}
    if updateInput.Weeks != nil {
        for _, weekInput := range *updateInput.Weeks {
            if existingWeek, exists := existingWeeksMap[*weekInput.WeekNumber]; exists {
                // Merge updates for existing week
                updatedWeek := mergeWeekUpdates(*existingWeek, weekInput)
                updatedWeeks = append(updatedWeeks, updatedWeek)
                // Remove the week from the map to track which weeks have been processed
                delete(existingWeeksMap, *weekInput.WeekNumber)
            } else {
                // Create new week since it doesn't exist in the existingPlan
                newWeek := createNewWeekFromInput(weekInput)
                updatedWeeks = append(updatedWeeks, newWeek)
            }
        }
    }

    // Add any remaining weeks from existingPlan that were not in updateInput
    for _, week := range existingWeeksMap {
        updatedWeeks = append(updatedWeeks, *week)
    }

    // Update the weeks in existingPlan with the merged weeks
    existingPlan.Weeks = updatedWeeks

    return existingPlan
}

func mergeWeekUpdates(existingWeek models.WorkoutWeek, weekInput models.WorkoutWeekInput) models.WorkoutWeek {
    // Assuming that there are no top-level fields to update in WorkoutWeek itself from WorkoutWeekInput
    // Update days within the week
    if weekInput.Days != nil {
        updatedDays := []models.WorkoutDay{}
        for _, dayInput := range *weekInput.Days {
            // Create new days for now; adjust based on your logic to match existing days if necessary
            newDay := createNewDayFromInput(dayInput)
            updatedDays = append(updatedDays, newDay)
        }
        existingWeek.Days = updatedDays
    }
    return existingWeek
}

func createNewWeekFromInput(weekInput models.WorkoutWeekInput) models.WorkoutWeek {
    newWeek := models.WorkoutWeek{
        ID:         primitive.NewObjectID(), // Assign a new ID
        WeekNumber: *weekInput.WeekNumber,
        Days:       []models.WorkoutDay{}, // Initialize Days slice for new weeks
    }
    // Create new days for the new week
    if weekInput.Days != nil {
        for _, dayInput := range *weekInput.Days {
            newDay := createNewDayFromInput(dayInput)
            newWeek.Days = append(newWeek.Days, newDay)
        }
    }
    return newWeek
}

func createNewDayFromInput(dayInput models.WorkoutDayInput) models.WorkoutDay {
    newDay := models.WorkoutDay{
        ID: primitive.NewObjectID(), // Assign a new ID for the day
    }

    // If WorkoutTimeRange is provided, use it; otherwise, initialize to a default range
    if dayInput.WorkoutTimeRange != nil {
        newDay.WorkoutTimeRange = *dayInput.WorkoutTimeRange
    } else {
        newDay.WorkoutTimeRange = [2]int{0, 0} // Example default value, adjust as needed
    }

    // Convert CircuitInputs to Circuits
    newDay.WarmUps = convertCircuitInputsToCircuits(dayInput.WarmUps)
    newDay.Workouts = convertCircuitInputsToCircuits(dayInput.Workouts)
    newDay.CoolDowns = convertCircuitInputsToCircuits(dayInput.CoolDowns)

    return newDay
}

func convertCircuitInputsToCircuits(circuitInputs *[]models.CircuitInput) []models.Circuit {
    circuits := []models.Circuit{}
    if circuitInputs == nil {
        return circuits // Return an empty slice if there's no input
    }

    for _, input := range *circuitInputs {
        newCircuit := models.Circuit{
            ID: primitive.NewObjectID(), // Assign a new ID for each circuit
        }

        // Set ExerciseIDs
        if input.ExerciseIDs != nil {
            newCircuit.ExerciseIDs = *input.ExerciseIDs
        } else {
            newCircuit.ExerciseIDs = []primitive.ObjectID{} // Ensure a valid slice is set
        }

        // Set RestTime
        if input.RestTime != nil {
            newCircuit.RestTime = input.RestTime
        } // else it remains nil, which is valid per the Circuit struct definition

        // Set ProposedLaps, defaulting to 1 if not provided
        if input.ProposedLaps != nil {
            newCircuit.ProposedLaps = *input.ProposedLaps
        } else {
            newCircuit.ProposedLaps = 1
        }

        circuits = append(circuits, newCircuit)
    }

    return circuits
}