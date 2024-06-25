package metrics

import (
	"context"
	"fmt"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottlmetric"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"sort"
)

type extractP95FromHistogramArguments struct {
}

func newExtractP95FromHistogramFactory() ottl.Factory[ottlmetric.TransformContext] {
	return ottl.NewFactory("extract_p95_from_histogram", &extractP95FromHistogramArguments{}, createExtractP95FromHistogramFunction)
}

func createExtractP95FromHistogramFunction(_ ottl.FunctionContext, oArgs ottl.Arguments) (ottl.ExprFunc[ottlmetric.TransformContext], error) {
	//args, ok := oArgs.(*extractP95FromHistogramArguments)

	//if !ok {
	//	return nil, fmt.Errorf("extractP95FromHistogramFactory args must be of type *extractP95FromHistogramArguments")
	//}

	return extractP95FromHistogram()
}

func extractP95FromHistogram() (ottl.ExprFunc[ottlmetric.TransformContext], error) {
	return func(_ context.Context, tCtx ottlmetric.TransformContext) (any, error) {
		metric := tCtx.GetMetric()
		newMetric := tCtx.GetMetrics().AppendEmpty()
		invalidMetricTypeError := fmt.Errorf("extract_p95_from_histogram requires an input metric of type Histogram or ExponentialHistogram got %s", metric.Type())

		newMetric.SetDescription(metric.Description())
		newMetric.SetName(metric.Name() + ".95percentile")
		newMetric.SetUnit(metric.Unit())
		newMetric.SetEmptyGauge()
		pmetric.MetricTypeExponentialHistogram
		switch metric.Type() {
		case pmetric.MetricTypeHistogram:
			dataPoints := metric.Histogram().DataPoints()
			var dataPointSlc []float64
			for i := 0; i < dataPoints.Len(); i++ {
				dataPointSlc = append(dataPointSlc, dataPoints.At(i).Max())
			}
			sort.Float64s(dataPointSlc)
			f := dataPointSlc[int(float64(dataPoints.Len())*0.95)]
			newMetric.Gauge().DataPoints().AppendEmpty().SetDoubleValue(f)
		case pmetric.MetricTypeExponentialHistogram:
			dataPoints := metric.ExponentialHistogram().DataPoints()
			var dataPointSlc []float64
			for i := 0; i < dataPoints.Len(); i++ {
				dataPointSlc = append(dataPointSlc, dataPoints.At(i).Max())
			}
			sort.Float64s(dataPointSlc)
			f := dataPointSlc[int(float64(dataPoints.Len())*0.95)]
			newMetric.Gauge().DataPoints().AppendEmpty().SetDoubleValue(f)
		default:
			return nil, invalidMetricTypeError
		}
		return nil, nil
	}, nil
}
