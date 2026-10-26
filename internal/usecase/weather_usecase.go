package usecase

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"go_26_10/domain"
)

type weatherUsecase struct {
	weatherRepo  domain.WeatherRepository
	locationRepo domain.LocationRepository
	userRepo     domain.UserRepository
}

// NewWeatherUsecase creates a new instance of WeatherUsecase
func NewWeatherUsecase(
	weatherRepo domain.WeatherRepository,
	locationRepo domain.LocationRepository,
	userRepo domain.UserRepository,
) domain.WeatherUsecase {
	return &weatherUsecase{
		weatherRepo:  weatherRepo,
		locationRepo: locationRepo,
		userRepo:     userRepo,
	}
}

func (u *weatherUsecase) resolveLocation(ctx context.Context, userID *uuid.UUID, locationID *uuid.UUID) *domain.Location {
	defaultState := "CA"
	defaultImage := "https://images.unsplash.com/photo-1580655653885-65763b2597d0?auto=format&fit=crop&w=800&q=80"
	targetLocation := &domain.Location{
		ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		City:      "Los Angeles",
		State:     &defaultState,
		Country:   "USA",
		Latitude:  34.052235,
		Longitude: -118.243683,
		ImageURL:  &defaultImage,
	}

	if locationID != nil {
		if loc, err := u.locationRepo.FindByID(ctx, *locationID); err == nil && loc != nil {
			return loc
		}
		log.Printf("[INFO] Location ID %s not found in DB, using default %s\n", locationID.String(), targetLocation.City)
	} else if userID != nil {
		if saved, err := u.locationRepo.GetDefaultUserLocation(ctx, *userID); err == nil && saved != nil && saved.Location != nil {
			return saved.Location
		}
	}

	return targetLocation
}

func (u *weatherUsecase) resolveTemperatureUnit(ctx context.Context, userID *uuid.UUID, unitParam string) string {
	if strings.EqualFold(unitParam, "fahrenheit") {
		return "Fahrenheit"
	}
	if strings.EqualFold(unitParam, "celsius") {
		return "Celsius"
	}

	if userID != nil {
		if user, err := u.userRepo.FindByID(ctx, *userID); err == nil && user != nil {
			if strings.EqualFold(user.TemperatureUnit, "fahrenheit") {
				return "Fahrenheit"
			}
		}
	}

	return "Celsius"
}

func toFahrenheit(celsius float64) float64 {
	return math.Round(((celsius*9.0/5.0)+32.0)*10) / 10
}

func toLocationInfo(loc *domain.Location) domain.WeatherLocationInfo {
	return domain.WeatherLocationInfo{
		ID:       loc.ID,
		City:     loc.City,
		State:    loc.State,
		Country:  loc.Country,
		ImageURL: loc.ImageURL,
	}
}

func formatForecastTemperatures(forecasts []domain.DayForecast, isFahrenheit bool) []domain.DayForecast {
	unitSymbol := "°C"
	if isFahrenheit {
		unitSymbol = "°F"
	}

	formatted := make([]domain.DayForecast, len(forecasts))
	for i, item := range forecasts {
		high := item.TempHigh
		low := item.TempLow
		if isFahrenheit {
			high = toFahrenheit(high)
			low = toFahrenheit(low)
		}

		highDisp := fmt.Sprintf("↑ %.0f%s", high, unitSymbol)
		lowDisp := fmt.Sprintf("↓ %.0f%s", low, unitSymbol)
		summary := fmt.Sprintf("%s, estimasi suhu %s dan %s", item.Condition, highDisp, lowDisp)

		itemCopy := item
		itemCopy.TempHigh = high
		itemCopy.TempLow = low
		itemCopy.TempHighDisplay = highDisp
		itemCopy.TempLowDisplay = lowDisp
		itemCopy.Summary = summary
		formatted[i] = itemCopy
	}
	return formatted
}

