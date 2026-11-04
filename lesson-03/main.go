package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"go-agent/consts"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

func main() {
	// question := "给我讲讲os.Args"
	// if len(os.Args) > 1 {
	// 	question = strings.Join(os.Args[1:], " ")
	// }
	chat, err := NewCoach()
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx := context.Background()
	// answer, usage, err := chat.Ask(ctx, question)
	err = chat.RunChat(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println("问题；", question, "模型回答:", answer)
	// fmt.Println("使用情况:", usage)
}

type Coach struct {
	client *openai.Client
	model  string
}

func NewCoach() (*Coach, error) {
	apiKey := strings.TrimSpace(os.Getenv("DEEPSEEK_API_KEY"))
	if apiKey == "" {
		return nil, errors.New("DEEPSEEK_API_KEY is required")
	}
	cfg := openai.DefaultConfig(apiKey)
	cfg.BaseURL = "https://api.deepseek.com"

	return &Coach{client: openai.NewClientWithConfig(cfg),
		model: consts.DeepSeekModelFlash}, nil
}

func (c *Coach) RunChat(ctx context.Context) error {
	message := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "你是一个Go后端的老师，回答要适合初学者，先给结论再解释原因",
		},
	}
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("开始对话 ，输入exit退出：")
	for {
		fmt.Print("\n你：")
		if !scanner.Scan() {
			return nil
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		if input == "exit" {
			fmt.Println("再见！")
			return nil
		}
		message = append(message, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: input,
		})
		resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			Model:               c.model,
			Messages:            message,
			MaxCompletionTokens: 100,
			Temperature:         0.4,
		})
		if err != nil {

			fmt.Println(err)
			message = message[:len(message)-1]
			continue
		}
		reply := resp.Choices[0].Message.Content
		fmt.Println("\nAI助手：", reply)
		message = append(message, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: reply,
		})

	}
}

func (c *Coach) Ask(ctx context.Context, question string) (string, openai.Usage, error) {
	resp, err := c.client.CreateChatCompletion(ctx,
		openai.ChatCompletionRequest{
			Model: c.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "你是一个Go后端的老师，回答要适合初学者，先给结论再解释原因",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: question,
				},
			},
			MaxTokens:           1024,
			MaxCompletionTokens: 0,
			Temperature:         0.3,
			TopP:                0,
		},
	)
	if err != nil {
		return "", openai.Usage{}, err
	}
	if len(resp.Choices) == 0 {
		return "", openai.Usage{}, errors.New("no choices returned")
	}
	return resp.Choices[0].Message.Content, resp.Usage, nil
}
