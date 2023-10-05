package input_enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type RetentionGroupEvent int32

const (
	// RETENTION_GROUP_EVENT_UNSPECIFIED.
	RetentionGroupEventUnspecified RetentionGroupEvent = iota
	// RETENTION_GROUP_EVENT_INFLOW represents the retention for the event
	// when the end user has the first interaction with our chatbot.
	RetentionGroupEventInflow
	// RETENTION_GROUP_EVENT_TAGGING represents the retention for the event
	// when the end user submit some data that we store in the EndUserProfile (common.EndUser).
	RetentionGroupEventTagging
	// RETENTION_GROUP_EVENT_CONVERSION represents the retention for the event
	// when the end user makes a CV by using our system.
	RetentionGroupEventConversion
)

// ToProto converts the RetentionGroupEvent to Protobuf version.
func (r RetentionGroupEvent) ToProto() delivery_settings_entities.RetentionGroupEvent {
	switch r {
	case RetentionGroupEventInflow:
		return delivery_settings_entities.RetentionGroupEvent_RETENTION_GROUP_EVENT_INFLOW
	case RetentionGroupEventTagging:
		return delivery_settings_entities.RetentionGroupEvent_RETENTION_GROUP_EVENT_TAGGING
	case RetentionGroupEventConversion:
		return delivery_settings_entities.RetentionGroupEvent_RETENTION_GROUP_EVENT_CONVERSION
	default:
		return delivery_settings_entities.RetentionGroupEvent_RETENTION_GROUP_EVENT_UNSPECIFIED
	}
}

// ProtoToRetentionGroupEvent converts from Protobuf version to the RetentionGroupEvent.
func ProtoToRetentionGroupEvent(r delivery_settings_entities.RetentionGroupEvent) RetentionGroupEvent {
	switch r {
	case delivery_settings_entities.RetentionGroupEvent_RETENTION_GROUP_EVENT_INFLOW:
		return RetentionGroupEventInflow
	case delivery_settings_entities.RetentionGroupEvent_RETENTION_GROUP_EVENT_TAGGING:
		return RetentionGroupEventTagging
	case delivery_settings_entities.RetentionGroupEvent_RETENTION_GROUP_EVENT_CONVERSION:
		return RetentionGroupEventConversion
	default:
		return RetentionGroupEventUnspecified
	}
}

func TestRetentionGroupEvent_Convert(t *testing.T) {
	type want struct {
		args      models.RetentionGroupEvent
		wantProto delivery_settings_entities.RetentionGroupEvent
	}

	type Context struct {
		testData *want
	}

	contextInitiateFunction := func(t *testing.T) *Context {
		return &Context{}
	}

	gt.Begin(t,
		contextInitiateFunction,
		gt.Run("Convert from model to Proto and then convert back to model", func(t *testing.T, ctx *Context) {
			// Convert from model to Proto.
			gotProto := ctx.testData.args.ToProto()
			assert.Equal(t, ctx.testData.wantProto, gotProto)

			// Then convert from Proto back to model
			gotModel := models.ProtoToRetentionGroupEvent(gotProto)
			assert.Equal(t, ctx.testData.args, gotModel)
		}).
			Using("given RetentionGroupEventUnspecified value", func(t *testing.T, ctx *Context) {
				ctx.testData = &want{
					args:      models.RetentionGroupEventUnspecified,
					wantProto: delivery_settings_entities.RetentionGroupEvent_RETENTION_GROUP_EVENT_UNSPECIFIED,
				}
			}).
			Using("given RetentionGroupEventInflow value", func(t *testing.T, ctx *Context) {
				ctx.testData = &want{
					args:      models.RetentionGroupEventInflow,
					wantProto: delivery_settings_entities.RetentionGroupEvent_RETENTION_GROUP_EVENT_INFLOW,
				}
			}).
			Using("given RetentionGroupEventTagging value", func(t *testing.T, ctx *Context) {
				ctx.testData = &want{
					args:      models.RetentionGroupEventTagging,
					wantProto: delivery_settings_entities.RetentionGroupEvent_RETENTION_GROUP_EVENT_TAGGING,
				}
			}).
			Using("given RetentionGroupEventConversion value", func(t *testing.T, ctx *Context) {
				ctx.testData = &want{
					args:      models.RetentionGroupEventConversion,
					wantProto: delivery_settings_entities.RetentionGroupEvent_RETENTION_GROUP_EVENT_CONVERSION,
				}
			}),
	)
}