func (u *weatherUsecase) buildHourlyChart(loc domain.WeatherLocationInfo, hourly []domain.HourlyTemperature, unit string) *domain.HourlyChartData {
	if len(hourly) == 0 {
		return &domain.HourlyChartData{
			Location:   loc,
			Unit:       unit,
			Labels:     []string{},
			DataPoints: []float64{},
			Series:     []domain.HourlyTemperature{},
		}
	}

	isFahrenheit := strings.EqualFold(unit, "fahrenheit")
	chartSeries := make([]domain.HourlyTemperature, len(hourly))
	labels := make([]string, len(hourly))
	dataPoints := make([]float64, len(hourly))

	minTemp := math.MaxFloat64
	maxTemp := -math.MaxFloat64
	totalTemp := 0.0
	peakTime := ""
	lowestTime := ""

	for i, item := range hourly {
		tempVal := item.Temp
		if isFahrenheit {
			if item.TempF != 0 {
				tempVal = item.TempF
			} else {
				tempVal = toFahrenheit(item.Temp)
			}
		} else {
			if item.TempC != 0 {
				tempVal = item.TempC
			}
		}

		itemCopy := item
		itemCopy.Temp = tempVal
		chartSeries[i] = itemCopy

		labels[i] = item.Time
		dataPoints[i] = tempVal
		totalTemp += tempVal

		if tempVal > maxTemp {
			maxTemp = tempVal
			peakTime = item.Time
		}
		if tempVal < minTemp {
			minTemp = tempVal
			lowestTime = item.Time
		}
	}

	avgTemp := math.Round((totalTemp/float64(len(hourly)))*10) / 10

	unitSymbol := "°C"
	unitName := "Celsius"
	if isFahrenheit {
		unitSymbol = "°F"
		unitName = "Fahrenheit"
	}

	trendDesc := fmt.Sprintf("Suhu harian berfluktuasi dari %.1f%s (%s) hingga titik puncak %.1f%s (%s) dengan rata-rata %.1f%s.",
		minTemp, unitSymbol, lowestTime, maxTemp, unitSymbol, peakTime, avgTemp, unitSymbol)

	return &domain.HourlyChartData{
		Location:   loc,
		Unit:       unitName,
		MinTemp:    minTemp,
		MaxTemp:    maxTemp,
		AvgTemp:    avgTemp,
		PeakTime:   peakTime,
		LowestTime: lowestTime,
		Trend:      trendDesc,
		Labels:     labels,
		DataPoints: dataPoints,
		Series:     chartSeries,
	}
}

