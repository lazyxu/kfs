import { getSysConfig } from "hox/sysConfig";
import useWindows, { APP_IMAGE_VIEWER, APP_VIDEO_VIEWER, newWindow } from "hox/windows";
import moment from "moment";
import styles from './image.module.scss';

export default function ({ metadata }) {
    const sysConfig = getSysConfig().sysConfig;
    let { hash, exif, fileType, shotTime, shotEquipment, videoMetadata } = metadata;
    let time = shotTime.isValid() ? shotTime.format("YYYY年MM月DD日 HH:mm:ss") : "未知时间";
    const [windows, setWindows] = useWindows();
    return (
        <>
            {fileType.type === "image" && <img style={{ width: "100%" }} className={styles.clickable}
                src={`${sysConfig.webServer}/thumbnail?size=256&cutSquare=true&hash=${hash}`} loading="lazy"
                title={time + "\n"
                    + shotEquipment + "\n"
                    + fileType.type + "/" + fileType.subType + "\n"
                    + (exif.GPSLatitudeRef ? (exif.GPSLatitudeRef + " " + exif.GPSLatitude + "\n") : "")
                    + (exif.GPSLongitudeRef ? (exif.GPSLongitudeRef + " " + exif.GPSLongitude + "\n") : "")
                    + hash}
                onClick={() => newWindow(setWindows, APP_IMAGE_VIEWER, { hash })}
            />}
            {fileType.type === "video" && <img style={{ width: "100%" }} className={styles.clickable}
                src={`${sysConfig.webServer}/thumbnail?size=256&cutSquare=true&hash=${hash}`} loading="lazy"
                title={time + "\n"
                    + shotEquipment + "\n"
                    + fileType.type + "/" + fileType.subType + "\n"
                    + (videoMetadata ? (videoMetadata.Codec + "\n") : "")
                    + (videoMetadata ? (moment(videoMetadata.Created / 1000 / 1000).format("YYYY年MM月DD日 HH:mm:ss") + "\n") : "")
                    + (videoMetadata ? (videoMetadata.Duration + "s\n") : "")
                    + hash}
                onClick={() => newWindow(setWindows, APP_VIDEO_VIEWER, { hash })}
            />}
        </>
    );
}
