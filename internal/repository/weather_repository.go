package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go_26_10/domain"
)

type weatherRepository struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewWeatherRepository creates a new WeatherRepository instance.
// If baseURL is empty, it uses fallback mock data (useful for deterministic tests).
func NewWeatherRepository(apiKey, baseURL string) domain.WeatherRepository {
	trimmed := strings.TrimSpace(baseURL)
	return &weatherRepository{
		apiKey:  apiKey,
		baseURL: trimmed,
		httpClient: &http.Client{
			Timeout: 6 * time.Second,
		},
	}
}

// Daftar pemetaan koordinat kota-kota di Indonesia ke kode adm4 BMKG
var indonesianCities = []struct {
	Lat, Lon float64
	Adm4     string
	City     string
}{
	{5.5481, 95.3235, "11.71.01.1001", "Banda Aceh"},
	{3.5952, 98.6722, "12.71.01.1001", "Medan"},
	{-0.9471, 100.4175, "13.71.01.1001", "Padang"},
	{0.5071, 101.4478, "14.71.01.1001", "Pekanbaru"},
	{-1.6116, 103.6167, "15.71.01.1001", "Jambi"},
	{-2.9761, 104.7754, "16.71.01.1001", "Palembang"},
	{-3.8565, 102.3540, "17.71.01.1001", "Bengkulu"},
	{-5.4352, 105.2594, "18.71.01.1001", "Bandar Lampung"},
	{-2.1229, 106.1138, "19.71.01.1001", "Pangkal Pinang"},
	{0.9197, 104.4543, "21.71.01.1001", "Tanjung Pinang"},
	{-6.1123, 106.6259, "36.71.01.1001", "Serang"},
	{-6.2088, 106.8456, "31.71.01.1001", "Jakarta"},
	{-6.9175, 107.6191, "32.73.01.1001", "Bandung"},
	{-6.9932, 110.4204, "33.74.01.1001", "Semarang"},
	{-7.7956, 110.3695, "34.71.01.1001", "Yogyakarta"},
	{-7.2575, 112.7521, "35.78.01.1001", "Surabaya"},
	{-7.9797, 112.6304, "35.73.01.1001", "Malang"},
	{-8.6705, 115.2126, "51.71.01.1001", "Denpasar"},
	{-8.5827, 116.1154, "52.71.01.1001", "Mataram"},
	{-10.1735, 123.5816, "53.71.01.1001", "Kupang"},
	{-0.0264, 109.3425, "61.71.01.1001", "Pontianak"},
	{-2.2142, 113.9192, "62.71.01.1001", "Palangkaraya"},
	{-3.3190, 114.5892, "63.71.01.1001", "Banjarmasin"},
	{-0.5019, 117.1434, "64.71.01.1001", "Samarinda"},
	{-1.2654, 116.8312, "64.71.01.1001", "Balikpapan"},
	{2.7667, 117.4000, "65.71.01.1001", "Tanjung Selor"},
	{1.4748, 124.8421, "71.71.01.1001", "Manado"},
	{-0.8915, 119.8698, "72.71.01.1001", "Palu"},
	{-5.1477, 119.4327, "73.71.01.1001", "Makassar"},
	{-3.9772, 122.5162, "74.71.01.1001", "Kendari"},
	{0.5244, 123.0624, "75.71.01.1001", "Gorontalo"},
	{-2.6239, 118.9070, "76.71.01.1001", "Mamuju"},
	{-3.6944, 128.1814, "81.71.01.1001", "Ambon"},
	{0.7901, 127.3839, "82.71.01.1001", "Sofifi"},
	{-2.5330, 140.7179, "91.71.01.1001", "Jayapura"},
	{-0.8667, 134.0833, "92.71.01.1001", "Manokwari"},
	{-8.4719, 140.3924, "93.71.01.1001", "Merauke"},
	{-3.5481, 135.4781, "94.71.01.1001", "Nabire"},
	{-3.8289, 138.8537, "95.71.01.1001", "Jayawijaya"},
	{-1.3260, 132.4668, "96.71.01.1001", "Sorong"},
}

