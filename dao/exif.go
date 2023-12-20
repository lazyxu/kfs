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

type VideoMetadata struct {
	Codec    string
	Created  int64
	Modified int64
	Duration float64 // 容器中媒体数据的持续时间（秒）
}

// VideoMetadataFfmpeg ffprobe.exe -v quiet -show_format -show_streams -print_format json 9638.mp4
type VideoMetadataFfmpeg struct {
	Height   string  // streams[0].height
	Width    string  // streams[0].width
	Created  int64   // format.tags.creation_time
	Duration float64 // format.duration
	Make     string  // format.tags. com.apple.quicktime.make
	Model    string  // format.tags. com.apple.quicktime.modal
}

type HeightWidth struct {
	Width  uint64
	Height uint64
}

type Metadata struct {
	Hash          string         `json:"hash"`
	FileType      *FileType      `json:"fileType"`
	Time          int64          `json:"time"`
	Year          int64          `json:"year"`
	Month         int64          `json:"month"`
	Day           int64          `json:"day"`
	HeightWidth   *HeightWidth   `json:"heightWidth"`
	Exif          *Exif          `json:"exif"`
	VideoMetadata *VideoMetadata `json:"videoMetadata"`
}
