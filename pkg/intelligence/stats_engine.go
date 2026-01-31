package intelligence

import (
	"math"
	"sort"
)

// StatisticalEngine provides data-driven analysis utilities
type StatisticalEngine struct{}

// NewStatisticalEngine creates a new statistical engine
func NewStatisticalEngine() *StatisticalEngine {
	return &StatisticalEngine{}
}

// Distribution represents a statistical distribution of values
type Distribution struct {
	Values      []float64
	Mean        float64
	StdDev      float64
	Percentiles map[int]float64 // p25, p50, p75, p90
	SampleSize  int
	Min         float64
	Max         float64
}

// CalculateDistribution computes statistics from a slice of values
func (e *StatisticalEngine) CalculateDistribution(values []float64) *Distribution {
	if len(values) == 0 {
		return &Distribution{
			SampleSize:  0,
			Percentiles: make(map[int]float64),
		}
	}

	// Sort for percentile calculation
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)

	// Calculate mean
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(len(values))

	// Calculate standard deviation
	sumSquares := 0.0
	for _, v := range values {
		sumSquares += (v - mean) * (v - mean)
	}
	stdDev := math.Sqrt(sumSquares / float64(len(values)))

	// Calculate percentiles
	percentiles := map[int]float64{
		25: e.percentile(sorted, 25),
		50: e.percentile(sorted, 50),
		75: e.percentile(sorted, 75),
		90: e.percentile(sorted, 90),
	}

	return &Distribution{
		Values:      values,
		Mean:        mean,
		StdDev:      stdDev,
		Percentiles: percentiles,
		SampleSize:  len(values),
		Min:         sorted[0],
		Max:         sorted[len(sorted)-1],
	}
}


// percentile calculates the p-th percentile of sorted values
func (e *StatisticalEngine) percentile(sorted []float64, p int) float64 {
	if len(sorted) == 0 {
		return 0
	}
	if len(sorted) == 1 {
		return sorted[0]
	}

	// Linear interpolation method
	rank := float64(p) / 100.0 * float64(len(sorted)-1)
	lower := int(math.Floor(rank))
	upper := int(math.Ceil(rank))

	if lower == upper {
		return sorted[lower]
	}

	// Interpolate between lower and upper
	weight := rank - float64(lower)
	return sorted[lower]*(1-weight) + sorted[upper]*weight
}

// ClassifyValue determines if a value is low/average/high relative to distribution
func (e *StatisticalEngine) ClassifyValue(value float64, dist *Distribution) string {
	if dist.SampleSize < 3 {
		return "insufficient_data"
	}

	p25 := dist.Percentiles[25]
	p75 := dist.Percentiles[75]

	if value < p25 {
		return "low"
	} else if value > p75 {
		return "high"
	}
	return "average"
}

// ClassifyTiming determines if a timing is early/average/late relative to distribution
func (e *StatisticalEngine) ClassifyTiming(value float64, dist *Distribution) string {
	if dist.SampleSize < 3 {
		return "insufficient_data"
	}

	p25 := dist.Percentiles[25]
	p75 := dist.Percentiles[75]

	if value < p25 {
		return "early"
	} else if value > p75 {
		return "late"
	}
	return "average"
}

// CalculateConfidence returns confidence score based on sample size and variance
// Returns a value between 0.0 and 1.0
func (e *StatisticalEngine) CalculateConfidence(dist *Distribution) float64 {
	if dist.SampleSize == 0 {
		return 0.0
	}

	// Base confidence from sample size (logarithmic scaling)
	// n=1 -> ~0.15, n=10 -> ~0.5, n=100 -> ~1.0
	sampleConfidence := math.Min(1.0, math.Log10(float64(dist.SampleSize)+1)/2.0)

	// Adjust for variance (lower variance = higher confidence)
	if dist.Mean != 0 {
		cv := dist.StdDev / math.Abs(dist.Mean) // Coefficient of variation
		varianceAdjustment := math.Max(0.5, 1.0-cv)
		return math.Min(1.0, sampleConfidence*varianceAdjustment)
	}

	return sampleConfidence
}

// GetPercentileRank returns where a value falls in the distribution (0-100)
func (e *StatisticalEngine) GetPercentileRank(value float64, dist *Distribution) int {
	if dist.SampleSize == 0 {
		return 50 // Default to median
	}

	// Count values below
	below := 0
	for _, v := range dist.Values {
		if v < value {
			below++
		}
	}

	return int(float64(below) / float64(dist.SampleSize) * 100)
}

// IsOutlier determines if a value is an outlier using IQR method
func (e *StatisticalEngine) IsOutlier(value float64, dist *Distribution) bool {
	if dist.SampleSize < 4 {
		return false
	}

	p25 := dist.Percentiles[25]
	p75 := dist.Percentiles[75]
	iqr := p75 - p25

	lowerBound := p25 - 1.5*iqr
	upperBound := p75 + 1.5*iqr

	return value < lowerBound || value > upperBound
}

// RemoveOutliers returns a new distribution with outliers removed
func (e *StatisticalEngine) RemoveOutliers(dist *Distribution) *Distribution {
	if dist.SampleSize < 4 {
		return dist
	}

	p25 := dist.Percentiles[25]
	p75 := dist.Percentiles[75]
	iqr := p75 - p25

	lowerBound := p25 - 1.5*iqr
	upperBound := p75 + 1.5*iqr

	filtered := make([]float64, 0, len(dist.Values))
	for _, v := range dist.Values {
		if v >= lowerBound && v <= upperBound {
			filtered = append(filtered, v)
		}
	}

	return e.CalculateDistribution(filtered)
}

// CompareDistributions returns a comparison between two distributions
type DistributionComparison struct {
	MeanDiff       float64 // Difference in means
	MedianDiff     float64 // Difference in medians
	OverlapPercent float64 // Percentage of overlap between distributions
	Significance   string  // "significant", "moderate", "minimal"
}

// CompareDistributions compares two distributions
func (e *StatisticalEngine) CompareDistributions(a, b *Distribution) *DistributionComparison {
	if a.SampleSize == 0 || b.SampleSize == 0 {
		return &DistributionComparison{Significance: "insufficient_data"}
	}

	meanDiff := a.Mean - b.Mean
	medianDiff := a.Percentiles[50] - b.Percentiles[50]

	// Calculate overlap using IQR ranges
	aLow, aHigh := a.Percentiles[25], a.Percentiles[75]
	bLow, bHigh := b.Percentiles[25], b.Percentiles[75]

	overlapLow := math.Max(aLow, bLow)
	overlapHigh := math.Min(aHigh, bHigh)

	var overlapPercent float64
	if overlapHigh > overlapLow {
		totalRange := math.Max(aHigh, bHigh) - math.Min(aLow, bLow)
		if totalRange > 0 {
			overlapPercent = (overlapHigh - overlapLow) / totalRange * 100
		}
	}

	// Determine significance based on effect size (Cohen's d approximation)
	pooledStdDev := math.Sqrt((a.StdDev*a.StdDev + b.StdDev*b.StdDev) / 2)
	var significance string
	if pooledStdDev > 0 {
		effectSize := math.Abs(meanDiff) / pooledStdDev
		if effectSize > 0.8 {
			significance = "significant"
		} else if effectSize > 0.5 {
			significance = "moderate"
		} else {
			significance = "minimal"
		}
	} else {
		significance = "minimal"
	}

	return &DistributionComparison{
		MeanDiff:       meanDiff,
		MedianDiff:     medianDiff,
		OverlapPercent: overlapPercent,
		Significance:   significance,
	}
}
