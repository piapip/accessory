// Code generated by accessory; DO NOT EDIT.

package input

import (
	"time"
)

// GetTimeTransition returns the DeliveryTiming's TimeTransition.
func (d *DeliveryTiming) GetTimeTransition() int32 {
	if d == nil {
		return 0
	}

	return d.TimeTransition
}

// GetNumber returns the DeliveryTiming's Number.
func (d *DeliveryTiming) GetNumber() *int64 {
	if d == nil {
		return nil
	}

	return d.Number
}

// GetTimeUnit returns the DeliveryTiming's TimeUnit.
func (d *DeliveryTiming) GetTimeUnit() int32 {
	if d == nil {
		return 0
	}

	return d.TimeUnit
}

// GetDeliveryTime returns the DeliveryTiming's DeliveryTime.
func (d *DeliveryTiming) GetDeliveryTime() *time.Time {
	if d == nil {
		return nil
	}

	return d.DeliveryTime
}

func TestDeliveryTiming_GetFunctions(t *testing.T) {
	type want struct {
		args               *models.DeliveryTiming
		wantTimeTransition int32
		wantNumber         *int64
		wantTimeUnit       int32
		wantDeliveryTime   *time.Time
		wantProto          *replaceMe.DeliveryTiming
	}

	type Context struct {
		testData *want
	}

	contextInitiateFunction := func(t *testing.T) *Context {
		return &Context{}
	}

	gt.Begin(t,
		contextInitiateFunction,
		gt.Run("Get functions return proper value", func(t *testing.T, ctx *Context) {
			// GET functions
			gotTimeTransition := ctx.testData.args.GetTimeTransition()
			assert.Equal(t, ctx.testData.wantTimeTransition, gotTimeTransition)

			gotNumber := ctx.testData.args.GetNumber()
			assert.Equal(t, ctx.testData.wantNumber, gotNumber)

			gotTimeUnit := ctx.testData.args.GetTimeUnit()
			assert.Equal(t, ctx.testData.wantTimeUnit, gotTimeUnit)

			gotDeliveryTime := ctx.testData.args.GetDeliveryTime()
			assert.Equal(t, ctx.testData.wantDeliveryTime, gotDeliveryTime)

			// Convert from models to Proto.
			gotProto := ctx.testData.args.ToProto()
			assert.Equal(t, ctx.testData.wantProto, gotProto)

			// Then convert from Proto back to model
			gotModel := models.ProtoToRetentionAvoid(gotProto)
			assert.Equal(t, ctx.testData.args, gotModel)
		}).
			Using("given nil value", func(t *testing.T, ctx *Context) {
				ctx.testData = &want{
					args:               nil,
					wantTimeTransition: 0,
					wantNumber:         nil,
					wantTimeUnit:       0,
					wantDeliveryTime:   nil,
					wantProto:          nil,
				}
			}).
			Using("given empty value", func(t *testing.T, ctx *Context) {
				ctx.testData = &want{
					args:               &models.DeliveryTiming{},
					wantTimeTransition: 0,
					wantNumber:         nil,
					wantTimeUnit:       0,
					wantDeliveryTime:   nil,
					wantProto:          &replaceMe.DeliveryTiming{},
				}
			}).
			Using("given NON nil value", func(t *testing.T, ctx *Context) {
				ctx.testData = &want{}
			}),
	)
}
