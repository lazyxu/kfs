import ImgCancelable from "@kfs/common/components/ImgCancelable";
import { getSysConfig } from "@kfs/common/hox/sysConfig";
import useWindows, { APP_IMAGE_VIEWER, APP_VIDEO_VIEWER, newWindow } from "@kfs/mui/hox/windows";
import { Box, Skeleton } from "@mui/material";
import moment from "moment";
import { useInView } from "react-intersection-observer";

export default function ({ metadata }) {
    const sysConfig = getSysConfig();
    let { hash, exif, fileType, shotTime, shotEquipment, videoMetadata } = metadata;
    let time = shotTime.isValid() ? shotTime.format("YYYY年MM月DD日 HH:mm:ss") : "未知时间";
    const [windows, setWindows] = useWindows();
    const { ref, inView } = useInView({ threshold: 0 });
    let title = "";
    if (fileType.type === "image") {
        title = time + "\n"
        + shotEquipment + "\n"
        + fileType.type + "/" + fileType.subType + "\n"
        + (exif.GPSLatitudeRef ? (exif.GPSLatitudeRef + " " + exif.GPSLatitude + "\n") : "")
        + (exif.GPSLongitudeRef ? (exif.GPSLongitudeRef + " " + exif.GPSLongitude + "\n") : "")
        + hash;
    }
    if (fileType.type === "video") {
        title = time + "\n"
        + shotEquipment + "\n"
        + fileType.type + "/" + fileType.subType + "\n"
        + (videoMetadata ? (videoMetadata.Codec + "\n") : "")
        + (videoMetadata ? (moment(videoMetadata.Created / 1000 / 1000).format("YYYY年MM月DD日 HH:mm:ss") + "\n") : "")
        + (videoMetadata ? (videoMetadata.Duration + "s\n") : "")
        + hash;
    }
    return (
        <Box ref={ref} sx={{ width: "100%", height: "100%" }}>
            {fileType.type === "image" && <ImgCancelable inView={inView}
                src={`${sysConfig.webServer}/thumbnail?size=256&cutSquare=true&hash=${hash}`}
                renderImg={(url) => <img style={{ width: "100%" }} title={title} src={url} loading="lazy" onClick={() =>
                    newWindow(setWindows, APP_IMAGE_VIEWER, { hash })
                } />}
                renderSkeleton={() => <Skeleton title={title} variant="rectangular" animation={false} width="100%" height="100%" />}
            />}
            {fileType.type === "video" && <ImgCancelable inView={inView} style={{ width: "100%" }}
                src={`${sysConfig.webServer}/thumbnail?size=256&cutSquare=true&hash=${hash}`}
                renderImg={(url) => <img style={{ width: "100%" }} title={title} src={url} loading="lazy" onClick={() =>
                    newWindow(setWindows, APP_VIDEO_VIEWER, { hash })
                } />}
                renderSkeleton={() => <Skeleton title={title} variant="rectangular" animation={false} width="100%" height="100%" />}
            />}
        </Box>
    );
}