func (u *weatherUsecase) buildRainAnalysis(loc domain.WeatherLocationInfo, rainList []domain.ChanceOfRain, dateStr string) *domain.RainAnalysis {
	if len(rainList) == 0 {
		return &domain.RainAnalysis{
			Location:       loc,
			Date:           dateStr,
			RiskLevel:      "Low",
			Advisory:       "Tidak ada data prakiraan hujan.",
			TimeRanges:     []domain.RainTimeRange{},
			HourlyForecast: []domain.ChanceOfRain{},
		}
	}

	totalPercentage := 0
	peakPercentage := -1
	peakTime := ""
	peakIntensity := "Sunny"

	type periodBucket struct {
		periodName string
		timeSlot   string
		items      []domain.ChanceOfRain
	}

	buckets := []*periodBucket{
		{periodName: "Pagi", timeSlot: "06:00 - 12:00", items: nil},
		{periodName: "Siang", timeSlot: "12:00 - 18:00", items: nil},
		{periodName: "Sore / Malam", timeSlot: "18:00 - 24:00", items: nil},
		{periodName: "Dini Hari", timeSlot: "00:00 - 06:00", items: nil},
	}

	for _, item := range rainList {
		totalPercentage += item.Percentage
		if item.Percentage > peakPercentage {
			peakPercentage = item.Percentage
			peakTime = item.Time
			peakIntensity = item.Intensity
		}

		h := item.Hour
		if h == 0 && item.Time != "12 am" && item.Time != "00:00" {
			if t, err := time.Parse("03 pm", strings.ToLower(item.Time)); err == nil {
				h = t.Hour()
			}
		}

		if h >= 6 && h < 12 {
			buckets[0].items = append(buckets[0].items, item)
		} else if h >= 12 && h < 18 {
			buckets[1].items = append(buckets[1].items, item)
		} else if h >= 18 && h < 24 {
			buckets[2].items = append(buckets[2].items, item)
		} else {
			buckets[3].items = append(buckets[3].items, item)
		}
	}

	dailyAverage := totalPercentage / len(rainList)

	var riskLevel string
	switch {
	case peakPercentage >= 75:
		riskLevel = "Very High"
	case peakPercentage >= 50:
		riskLevel = "High"
	case peakPercentage >= 25:
		riskLevel = "Moderate"
	default:
		riskLevel = "Low"
	}

	var advisory string
	switch riskLevel {
	case "Very High":
		advisory = fmt.Sprintf("Waspada potensi hujan lebat dengan peluang mencapai %d%% pada pukul %s (%s). Sangat disarankan membawa payung atau jas hujan saat beraktivitas.", peakPercentage, peakTime, peakIntensity)
	case "High":
		advisory = fmt.Sprintf("Peluang hujan cukup tinggi hingga %d%% pada sekitar pukul %s (%s). Siapkan perlengkapan hujan sebelum bepergian.", peakPercentage, peakTime, peakIntensity)
	case "Moderate":
		advisory = fmt.Sprintf("Potensi hujan ringan pada sekitar pukul %s (%d%%). Cuaca relatif aman untuk kegiatan luar ruangan.", peakTime, peakPercentage)
	default:
		advisory = fmt.Sprintf("Peluang hujan rendah sepanjang hari (maksimal %d%%). Cuaca diperkirakan cerah berawan.", peakPercentage)
	}

	var timeRanges []domain.RainTimeRange
	for _, b := range buckets {
		if len(b.items) == 0 {
			continue
		}
		bucketTotal := 0
		bucketMax := -1
		dominantIntensity := "Sunny"

		for _, it := range b.items {
			bucketTotal += it.Percentage
			if it.Percentage > bucketMax {
				bucketMax = it.Percentage
				dominantIntensity = it.Intensity
			}
		}
		bucketAvg := bucketTotal / len(b.items)

		var desc string
		if bucketMax > 70 {
			desc = fmt.Sprintf("Potensi hujan lebat signifikan (hingga %d%%) dengan intensitas %s.", bucketMax, dominantIntensity)
		} else if bucketMax > 30 {
			desc = fmt.Sprintf("Berpotensi hujan ringan hingga sedang (hingga %d%%).", bucketMax)
		} else {
			desc = fmt.Sprintf("Kondisi cerah berawan dengan peluang hujan minim (maksimal %d%%).", bucketMax)
		}

		timeRanges = append(timeRanges, domain.RainTimeRange{
			Period:        b.periodName,
			TimeSlot:      b.timeSlot,
			AverageChance: bucketAvg,
			MaxChance:     bucketMax,
			Intensity:     dominantIntensity,
			Description:   desc,
		})
	}

	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	return &domain.RainAnalysis{
		Location:       loc,
		Date:           dateStr,
		DailyAverage:   dailyAverage,
		PeakTime:       peakTime,
		PeakPercentage: peakPercentage,
		PeakIntensity:  peakIntensity,
		RiskLevel:      riskLevel,
		Advisory:       advisory,
		TimeRanges:     timeRanges,
		HourlyForecast: rainList,
	}
}

