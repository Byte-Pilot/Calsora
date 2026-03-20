package nutrition

type MealAnalysis struct {
	Name  string              `json:"name" jsonschema_description:"The name of the analyzed meal or dish"`
	Items []MealAnalysisItems `json:"items" jsonschema_description:"List of food items detected in the meal"`
}

type MealAnalysisItems struct {
	Name       string  `json:"name" jsonschema_description:"Food item name"`
	Grams      int     `json:"grams" jsonschema_description:"Estimated weight in grams"`
	Calories   int     `json:"cal" jsonschema_description:"Total energy value of the meal in kilocalories (kcal)"`
	Protein    float64 `json:"protein" jsonschema_description:"Amount of protein in the meal, measured in grams"`
	Carbs      float64 `json:"carbs" jsonschema_description:"Amount of carbohydrates in the meal, measured in grams"`
	Fats       float64 `json:"fats" jsonschema_description:"Amount of fats in the meal, measured in grams"`
	Confidence float64 `json:"confidence" jsonschema_description:"Confidence level of the analysis accuracy, from 0.0 (low) to 1.0 (high)"`
}
