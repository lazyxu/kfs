package dao

type ExifData struct {
	Version         string
	DateTime        uint64
	HostComputer    string
	OffsetTime      string
	GPSLatitudeRef  string
	GPSLatitude     float64
	GPSLongitudeRef string
	GPSLongitude    float64
}
