package common

// DataBaseConfig database info
type DataBaseConfig struct {
	DataBaseMode        string
	DataBaseReplsetName string
	DataBaseConnURLs    []string
	DataBaseName        string
	UseAuthentication   bool
	User                string
	Pwd                 string
}
