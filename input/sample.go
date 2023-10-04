package input

import "time"

type RetentionAvoid struct {
	AvoidWeekend   bool       `bson:"avoid_weekend"`
	AvoidStartTime *time.Time `bson:"avoid_start_time"`
	AvoidEndTime   *time.Time `bson:"avoid_end_time"`
}
