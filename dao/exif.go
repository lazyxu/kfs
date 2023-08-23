package dao

type Exif struct {
	ExifVersion         string
	ImageDescription    string
	Orientation         uint16 // 方向
	DateTime            string // 修改时间 YYYY:MM:DD HH:MM:SS
	DateTimeOriginal    string // 拍摄时间
	DateTimeDigitized   string // 写入时间
	OffsetTime          string // 时区 +01:00
	OffsetTimeOriginal  string
	OffsetTimeDigitized string
	SubsecTime          string // 亚秒 长度不确定
	SubsecTimeOriginal  string
	SubsecTimeDigitized string
	HostComputer        string
	Make                string
	Model               string
	ExifImageWidth      uint64
	ExifImageLength     uint64
	GPSLatitudeRef      string
	GPSLatitude         float64 // 纬度
	GPSLongitudeRef     string
	GPSLongitude        float64 // 经度
}

type ExifWithFileType struct {
	Hash     string   `json:"hash"`
	Exif     Exif     `json:"exif"`
	FileType FileType `json:"fileType"`
}
