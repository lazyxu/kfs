export default function ({ hash, exif }) {
    let time = exif.shotTime.isValid() ? exif.shotTime.format("YYYY年MM月DD日 HH:mm:ss") : "未知时间";
    let transform = "";
    switch (exif.Orientation) {
        case 2:
            transform = "rotateY(180deg)";
            break;
        case 3:
            transform = "rotate(180deg)";
            break;
        case 4:
            transform = "rotate(180deg)rotateY(180deg)";
            break;
        case 5:
            transform = "rotate(270deg)rotateY(180deg)";
            break;
        case 6:
            transform = "rotate(90deg)";
            break;
        case 7:
            transform = "rotate(90deg)rotateY(180deg)";
            break;
        case 8:
            transform = "rotate(270deg)";
            break;
    }
    return (
        <img style={{ width: "100%", transform }} src={"http://127.0.0.1:1123/thumbnail?size=256&cutSquare=true&hash=" + hash} loading="lazy"
            title={time + "\n"
                + (exif.Model ? (exif.Model + "\n") : "")
                + (exif.GPSLatitudeRef ? (exif.GPSLatitudeRef + " " + exif.GPSLatitude + "\n") : "")
                + (exif.GPSLongitudeRef ? (exif.GPSLongitudeRef + " " + exif.GPSLongitude + "\n") : "")
                + hash}
        />
    );
}
