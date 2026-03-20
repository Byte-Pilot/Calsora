package nutrition

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/invopop/jsonschema"
	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
	"os"
	"time"
)

type GPTClient struct {
	client openai.Client
}

func NewGPTClient() *GPTClient {
	apiKey := os.Getenv("OPENAI_API_KEY")
	return &GPTClient{
		client: openai.NewClient(
			option.WithAPIKey(apiKey),
		)}
}

func GenerateSchema[T any]() map[string]any {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	data, _ := json.Marshal(schema)
	var result map[string]any
	_ = json.Unmarshal(data, &result)
	return result
}

var MealAnalysisResponseSchema = GenerateSchema[MealAnalysis]()

func (g *GPTClient) AnalyzeMeal(text, image string) (*MealAnalysis, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	var content []responses.ResponseInputContentUnionParam

	if text != "" {
		content = append(content, responses.ResponseInputContentUnionParam{
			OfInputText: &responses.ResponseInputTextParam{
				Type: "input_text",
				Text: text,
			},
		})
	}

	if image != "" {
		content = append(content, responses.ResponseInputContentUnionParam{
			OfInputImage: &responses.ResponseInputImageParam{
				Type:     "input_image",
				ImageURL: openai.String("data:image/jpeg;base64," + image),
				Detail:   responses.ResponseInputImageDetailOriginal,
			},
		})
	}

	resp, err := g.client.Responses.New(ctx, responses.ResponseNewParams{
		Model:        openai.ChatModelGPT5ChatLatest,
		Instructions: openai.String(Prompt),
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: []responses.ResponseInputItemUnionParam{
				{
					OfMessage: &responses.EasyInputMessageParam{
						Role: "user",
						Content: responses.EasyInputMessageContentUnionParam{
							OfInputItemContentList: responses.ResponseInputMessageContentListParam(content),
						},
					},
				},
			},
		},
		Text: responses.ResponseTextConfigParam{
			Format: responses.ResponseFormatTextConfigParamOfJSONSchema(
				"meal_analysis",
				MealAnalysisResponseSchema,
			),
		},
	})
	if err != nil {
		return nil, errors.New("GPT request failed")
	}

	if len(resp.Output) == 0 || len(resp.Output[0].Content) == 0 {
		return nil, errors.New("empty response from GPT")
	}

	var result MealAnalysis

	err = json.Unmarshal(
		[]byte(resp.Output[0].Content[0].Text),
		&result,
	)
	if err != nil {
		return nil, errors.New("failed to unmarshal GPT response:" + resp.Output[0].Content[0].Text)
	}

	return &result, nil
}