func (u *weatherUsecase) getDashboardForLocation(ctx context.Context, targetLocation *domain.Location, unit string) (*domain.WeatherDashboardData, error) {
	// 1. Ambil kondisi cuaca saat ini (Current Weather)
	currentWeather, err := u.weatherRepo.GetRealtimeWeather(ctx, targetLocation.Latitude, targetLocation.Longitude)
	if err != nil {
		return nil, err
	}

	// 2. Ambil data Today's Highlights (Precipitation, Humidity, Wind, Sunrise, Sunset)
	highlights, err := u.weatherRepo.GetHighlights(ctx, targetLocation.Latitude, targetLocation.Longitude)
	if err != nil {
		return nil, err
	}

	// 3. Ambil Hourly Temperature & Chance of Rain
	hourly, chanceOfRain, err := u.weatherRepo.GetHourlyForecast(ctx, targetLocation.Latitude, targetLocation.Longitude)
	if err != nil {
		log.Printf("[WARN] Failed to fetch hourly forecast: %v\n", err)
	}

	// 4. Ambil 3 Days Forecast
	threeDays, err := u.weatherRepo.GetMultiDayForecast(ctx, targetLocation.Latitude, targetLocation.Longitude, 3)
	if err != nil {
		log.Printf("[WARN] Failed to fetch 3-day forecast: %v\n", err)
	}

	locInfo := toLocationInfo(targetLocation)
	isFahrenheit := strings.EqualFold(unit, "fahrenheit")

	// Sesuaikan unit suhu jika Fahrenheit
	if isFahrenheit {
		currentWeather.Temperature = toFahrenheit(currentWeather.Temperature)
		for i := range hourly {
			if hourly[i].TempF != 0 {
				hourly[i].Temp = hourly[i].TempF
			} else {
				hourly[i].Temp = toFahrenheit(hourly[i].Temp)
			}
		}
	}

	formattedForecast := formatForecastTemperatures(threeDays, isFahrenheit)
	chartData := u.buildHourlyChart(locInfo, hourly, unit)
	rainData := u.buildRainAnalysis(locInfo, chanceOfRain, time.Now().Format("2006-01-02"))

	return &domain.WeatherDashboardData{
		Location:          locInfo,
		CurrentWeather:    *currentWeather,
		Highlights:        *highlights,
		HourlyTemperature: hourly,
		ChanceOfRain:      chanceOfRain,
		ThreeDaysForecast: formattedForecast,
		HourlyChart:       chartData,
		RainAnalysis:      rainData,
	}, nil
}

func (u *weatherUsecase) resolveLocationByQueryOrID(
	ctx context.Context,
	query string,
	locID *uuid.UUID,
	fallback *domain.Location,
) *domain.Location {
	if locID != nil {
		if loc, err := u.locationRepo.FindByID(ctx, *locID); err == nil && loc != nil {
			return loc
		}
	}

	trimmed := strings.TrimSpace(query)
	if trimmed != "" {
		if parsedUUID, err := uuid.Parse(trimmed); err == nil {
			if loc, err := u.locationRepo.FindByID(ctx, parsedUUID); err == nil && loc != nil {
				return loc
			}
		}

		if locs, err := u.locationRepo.Search(ctx, trimmed); err == nil && len(locs) > 0 {
			for _, l := range locs {
				if strings.EqualFold(l.City, trimmed) {
					locCopy := l
					return &locCopy
				}
			}
			locCopy := locs[0]
			return &locCopy
		}
	}

	return fallback
}

func (u *weatherUsecase) GetDashboardData(ctx context.Context, userID *uuid.UUID, locationID *uuid.UUID) (*domain.WeatherDashboardData, error) {
	targetLocation := u.resolveLocation(ctx, userID, locationID)
	unit := u.resolveTemperatureUnit(ctx, userID, "")
	return u.getDashboardForLocation(ctx, targetLocation, unit)
}