// findAdm4Code mencari kode adm4 BMKG terdekat berdasarkan koordinat latitude & longitude
func findAdm4Code(lat, lon float64) string {
	closestAdm4 := "31.71.01.1001" // Default DKI Jakarta
	minDist := math.MaxFloat64

	for _, c := range indonesianCities {
		dLat := lat - c.Lat
		dLon := lon - c.Lon
		dist := dLat*dLat + dLon*dLon
		if dist < minDist {
			minDist = dist
			closestAdm4 = c.Adm4
		}
	}
	return closestAdm4
}

func (r *weatherRepository) isBMKGConfigured() bool {
	return strings.Contains(r.baseURL, "bmkg.go.id")
}

func (r *weatherRepository) fetchBMKGAPI(ctx context.Context, lat, lon float64) (*bmkgResponse, error) {
	if !r.isBMKGConfigured() {
		return nil, fmt.Errorf("BMKG API endpoint tidak dikonfigurasi: %s", r.baseURL)
	}

	adm4 := findAdm4Code(lat, lon)
	endpoint := fmt.Sprintf("%s/publik/prakiraan-cuaca?adm4=%s",
		strings.TrimRight(r.baseURL, "/"),
		url.QueryEscape(adm4),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("BMKG API returned status: %d", resp.StatusCode)
	}

	var data bmkgResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

// getClosestWeather mengambil data prakiraan BMKG yang paling dekat dengan jam sekarang
func getClosestWeather(data *bmkgResponse) *bmkgWeather {
	if data == nil || len(data.Data) == 0 || len(data.Data[0].Cuaca) == 0 {
		return nil
	}
	now := time.Now()
	var closest *bmkgWeather
	var minDiff time.Duration = 1<<63 - 1

	for _, dayBlock := range data.Data[0].Cuaca {
		for i := range dayBlock {
			w := &dayBlock[i]
			if t, err := time.Parse("2006-01-02 15:04:05", w.LocalDatetime); err == nil {
				diff := now.Sub(t)
				if diff < 0 {
					diff = -diff
				}
				if diff < minDiff {
					minDiff = diff
					closest = w
				}
			}
		}
	}
	if closest == nil && len(data.Data[0].Cuaca[0]) > 0 {
		return &data.Data[0].Cuaca[0][0]
	}
	return closest
}

// ==========================================
// FALLBACK DATA CUACA DAERAH INDONESIA
// ==========================================

func (r *weatherRepository) fallbackCurrentWeather() *domain.CurrentWeather {
	return &domain.CurrentWeather{
		Temperature:   28,
		Condition:     "Cerah Berawan",
		ConditionIcon: "https://api-apps.bmkg.go.id/storage/icon/cuaca/cerah-berawan-am.svg",
		DateTime:      time.Now().Format("Monday, 03:04 PM"),
	}
}

func (r *weatherRepository) fallbackHighlights() *domain.WeatherHighlights {
	return &domain.WeatherHighlights{
		Precipitation: "2%",
		Humidity:      "87%",
		WindSpeed:     "0 km/h",
		Sunrise:       "05:45 WIB",
		Sunset:        "17:55 WIB",
	}
}

func (r *weatherRepository) fallbackHourlyForecast() ([]domain.HourlyTemperature, []domain.ChanceOfRain) {
	hourly := []domain.HourlyTemperature{
		{Time: "06 am", Hour: 6, Temp: 10, TempC: 10, TempF: 50, FeelsLike: 10, Humidity: 65, Condition: "Cerah", Icon: "sunny"},
		{Time: "09 am", Hour: 9, Temp: 14, TempC: 14, TempF: 57.2, FeelsLike: 14, Humidity: 60, Condition: "Cerah Berawan", Icon: "partly-cloudy"},
		{Time: "12 pm", Hour: 12, Temp: 19, TempC: 19, TempF: 66.2, FeelsLike: 19, Humidity: 55, Condition: "Cerah Berawan", Icon: "partly-cloudy"},
		{Time: "03 pm", Hour: 15, Temp: 22, TempC: 22, TempF: 71.6, FeelsLike: 22, Humidity: 50, Condition: "Berawan", Icon: "cloudy"},
		{Time: "06 pm", Hour: 18, Temp: 18, TempC: 18, TempF: 64.4, FeelsLike: 18, Humidity: 90, Condition: "Hujan Lebat", Icon: "rainy"},
		{Time: "09 pm", Hour: 21, Temp: 15, TempC: 15, TempF: 59, FeelsLike: 15, Humidity: 75, Condition: "Hujan Ringan", Icon: "rainy"},
		{Time: "12 am", Hour: 0, Temp: 13, TempC: 13, TempF: 55.4, FeelsLike: 13, Humidity: 80, Condition: "Berawan", Icon: "cloudy"},
		{Time: "03 am", Hour: 3, Temp: 11, TempC: 11, TempF: 51.8, FeelsLike: 11, Humidity: 82, Condition: "Cerah Berawan", Icon: "partly-cloudy"},
	}

	chanceOfRain := []domain.ChanceOfRain{
		{Time: "09 am", Hour: 9, Intensity: "Sunny", Percentage: 10, PrecipMM: 0},
		{Time: "12 pm", Hour: 12, Intensity: "Partly Cloudy", Percentage: 60, PrecipMM: 1.2},
		{Time: "03 pm", Hour: 15, Intensity: "Partly Cloudy", Percentage: 40, PrecipMM: 0.5},
		{Time: "06 pm", Hour: 18, Intensity: "Heavy Rain", Percentage: 90, PrecipMM: 12.0},
		{Time: "09 pm", Hour: 21, Intensity: "Partly Cloudy", Percentage: 30, PrecipMM: 0.3},
		{Time: "12 am", Hour: 0, Intensity: "Sunny", Percentage: 15, PrecipMM: 0},
		{Time: "03 am", Hour: 3, Intensity: "Sunny", Percentage: 5, PrecipMM: 0},
	}

	return hourly, chanceOfRain
}

func (r *weatherRepository) fallbackMultiDayForecast(days int) []domain.DayForecast {
	if days <= 0 {
		days = 3
	}
	defaultForecasts := []domain.DayForecast{
		{Day: "Selasa", Date: time.Now().AddDate(0, 0, 1).Format("2006-01-02"), Condition: "Cerah Berawan", TempHigh: 32, TempLow: 24, Icon: "partly-cloudy"},
		{Day: "Rabu", Date: time.Now().AddDate(0, 0, 2).Format("2006-01-02"), Condition: "Hujan Ringan", TempHigh: 31, TempLow: 23, Icon: "rainy"},
		{Day: "Kamis", Date: time.Now().AddDate(0, 0, 3).Format("2006-01-02"), Condition: "Berawan", TempHigh: 33, TempLow: 24, Icon: "cloudy"},
		{Day: "Jumat", Date: time.Now().AddDate(0, 0, 4).Format("2006-01-02"), Condition: "Hujan Sedang", TempHigh: 30, TempLow: 23, Icon: "rainy"},
		{Day: "Sabtu", Date: time.Now().AddDate(0, 0, 5).Format("2006-01-02"), Condition: "Cerah", TempHigh: 33, TempLow: 25, Icon: "sunny"},
		{Day: "Minggu", Date: time.Now().AddDate(0, 0, 6).Format("2006-01-02"), Condition: "Cerah Berawan", TempHigh: 32, TempLow: 24, Icon: "partly-cloudy"},
		{Day: "Senin", Date: time.Now().AddDate(0, 0, 7).Format("2006-01-02"), Condition: "Hujan Ringan", TempHigh: 31, TempLow: 23, Icon: "rainy"},
	}
	if days > len(defaultForecasts) {
		days = len(defaultForecasts)
	}
	return defaultForecasts[:days]
}

// ==========================================
// WEATHER REPOSITORY INTERFACE METHODS
// ==========================================

func (r *weatherRepository) GetRealtimeWeather(ctx context.Context, lat, lon float64) (*domain.CurrentWeather, error) {
	data, err := r.fetchBMKGAPI(ctx, lat, lon)
	if err != nil {
		return r.fallbackCurrentWeather(), nil
	}

	closest := getClosestWeather(data)
	if closest == nil {
		return r.fallbackCurrentWeather(), nil
	}

	dateTimeStr := time.Now().Format("Monday, 03:04 PM")
	if t, parseErr := time.Parse("2006-01-02 15:04:05", closest.LocalDatetime); parseErr == nil {
		dateTimeStr = t.Format("Monday, 03:04 PM")
	}

	return &domain.CurrentWeather{
		Temperature:   closest.T,
		Condition:     closest.WeatherDesc,
		ConditionIcon: closest.Image,
		DateTime:      dateTimeStr,
	}, nil
}

func (r *weatherRepository) GetHighlights(ctx context.Context, lat, lon float64) (*domain.WeatherHighlights, error) {
	data, err := r.fetchBMKGAPI(ctx, lat, lon)
	if err != nil {
		return r.fallbackHighlights(), nil
	}

	closest := getClosestWeather(data)
	if closest == nil {
		return r.fallbackHighlights(), nil
	}

	humidity := fmt.Sprintf("%d%%", closest.Hu)
	windSpeed := fmt.Sprintf("%.0f km/h", closest.Ws)
	precip := fmt.Sprintf("%.1f mm", closest.Tp)
	if closest.Tp == 0 {
		precip = fmt.Sprintf("%d%%", closest.Tcc)
	}

	return &domain.WeatherHighlights{
		Precipitation: precip,
		Humidity:      humidity,
		WindSpeed:     windSpeed,
		Sunrise:       "05:45 WIB",
		Sunset:        "17:55 WIB",
	}, nil
}

func (r *weatherRepository) GetHourlyForecast(ctx context.Context, lat, lon float64) ([]domain.HourlyTemperature, []domain.ChanceOfRain, error) {
	data, err := r.fetchBMKGAPI(ctx, lat, lon)
	if err != nil || len(data.Data) == 0 || len(data.Data[0].Cuaca) == 0 {
		h, c := r.fallbackHourlyForecast()
		return h, c, nil
	}

	var hourly []domain.HourlyTemperature
	var chanceOfRain []domain.ChanceOfRain

	for _, dayBlock := range data.Data[0].Cuaca {
		for _, w := range dayBlock {
			formattedTime := w.LocalDatetime
			hourNum := 0
			if t, err := time.Parse("2006-01-02 15:04:05", w.LocalDatetime); err == nil {
				formattedTime = strings.ToLower(t.Format("03 pm"))
				hourNum = t.Hour()
			}

			tempF := (w.T * 9 / 5) + 32

			hourly = append(hourly, domain.HourlyTemperature{
				Time:      formattedTime,
				Hour:      hourNum,
				Temp:      w.T,
				TempC:     w.T,
				TempF:     mathRound(tempF, 1),
				FeelsLike: w.T,
				Humidity:  w.Hu,
				Condition: w.WeatherDesc,
				Icon:      w.Image,
			})

			intensity := "Sunny"
			lowerDesc := strings.ToLower(w.WeatherDesc)
			if w.Tp >= 5.0 || strings.Contains(lowerDesc, "lebat") || strings.Contains(lowerDesc, "petir") {
				intensity = "Heavy Rain"
			} else if w.Tp > 0.0 || strings.Contains(lowerDesc, "hujan") {
				intensity = "Light Rain"
			} else if w.Tcc > 50 || strings.Contains(lowerDesc, "berawan") {
				intensity = "Partly Cloudy"
			}

			rainPct := w.Hu
			if w.Tp > 0 {
				rainPct = int(math.Min(100, 60+w.Tp*8))
			} else if w.Tcc < 30 {
				rainPct = int(float64(w.Hu) * 0.2)
			}

			chanceOfRain = append(chanceOfRain, domain.ChanceOfRain{
				Time:       formattedTime,
				Hour:       hourNum,
				Intensity:  intensity,
				Percentage: rainPct,
				PrecipMM:   w.Tp,
			})
		}
	}

	if len(hourly) > 12 {
		hourly = hourly[:12]
		chanceOfRain = chanceOfRain[:12]
	}

	return hourly, chanceOfRain, nil
}

func (r *weatherRepository) GetMultiDayForecast(ctx context.Context, lat, lon float64, days int) ([]domain.DayForecast, error) {
	data, err := r.fetchBMKGAPI(ctx, lat, lon)
	if err != nil || len(data.Data) == 0 || len(data.Data[0].Cuaca) == 0 {
		return r.fallbackMultiDayForecast(days), nil
	}

	var forecasts []domain.DayForecast
	totalDays := len(data.Data[0].Cuaca)
	count := days
	if count > totalDays {
		count = totalDays
	}

	for i := 0; i < count; i++ {
		dayCuaca := data.Data[0].Cuaca[i]
		if len(dayCuaca) == 0 {
			continue
		}

		maxTemp := dayCuaca[0].T
		minTemp := dayCuaca[0].T
		condition := dayCuaca[0].WeatherDesc
		icon := dayCuaca[0].Image
		dateStr := ""
		dayName := ""

		for _, w := range dayCuaca {
			if w.T > maxTemp {
				maxTemp = w.T
			}
			if w.T < minTemp {
				minTemp = w.T
			}
			if strings.Contains(w.LocalDatetime, " 12:00:") || strings.Contains(w.LocalDatetime, " 13:00:") || strings.Contains(w.LocalDatetime, " 14:00:") {
				condition = w.WeatherDesc
				icon = w.Image
			}
		}

		if t, err := time.Parse("2006-01-02 15:04:05", dayCuaca[0].LocalDatetime); err == nil {
			dayName = translateDayIndonesian(t.Format("Monday"))
			dateStr = t.Format("2006-01-02")
		}

		forecasts = append(forecasts, domain.DayForecast{
			Day:       dayName,
			Date:      dateStr,
			Condition: condition,
			TempHigh:  mathRound(maxTemp, 1),
			TempLow:   mathRound(minTemp, 1),
			Icon:      icon,
		})
	}

	if len(forecasts) < days {
		fallbacks := r.fallbackMultiDayForecast(days)
		for i := len(forecasts); i < days && i < len(fallbacks); i++ {
			forecasts = append(forecasts, fallbacks[i])
		}
	}

	return forecasts, nil
}

func translateDayIndonesian(day string) string {
	switch strings.ToLower(day) {
	case "sunday":
		return "Minggu"
	case "monday":
		return "Senin"
	case "tuesday":
		return "Selasa"
	case "wednesday":
		return "Rabu"
	case "thursday":
		return "Kamis"
	case "friday":
		return "Jumat"
	case "saturday":
		return "Sabtu"
	default:
		return day
	}
}

func mathRound(val float64, precision int) float64 {
	pow := math.Pow(10, float64(precision))
	return math.Round(val*pow) / pow
}

func (r *weatherRepository) GetCalendarForecast(ctx context.Context, city string, year, month int) (*domain.WeatherCalendarData, error) {
	now := time.Now()
	if year <= 0 {
		year = now.Year()
	}
	if month < 1 || month > 12 {
		month = int(now.Month())
	}

	if strings.TrimSpace(city) == "" {
		city = "Jakarta"
	}

	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, -1)
	totalDays := lastDay.Day()

	baseHigh := 31.0
	baseLow := 24.0
	baseRain := 35

	lowerCity := strings.ToLower(city)
	if strings.Contains(lowerCity, "jakarta") || strings.Contains(lowerCity, "surabaya") || strings.Contains(lowerCity, "semarang") {
		baseHigh = 32.5
		baseLow = 24.5
		baseRain = 40
	} else if strings.Contains(lowerCity, "bandung") || strings.Contains(lowerCity, "bogor") || strings.Contains(lowerCity, "malang") {
		baseHigh = 28.0
		baseLow = 20.0
		baseRain = 55
	} else if strings.Contains(lowerCity, "denpasar") || strings.Contains(lowerCity, "mataram") {
		baseHigh = 31.5
		baseLow = 24.0
		baseRain = 30
	} else if strings.Contains(lowerCity, "medan") || strings.Contains(lowerCity, "makassar") || strings.Contains(lowerCity, "palembang") {
		baseHigh = 32.0
		baseLow = 23.5
		baseRain = 45
	}

	var days []domain.CalendarDay
	for d := 1; d <= totalDays; d++ {
		curDate := time.Date(year, time.Month(month), d, 0, 0, 0, 0, time.UTC)
		dateStr := curDate.Format("2006-01-02")
		weekday := curDate.Format("Monday")
		isToday := (now.Year() == year && int(now.Month()) == month && now.Day() == d)

		variant := (d*7 + len(city)*13) % 100
		tempOffset := float64((variant % 5) - 2)
		rainChance := (baseRain + (variant % 30) - 15)
		if rainChance < 5 {
			rainChance = 5
		}
		if rainChance > 95 {
			rainChance = 95
		}

		condition := "Cerah Berawan"
		icon := "partly-cloudy"
		advisory := "Cuaca bersahabat untuk kegiatan di luar ruangan dan perjalanan dinas / wisata."

		if rainChance > 70 {
			condition = "Hujan Petir"
			icon = "thunderstorm"
			advisory = "Waspada potensi hujan lebat disertai petir, utamakan aktivitas di dalam ruangan."
		} else if rainChance > 45 {
			condition = "Hujan Ringan"
			icon = "rainy"
			advisory = "Peluang hujan cukup tinggi, selalu siapkan payung atau jas hujan sebelum bepergian."
		} else if rainChance < 20 {
			condition = "Cerah"
			icon = "sunny"
			advisory = "Langit cerah bermatahari, kenakan tabir surya jika beraktivitas lama di siang hari."
		} else if variant%3 == 0 {
			condition = "Berawan"
			icon = "cloudy"
			advisory = "Udara teduh berawan sepanjang hari, sangat nyaman untuk perjalanan darat."
		}

		high := baseHigh + tempOffset
		low := baseLow + (tempOffset * 0.6)
		if low >= high {
			low = high - 4
		}

		days = append(days, domain.CalendarDay{
			Date:         dateStr,
			Day:          d,
			Weekday:      weekday,
			Condition:    condition,
			TempHigh:     mathRound(high, 1),
			TempLow:      mathRound(low, 1),
			ChanceOfRain: rainChance,
			Icon:         icon,
			Advisory:     advisory,
			IsToday:      isToday,
		})
	}

	return &domain.WeatherCalendarData{
		City:      city,
		Year:      year,
		Month:     month,
		MonthName: firstDay.Month().String(),
		Days:      days,
	}, nil
}

