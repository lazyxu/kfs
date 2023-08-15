package dao

type Exif struct {
	Version         string `json:"version"`
	DateTime        uint64 `json:"dateTime"`
	HostComputer    string `json:"hostComputer"`
	OffsetTime      string `json:"offsetTime"`
	GPSLatitudeRef  string
	GPSLatitude     float64
	GPSLongitudeRef string
	GPSLongitude    float64
}