func (u *weatherUsecase) CompareCities(
	ctx context.Context,
	userID *uuid.UUID,
	query1, query2 string,
	locID1, locID2 *uuid.UUID,
	unitParam string,
) (*domain.WeatherComparison, error) {
	unit := u.resolveTemperatureUnit(ctx, userID, unitParam)

	// Fallback location 1: Los Angeles
	defaultStateCA := "CA"
	defaultImageLA := "https://images.unsplash.com/photo-1580655653885-65763b2597d0?auto=format&fit=crop&w=800&q=80"
	defaultLoc1 := &domain.Location{
		ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		City:      "Los Angeles",
		State:     &defaultStateCA,
		Country:   "USA",
		Latitude:  34.052235,
		Longitude: -118.243683,
		ImageURL:  &defaultImageLA,
	}

	// Fallback location 2: Jakarta
	defaultStateJKT := "DKI Jakarta"
	defaultImageJKT := "https://images.unsplash.com/photo-1555899434-94d1368aa7af?auto=format&fit=crop&w=800&q=80"
	defaultLoc2 := &domain.Location{
		ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		City:      "Jakarta",
		State:     &defaultStateJKT,
		Country:   "Indonesia",
		Latitude:  -6.208763,
		Longitude: 106.845599,
		ImageURL:  &defaultImageJKT,
	}

	loc1 := u.resolveLocationByQueryOrID(ctx, query1, locID1, defaultLoc1)
	loc2 := u.resolveLocationByQueryOrID(ctx, query2, locID2, defaultLoc2)

	// If both locations resolved to the same city and query2 was empty, pick Tokyo as loc2
	if loc1.ID == loc2.ID && strings.TrimSpace(query2) == "" && locID2 == nil {
		if tokyo, err := u.locationRepo.FindByID(ctx, uuid.MustParse("00000000-0000-0000-0000-000000000003")); err == nil && tokyo != nil {
			loc2 = tokyo
		}
	}

	d1, err := u.getDashboardForLocation(ctx, loc1, unit)
	if err != nil {
		return nil, fmt.Errorf("gagal memuat data cuaca untuk %s: %w", loc1.City, err)
	}

	d2, err := u.getDashboardForLocation(ctx, loc2, unit)
	if err != nil {
		return nil, fmt.Errorf("gagal memuat data cuaca untuk %s: %w", loc2.City, err)
	}

	unitSymbol := "°C"
	if strings.EqualFold(unit, "fahrenheit") {
		unitSymbol = "°F"
	}

	// 1. Temperature Comparison
	temp1 := d1.CurrentWeather.Temperature
	temp2 := d2.CurrentWeather.Temperature
	tempDiff := math.Round((temp1-temp2)*10) / 10

	var warmerCity, colderCity, tempDiffFormatted string
	if tempDiff > 0 {
		warmerCity = d1.Location.City
		colderCity = d2.Location.City
		tempDiffFormatted = fmt.Sprintf("%s lebih hangat %.1f%s dibanding %s", warmerCity, tempDiff, unitSymbol, colderCity)
	} else if tempDiff < 0 {
		warmerCity = d2.Location.City
		colderCity = d1.Location.City
		tempDiffFormatted = fmt.Sprintf("%s lebih hangat %.1f%s dibanding %s", warmerCity, math.Abs(tempDiff), unitSymbol, colderCity)
	} else {
		warmerCity = d1.Location.City
		colderCity = d2.Location.City
		tempDiffFormatted = fmt.Sprintf("%s dan %s memiliki suhu yang sama (%.1f%s)", d1.Location.City, d2.Location.City, temp1, unitSymbol)
	}

	// 2. Humidity Comparison
	parsePercent := func(s string) int {
		s = strings.TrimSuffix(strings.TrimSpace(s), "%")
		val, _ := strconv.Atoi(s)
		return val
	}
	h1 := parsePercent(d1.Highlights.Humidity)
	h2 := parsePercent(d2.Highlights.Humidity)
	var humidityComp string
	if h1 > h2 {
		humidityComp = fmt.Sprintf("%s lebih lembap (+%d%%) dibanding %s", d1.Location.City, h1-h2, d2.Location.City)
	} else if h2 > h1 {
		humidityComp = fmt.Sprintf("%s lebih lembap (+%d%%) dibanding %s", d2.Location.City, h2-h1, d1.Location.City)
	} else {
		humidityComp = fmt.Sprintf("Kedua kota memiliki kelembapan serupa (%d%%)", h1)
	}

	// 3. Rainier City Comparison
	r1 := parsePercent(d1.Highlights.Precipitation)
	r2 := parsePercent(d2.Highlights.Precipitation)
	var rainierCity string
	if r1 > r2 {
		rainierCity = d1.Location.City
	} else if r2 > r1 {
		rainierCity = d2.Location.City
	} else {
		rainierCity = "Kedua kota seimbang"
	}

	// 4. Windier City Comparison
	parseSpeed := func(s string) float64 {
		s = strings.TrimSuffix(s, " km/h")
		s = strings.TrimSpace(s)
		val, _ := strconv.ParseFloat(s, 64)
		return val
	}
	w1 := parseSpeed(d1.Highlights.WindSpeed)
	w2 := parseSpeed(d2.Highlights.WindSpeed)
	var windierCity string
	if w1 > w2 {
		windierCity = fmt.Sprintf("%s (%.0f km/h vs %.0f km/h)", d1.Location.City, w1, w2)
	} else if w2 > w1 {
		windierCity = fmt.Sprintf("%s (%.0f km/h vs %.0f km/h)", d2.Location.City, w2, w1)
	} else {
		windierCity = fmt.Sprintf("Kecepatan angin serupa (%.0f km/h)", w1)
	}

	// 5. Condition Summary
	conditionSummary := fmt.Sprintf("%s (%s, %.1f%s) vs %s (%s, %.1f%s)",
		d1.Location.City, d1.CurrentWeather.Condition, temp1, unitSymbol,
		d2.Location.City, d2.CurrentWeather.Condition, temp2, unitSymbol)

	// 6. Hourly Comparison Series
	var hourlyComp []domain.HourlyComparisonItem
	maxItems := len(d1.HourlyTemperature)
	if len(d2.HourlyTemperature) < maxItems {
		maxItems = len(d2.HourlyTemperature)
	}
	for i := 0; i < maxItems; i++ {
		item1 := d1.HourlyTemperature[i]
		item2 := d2.HourlyTemperature[i]
		hourlyComp = append(hourlyComp, domain.HourlyComparisonItem{
			Time:  item1.Time,
			Temp1: item1.Temp,
			Temp2: item2.Temp,
			Diff:  math.Round((item1.Temp-item2.Temp)*10) / 10,
		})
	}

	// 7. Practical Travel / Activity Recommendation
	var rec strings.Builder
	if warmerCity != "" && math.Abs(tempDiff) >= 5 {
		rec.WriteString(fmt.Sprintf("%s terasa jauh lebih sejuk/dingin daripada %s (selisih %.1f%s). ", colderCity, warmerCity, math.Abs(tempDiff), unitSymbol))
	}
	if r1 > 40 || r2 > 40 {
		if r1 > r2 {
			rec.WriteString(fmt.Sprintf("Disarankan menyiapkan payung di %s karena peluang hujan tinggi (%d%%). ", d1.Location.City, r1))
		} else if r2 > r1 {
			rec.WriteString(fmt.Sprintf("Disarankan menyiapkan payung di %s karena peluang hujan tinggi (%d%%). ", d2.Location.City, r2))
		}
	}
	if strings.Contains(strings.ToLower(d1.CurrentWeather.Condition), "rain") || strings.Contains(strings.ToLower(d2.CurrentWeather.Condition), "rain") {
		rec.WriteString("Perhatikan kondisi jalanan licin untuk perjalanan antar kota.")
	} else {
		rec.WriteString("Kondisi cuaca secara umum kondusif untuk aktivitas luar ruangan.")
	}

	summary := domain.ComparisonSummary{
		TempDifference:     tempDiff,
		TempDiffFormatted:  tempDiffFormatted,
		WarmerCity:         warmerCity,
		ColderCity:         colderCity,
		HumidityComparison: humidityComp,
		RainierCity:        rainierCity,
		WindierCity:        windierCity,
		ConditionSummary:   conditionSummary,
		Recommendation:     rec.String(),
		HourlyComparison:   hourlyComp,
	}

	return &domain.WeatherComparison{
		City1:      *d1,
		City2:      *d2,
		Comparison: summary,
	}, nil
}

