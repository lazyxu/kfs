import { parseDateTime } from 'api/utils/api';

export default function ({ hash, exif }) {
    let time = parseDateTime(exif).format("YYYY年MM月DD日 HH:mm:ss");
    return (
        <img style={{ width: "100%" }} src={"http://127.0.0.1:1123/thumbnail?size=256&cutSquare=true&hash=" + hash} loading="lazy"
            title={time + " " + exif.OffsetTime + "\n"
                + (exif.Model ? (exif.Model + "\n") : "")
                + (exif.GPSLatitudeRef ? (exif.GPSLatitudeRef + " " + exif.GPSLatitude + "\n") : "")
                + (exif.GPSLongitudeRef ? (exif.GPSLongitudeRef + " " + exif.GPSLongitude + "\n") : "")
                + hash}
        />
    );
}
