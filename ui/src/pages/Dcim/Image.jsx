import { useState } from "react";
import styles from './image.module.scss';
import ImageViewer from "components/FileViewer/ImageViewer";
import { getTransform } from "api/utils/api";
import moment from "moment";
import VideoViewer from "components/FileViewer/VideoViewer";
import { getSysConfig } from "hox/sysConfig";

export default function ({ metadata }) {
    const sysConfig = getSysConfig().sysConfig;
    const [open, setOpen] = useState(false);
    let { hash, exif, fileType, shotTime, shotEquipment, videoMetadata } = metadata;
    let time = shotTime.isValid() ? shotTime.format("YYYY年MM月DD日 HH:mm:ss") : "未知时间";
    return (
        <>
            {fileType.type == "image" && <img style={{ width: "100%", transform: getTransform(exif.Orientation) }} className={styles.clickable}
                src={`${sysConfig.webServer}/thumbnail?size=256&cutSquare=true&hash=${hash}`} loading="lazy"
                title={time + "\n"
                    + shotEquipment + "\n"
                    + fileType.type + "/" + fileType.subType + "\n"
                    + (exif.GPSLatitudeRef ? (exif.GPSLatitudeRef + " " + exif.GPSLatitude + "\n") : "")
                    + (exif.GPSLongitudeRef ? (exif.GPSLongitudeRef + " " + exif.GPSLongitude + "\n") : "")
                    + hash}
                onClick={() => setOpen(true)}
            />}
            {fileType.type == "video" && <img style={{ width: "100%" }} className={styles.clickable}
                src={`${sysConfig.webServer}/thumbnail?size=256&cutSquare=true&hash=${hash}`} loading="lazy"
                title={time + "\n"
                    + shotEquipment + "\n"
                    + fileType.type + "/" + fileType.subType + "\n"
                    + (videoMetadata ? (videoMetadata.Codec + "\n") : "")
                    + (videoMetadata ? (moment(videoMetadata.Created / 1000 / 1000).format("YYYY年MM月DD日 HH:mm:ss") + "\n") : "")
                    + (videoMetadata ? (videoMetadata.Duration + "s\n") : "")
                    + hash}
                onClick={() => setOpen(true)}
            />}
            {fileType.type == "video" && <VideoViewer open={open} setOpen={setOpen} metadata={metadata} hash={hash} />}
            {fileType.type == "image" && <ImageViewer open={open} setOpen={setOpen} metadata={metadata} hash={hash} />}
        </>
    );
}
