package nutrition

import "Calsora/internal/models"

const (
	Sedentary  = 1.2
	Light      = 1.375
	Moderate   = 1.55
	Active     = 1.725
	VeryActive = 1.9
)

const (
	calPerKgFat    = 7700
	calPerKgMuscle = 3200
)

func CalculateNorm(profile *models.UserProfile, goal *models.UserGoal, age int) *models.NutritionTarget {
	dailyIntake := &models.NutritionTarget{}
	var bmr float64

	if profile.Sex == "M" {
		bmr = (10.0 * float64(profile.Weight)) + (6.5 * float64(profile.Height)) - (5 * float64(age)) + 5
	} else {
		bmr = (10.0 * float64(profile.Weight)) + (6.5 * float64(profile.Height)) - (5 * float64(age)) - 161
	}

	var tdee float64
	switch profile.Activity {
	case 1:
		tdee = bmr * Sedentary
	case 2:
		tdee = bmr * Light
	case 3:
		tdee = bmr * Moderate
	case 4:
		tdee = bmr * Active
	case 5:
		tdee = bmr * VeryActive
	}

	switch goal.Type {
	case "lose":
		dailyIntake.Cal = int(tdee - ((calPerKgFat * goal.WeeklyRate) / 7))
		dailyIntake.Protein = (float64(dailyIntake.Cal) * 0.22) / 4
		dailyIntake.Carbs = (float64(dailyIntake.Cal) * 0.51) / 4
		dailyIntake.Fats = (float64(dailyIntake.Cal) * 0.27) / 9
	case "gain":
		dailyIntake.Cal = int(tdee + ((calPerKgMuscle * goal.WeeklyRate) / 7))
		dailyIntake.Protein = (float64(dailyIntake.Cal) * 0.22) / 4
		dailyIntake.Carbs = (float64(dailyIntake.Cal) * 0.53) / 4
		dailyIntake.Fats = (float64(dailyIntake.Cal) * 0.25) / 9
	case "maintain":
		dailyIntake.Cal = int(tdee)
		dailyIntake.Protein = (float64(dailyIntake.Cal) * 0.15) / 4
		dailyIntake.Carbs = (float64(dailyIntake.Cal) * 0.55) / 4
		dailyIntake.Fats = (float64(dailyIntake.Cal) * 0.30) / 9
	}

	return dailyIntake
}
