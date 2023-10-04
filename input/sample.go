package input

import "time"

type DeliveryTiming struct {
	TimeTransition int32      `bson:"time_transition"`
	Number         *int64     `bson:"number"`
	TimeUnit       int32      `bson:"time_unit"`
	DeliveryTime   *time.Time `bson:"delivery_time"`
}
