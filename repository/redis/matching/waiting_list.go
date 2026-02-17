package matching

import (
	"context"
	"fmt"
	"game_app/entity"
	"game_app/pkg/richerror"
	"time"

	"github.com/redis/go-redis/v9"
)

const WaitingListPrefix = "waitinglist"

func (d DB) AddToWaitingList(userID uint, category entity.Category) error {
	const op = richerror.Op("matching.AddToWaitingList")
	var ctx = context.Background()
	zKey := fmt.Sprintf("%s:%s", WaitingListPrefix, category)
	_, err := d.adapter.Client.ZAdd(ctx, zKey, redis.Z{
		Score:  float64(time.Now().UnixMicro()),
		Member: fmt.Sprintf("%d", userID),
	}).Result()

	if err != nil {
		return richerror.New(op).SetWrappedError(err).SetKind(richerror.KindUnexpected)
	}

	return nil
}
