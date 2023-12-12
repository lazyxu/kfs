import { getSysConfig } from "@/hox/sysConfig";
import useWindows, { APP_IMAGE_VIEWER, APP_VIDEO_VIEWER, newWindow } from "@/hox/windows";
import ImgCancelable from "@/pages/Files/DriverFiles/File/ImgCancelable";
import { Box } from "@mui/material";
import moment from "moment";
import { useInView } from "react-intersection-observer";
import styles from './image.module.scss';

export default function ({ metadata }) {
    const sysConfig = getSysConfig().sysConfig;
    let { hash, exif, fileType, shotTime, shotEquipment, videoMetadata } = metadata;
    let time = shotTime.isValid() ? shotTime.format("YYYY年MM月DD日 HH:mm:ss") : "未知时间";
    const [windows, setWindows] = useWindows();
    const { ref, inView } = useInView({ threshold: 0 });
    return (
        <Box ref={ref} sx={{ width: "100%", height: "100%" }}>
            {fileType.type === "image" && <ImgCancelable inView={inView} style={{ width: "100%" }} className={styles.clickable}
                src={`${sysConfig.webServer}/thumbnail?size=256&cutSquare=true&hash=${hash}`} loading="lazy"
                title={time + "\n"
                    + shotEquipment + "\n"
                    + fileType.type + "/" + fileType.subType + "\n"
                    + (exif.GPSLatitudeRef ? (exif.GPSLatitudeRef + " " + exif.GPSLatitude + "\n") : "")
                    + (exif.GPSLongitudeRef ? (exif.GPSLongitudeRef + " " + exif.GPSLongitude + "\n") : "")
                    + hash}
                onClick={() => newWindow(setWindows, APP_IMAGE_VIEWER, { hash })}
            />}
            {fileType.type === "video" && <ImgCancelable inView={inView} style={{ width: "100%" }} className={styles.clickable}
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
        </Box>
    );
}