func (u *weatherUsecase) GetHourlyChart(ctx context.Context, userID *uuid.UUID, locationID *uuid.UUID, unitParam string) (*domain.HourlyChartData, error) {
	targetLocation := u.resolveLocation(ctx, userID, locationID)
	unit := u.resolveTemperatureUnit(ctx, userID, unitParam)

	hourly, _, err := u.weatherRepo.GetHourlyForecast(ctx, targetLocation.Latitude, targetLocation.Longitude)
	if err != nil {
		return nil, err
	}

	locInfo := toLocationInfo(targetLocation)
	chartData := u.buildHourlyChart(locInfo, hourly, unit)
	return chartData, nil
}

func (u *weatherUsecase) GetChanceOfRainAnalysis(ctx context.Context, userID *uuid.UUID, locationID *uuid.UUID) (*domain.RainAnalysis, error) {
	targetLocation := u.resolveLocation(ctx, userID, locationID)

	_, chanceOfRain, err := u.weatherRepo.GetHourlyForecast(ctx, targetLocation.Latitude, targetLocation.Longitude)
	if err != nil {
		return nil, err
	}

	locInfo := toLocationInfo(targetLocation)
	rainAnalysis := u.buildRainAnalysis(locInfo, chanceOfRain, time.Now().Format("2006-01-02"))
	return rainAnalysis, nil
}

