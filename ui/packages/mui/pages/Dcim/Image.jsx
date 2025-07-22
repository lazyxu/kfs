import ImgCancelable from "@kfs/common/components/ImgCancelable";
import { getSysConfig } from "@kfs/common/hox/sysConfig";
import useWindows, { APP_IMAGE_VIEWER, APP_VIDEO_VIEWER, newWindow } from "@kfs/mui/hox/windows";
import { Box, Skeleton } from "@mui/material";
import moment from "moment";
import { useInView } from "react-intersection-observer";

export default function ({ metadata }) {
    const sysConfig = getSysConfig();
    let { id, size, md5, sha256, time, shotEquipment } = metadata;
    const [windows, setWindows] = useWindows();
    const { ref, inView } = useInView({ threshold: 0 });
    let title = "";
    if (metadata.type === "image") {
        title = time + "\n"
        + shotEquipment + "\n"
        + metadata.type + "/" + metadata.sub_type + "\n"
        + (metadata.gps_latitude_ref ? (metadata.gps_latitude_ref + " " + metadata.gps_latitude + "\n") : "")
        + (metadata.gps_longitude_ref ? (metadata.gps_longitude_ref + " " + metadata.gps_longitude + "\n") : "")
        + id;
    }
    if (metadata.type === "video") {
        title = time + "\n"
        + shotEquipment + "\n"
        + metadata.type + "/" + metadata.sub_type + "\n"
        + (metadata.codec ? (metadata.codec + "\n") : "")
        + (metadata.created ? (metadata.created + "\n") : "")
        + (metadata.duration ? (metadata.duration + "s\n") : "")
        + id;
    }
    return (
        <Box ref={ref} sx={{ width: "100%", height: "100%" }}>
            {metadata.type === "image" && <ImgCancelable inView={inView}
                src={`${sysConfig.webServer}/thumbnail?id=${id}&size=${size}&md5=${md5}&sha256=${sha256}&hw=256&cutSquare=true`}
                renderImg={(url) => <img style={{ width: "100%" }} title={title} src={url} loading="lazy" onClick={() =>
                    newWindow(setWindows, APP_IMAGE_VIEWER, { metadata })
                } />}
                renderSkeleton={() => <Skeleton title={title} variant="rectangular" animation={false} width={256} height={256} />}
            />}
            {metadata.type === "video" && <ImgCancelable inView={inView} style={{ width: "100%" }}
                src={`${sysConfig.webServer}/thumbnail?id=${id}&size=${size}&md5=${md5}&sha256=${sha256}&hw=256&cutSquare=true`}
                renderImg={(url) => <img style={{ width: "100%" }} title={title} src={url} loading="lazy" onClick={() =>
                    newWindow(setWindows, APP_VIDEO_VIEWER, { metadata })
                } />}
                renderSkeleton={() => <Skeleton title={title} variant="rectangular" animation={false} width={256} height={256} />}
            />}
        </Box>
    );
}
