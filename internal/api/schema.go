package api

// QueryTypeInfo describes a query type and its available metrics.
type QueryTypeInfo struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Metrics     []string `json:"metrics"`
}

// Schema holds all discoverable API parameter information.
type Schema struct {
	QueryTypes []QueryTypeInfo    `json:"query_types"`
	DeviceTypes []string          `json:"device_types"`
	FunNames    []string          `json:"fun_names"`
	FilterOps   []string          `json:"filter_ops"`
	FilterNames FilterNamesSchema `json:"filter_names"`
}

// FilterNamesSchema separates fixed and dynamic filter names.
type FilterNamesSchema struct {
	Fixed   []string `json:"fixed"`
	Dynamic []string `json:"dynamic"`
}

var queryTypes = []QueryTypeInfo{
	{
		Name:        "page_metrics",
		Description: "页面整体指标查询",
		Metrics:     []string{"visits", "pv", "uv", "newVisitsRate", "entrances", "fvRate", "timeOnPage", "clicks", "clickRate", "ctaClicks", "ctaClickRate", "bounceRate", "avgPageViews", "completions", "conversionRate"},
	},
	{
		Name:        "page_insight",
		Description: "页面分组洞察（需要 funName 参数）",
		Metrics:     []string{"visits", "pv", "uv", "newVisitsRate", "entrances", "fvRate", "timeOnPage", "clicks", "clickRate", "ctaClicks", "ctaClickRate", "bounceRate", "avgPageViews", "completions", "conversionRate"},
	},
	{
		Name:        "block_metrics",
		Description: "区块指标（需页面已配置区块）",
		Metrics:     []string{"impression", "impressionRate", "dropoff", "dropoffRate", "avgDuration", "completions", "conversionRate"},
	},
	{
		Name:        "element_metrics",
		Description: "元素指标（需页面已配置元素）",
		Metrics:     []string{"impression", "impressionRate", "click", "clickRate", "completions", "conversionRate"},
	},
}

// ValidQueryTypes returns valid query type names.
var ValidQueryTypes = []string{"page_metrics", "page_insight", "block_metrics", "element_metrics"}

// ValidDeviceTypes returns valid device type values.
var ValidDeviceTypes = []string{"ALL", "PC", "MOBILE", "TABLET"}

// ValidFunNames returns valid funName values for page_insight.
var ValidFunNames = []string{
	"terminalType", "sourceType", "visitType", "aiName",
	"utmCampaign", "utmSource", "utmMedium", "utmTerm", "utmContent",
	"week", "day",
}

// ValidFilterOps returns valid filter operations.
var ValidFilterOps = []string{"include", "exclude"}

var fixedFilterNames = []string{"deviceType", "sourceType", "visitType", "exitType"}

var dynamicFilterNames = []string{
	"os", "osVersion", "browser", "browserVersion", "screenResolution", "deviceBrand",
	"country", "region", "searchEngine", "socialNetwork", "socialUrl", "aiName",
	"referralSource", "referralUrl", "campaignUrl",
	"utmCampaign", "utmSource", "utmMedium", "utmTerm", "utmContent",
	"combinedPages", "originalPages", "conversionName", "eventName",
	"customDimension", "eventVariable",
}

// GetSchema returns the full API schema for the describe command.
func GetSchema() *Schema {
	return &Schema{
		QueryTypes:  queryTypes,
		DeviceTypes: ValidDeviceTypes,
		FunNames:    ValidFunNames,
		FilterOps:   ValidFilterOps,
		FilterNames: FilterNamesSchema{
			Fixed:   fixedFilterNames,
			Dynamic: dynamicFilterNames,
		},
	}
}

// GetMetricsForQueryType returns available metrics for a specific query type.
func GetMetricsForQueryType(qt string) []string {
	for _, q := range queryTypes {
		if q.Name == qt {
			return q.Metrics
		}
	}
	return nil
}