func (u *weatherUsecase) GetForecast(ctx context.Context, userID *uuid.UUID, locationID *uuid.UUID, days int, unitParam string) (*domain.ForecastData, error) {
	if days <= 0 {
		days = 3
	}
	targetLocation := u.resolveLocation(ctx, userID, locationID)
	unit := u.resolveTemperatureUnit(ctx, userID, unitParam)
	isFahrenheit := strings.EqualFold(unit, "fahrenheit")

	forecasts, err := u.weatherRepo.GetMultiDayForecast(ctx, targetLocation.Latitude, targetLocation.Longitude, days)
	if err != nil {
		return nil, err
	}

	formattedForecast := formatForecastTemperatures(forecasts, isFahrenheit)
	locInfo := toLocationInfo(targetLocation)

	return &domain.ForecastData{
		Location: locInfo,
		Unit:     unit,
		Days:     days,
		Forecast: formattedForecast,
	}, nil
}

func (u *weatherUsecase) GetCalendar(ctx context.Context, userID *uuid.UUID, city string, year, month int, unitParam string) (*domain.WeatherCalendarData, error) {
	if city == "" {
		targetLocation := u.resolveLocation(ctx, userID, nil)
		city = targetLocation.City
	}

	unit := u.resolveTemperatureUnit(ctx, userID, unitParam)
	isFahrenheit := strings.EqualFold(unit, "fahrenheit")

	calendarData, err := u.weatherRepo.GetCalendarForecast(ctx, city, year, month)
	if err != nil {
		return nil, err
	}

	if isFahrenheit {
		for i := range calendarData.Days {
			calendarData.Days[i].TempHigh = math.Round(((calendarData.Days[i].TempHigh*9/5)+32)*10) / 10
			calendarData.Days[i].TempLow = math.Round(((calendarData.Days[i].TempLow*9/5)+32)*10) / 10
		}
	}

	return calendarData, nil
}