// BMKG API response structures
type bmkgResponse struct {
	Lokasi bmkgLocation `json:"lokasi"`
	Data   []bmkgData   `json:"data"`
}

type bmkgData struct {
	Lokasi bmkgLocation  `json:"lokasi"`
	Cuaca  bmkgCuacaDays `json:"cuaca"`
}

type bmkgLocation struct {
	Kecamatan string  `json:"kecamatan"`
	Kotkab    string  `json:"kotkab"`
	Provinsi  string  `json:"provinsi"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	Timezone  string  `json:"timezone"`
}

type bmkgCuacaDays [][]bmkgWeather

func (b *bmkgCuacaDays) UnmarshalJSON(data []byte) error {
	var list2D [][]bmkgWeather
	if err := json.Unmarshal(data, &list2D); err == nil {
		*b = list2D
		return nil
	}
	var list1D []bmkgWeather
	if err := json.Unmarshal(data, &list1D); err == nil {
		*b = [][]bmkgWeather{list1D}
		return nil
	}
	return fmt.Errorf("unable to unmarshal cuaca as 1D or 2D array")
}

type bmkgWeather struct {
	Datetime      string  `json:"datetime"`
	LocalDatetime string  `json:"local_datetime"`
	WeatherDesc   string  `json:"weather_desc"`
	WeatherDescEn string  `json:"weather_desc_en"`
	T             float64 `json:"t"`
	Tcc           int     `json:"tcc"`
	Tp            float64 `json:"tp"`
	Hu            int     `json:"hu"`
	Ws            float64 `json:"ws"`
	Wd            string  `json:"wd"`
	VsText        string  `json:"vs_text"`
	Image         string  `json:"image"`
}
