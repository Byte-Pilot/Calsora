package nutrition

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go/v3"
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

func GenerateSchema[T any]() interface{} {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}

var MealAnalysisResponseSchema = GenerateSchema[MealAnalysis]()

func (g *GPTClient) AnalyzeMeal(ctx context.Context, text, image string) (*MealAnalysis, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "MealAnalysis",
		Description: openai.String("Schema for meal analysis response"),
		Schema:      MealAnalysisResponseSchema,
		Strict:      openai.Bool(true),
	}

	content := []responses.ResponseInputContentUnionParam{}

	content = append(content, responses.ResponseInputContentUnionParam{
		OfInputText: &responses.ResponseInputTextParam{
			Type: "input_text",
			Text: Prompt,
		},
	})

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
				Type: "input_image",
				//ImageURL: image,
			},
		})
	}
	g.client.Embeddings.New()
	resp, err := g.client.Responses.New(ctx, responses.ResponseNewParams{
		Model: openai.ChatModelGPT5_2,
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: responses.ResponseInputParam{
				responses.ResponseInputItemUnionParam{
					content,
					"user",
				},
			},
		},
		ResponseFormat: &responses.ResponseFormatJSONSchemaParam{
			Type: "json_schema",
			JSONSchema: responses.ResponseFormatJSONSchemaJSONSchemaParam{
				Name:   "meal_analysis",
				Schema: schemaParam,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Output) == 0 {
		return nil, errors.New("empty response from GPT")
	}

	var result MealAnalysis

	for _, item := range resp.Output {
		if item.Type == "output_text" {
			if err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result); err != nil {
				return nil, err
			}
			return &result, nil
		}
	}
}
