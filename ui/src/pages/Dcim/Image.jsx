import { useState } from "react";
import styles from './image.module.scss';
import ImageViewer from "components/FileViewer/ImageViewer";
import { getTransform } from "api/utils/api";
import moment from "moment";

export default function ({ metadata }) {
    const [open, setOpen] = useState(false);
    let { hash, exif, fileType, shotTime, shotEquipment, videoMetadata } = metadata;
    let time = shotTime.isValid() ? shotTime.format("YYYY年MM月DD日 HH:mm:ss") : "未知时间";
    return (
        <div>
            <img style={{ width: "100%", transform: getTransform(exif.Orientation) }} className={styles.clickable}
                src={`http://127.0.0.1:1123/thumbnail?size=256&cutSquare=true&hash=${hash}`} loading="lazy"
                title={time + "\n"
                    + shotEquipment + "\n"
                    + fileType.type + "/" + fileType.subType + "\n"
                    + (exif.GPSLatitudeRef ? (exif.GPSLatitudeRef + " " + exif.GPSLatitude + "\n") : "")
                    + (exif.GPSLongitudeRef ? (exif.GPSLongitudeRef + " " + exif.GPSLongitude + "\n") : "")
                    + (videoMetadata ? (videoMetadata.Codec + "\n") : "")
                    + (videoMetadata ? (moment(videoMetadata.Created / 1000 / 1000).format("YYYY年MM月DD日 HH:mm:ss") + "\n") : "")
                    + (videoMetadata ? (videoMetadata.Duration + "s\n") : "")
                    + hash}
                onClick={() => setOpen(true)}
            />
            <ImageViewer open={open} setOpen={setOpen} metadata={metadata} hash={hash} />
        </div>
    );
}
