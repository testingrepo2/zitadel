package projection

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/feature"
	"github.com/zitadel/zitadel/internal/repository/feature/feature_v2"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	InstanceFeatureTable = "projections.instance_features"

	InstanceFeatureInstanceIDCol   = "instance_id"
	InstanceFeatureKeyCol          = "key"
	InstanceFeatureCreationDateCol = "creation_date"
	InstanceFeatureChangeDateCol   = "change_date"
	InstanceFeatureSequenceCol     = "sequence"
	InstanceFeatureValueCol        = "value"
)

type instanceFeatureProjection struct{}

func newInstanceFeatureProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(instanceFeatureProjection))
}

func (*instanceFeatureProjection) Name() string {
	return InstanceFeatureTable
}

func (*instanceFeatureProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(handler.NewTable(
		[]*handler.InitColumn{
			handler.NewColumn(InstanceFeatureInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(InstanceFeatureKeyCol, handler.ColumnTypeText),
			handler.NewColumn(InstanceFeatureCreationDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(InstanceFeatureChangeDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(InstanceFeatureSequenceCol, handler.ColumnTypeInt64),
			handler.NewColumn(InstanceFeatureValueCol, handler.ColumnTypeJSONB),
		},
		handler.NewPrimaryKey(InstanceFeatureInstanceIDCol, InstanceFeatureKeyCol),
	))
}

func (*instanceFeatureProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{{
		Aggregate: feature_v2.AggregateType,
		EventReducers: []handler.EventReducer{
			{
				Event:  feature_v2.InstanceResetEventType,
				Reduce: reduceInstanceResetFeatures,
			},
			{
				Event:  feature_v2.InstanceDefaultLoginInstanceEventType,
				Reduce: reduceInstanceSetFeature[bool],
			},
			{
				Event:  feature_v2.InstanceTriggerIntrospectionProjectionsEventType,
				Reduce: reduceInstanceSetFeature[bool],
			},
			{
				Event:  feature_v2.InstanceLegacyIntrospectionEventType,
				Reduce: reduceInstanceSetFeature[bool],
			},
		},
	}}
}

func featureKeyFromEventType(eventType eventstore.EventType) (string, error) {
	ss := strings.Split(string(eventType), ".")
	if len(ss) != 4 {
		return "", zerrors.ThrowInternalf(nil, "PROJE-Ahs4m", "reduce.wrong.event.type %s", eventType)
	}
	if _, err := feature.FeatureString(ss[2]); err != nil {
		return "", zerrors.ThrowInternalf(err, "PROJE-Boo2i", "reduce.wrong.event.type %s", eventType)
	}
	return ss[2], nil
}

func reduceInstanceSetFeature[T any](event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*feature_v2.SetEvent[T])
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-uPh8O", "reduce.wrong.event.type %T", event)
	}
	key, err := featureKeyFromEventType(e.EventType)
	if err != nil {
		return nil, err
	}
	columns := []handler.Column{
		handler.NewCol(InstanceFeatureInstanceIDCol, e.Aggregate().InstanceID),
		handler.NewCol(InstanceFeatureKeyCol, key),
		handler.NewCol(InstanceFeatureCreationDateCol, handler.OnlySetValueOnInsert(InstanceFeatureTable, e.CreationDate())),
		handler.NewCol(InstanceFeatureChangeDateCol, e.CreationDate()),
		handler.NewCol(InstanceFeatureSequenceCol, e.Sequence()),
		handler.NewCol(InstanceFeatureValueCol, e.Value),
	}
	return handler.NewUpsertStatement(e, columns[0:2], columns), nil
}

func reduceInstanceResetFeatures(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*feature_v2.ResetEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-roo6A", "reduce.wrong.event.type %T", event)
	}
	return handler.NewDeleteStatement(e, []handler.Condition{
		handler.NewCond(InstanceFeatureInstanceIDCol, e.Aggregate().InstanceID),
	}), nil
}