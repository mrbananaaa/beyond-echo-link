package redis

import "context"

type Message struct {
	Channel string
	Payload []byte
}

type PubSub interface {
	Publish(ctx context.Context, channel string, message []byte) error
	Subscribe(ctx context.Context, channel string) (<-chan Message, error)
}

type RedisPubSub struct {
	redis *Redis
}

func NewPubSub(r *Redis) *RedisPubSub {
	return &RedisPubSub{
		redis: r,
	}
}

func (p *RedisPubSub) Publish(
	ctx context.Context,
	channel string,
	message []byte,
) error {
	ctx, cancel := p.redis.withTimeout(ctx)
	defer cancel()

	return p.redis.client.Publish(ctx, channel, message).Err()
}

func (p *RedisPubSub) Subscribe(
	ctx context.Context,
	channel string,
) (<-chan Message, error) {

	pubsub := p.redis.client.Subscribe(ctx, channel)

	ch := make(chan Message)

	go func() {
		defer close(ch)
		defer pubsub.Close()

		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-pubsub.Channel():
				if !ok {
					return
				}

				ch <- Message{
					Channel: msg.Channel,
					Payload: []byte(msg.Payload),
				}
			}
		}
	}()

	return ch, nil
}
