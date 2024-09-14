package queue

import (
	"context"
	"fmt"

	"github.com/iamsorryprincess/project-layout/internal/domain"
)

type ClickConsumer struct {
}

func NewClickConsumer() *ClickConsumer {
	return &ClickConsumer{}
}

func (c *ClickConsumer) Consume(ctx context.Context, clicks []domain.Click) error {
	fmt.Println(clicks)
	return nil
}
