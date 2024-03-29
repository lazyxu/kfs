import useWindows, { APP_IMAGE_VIEWER, APP_LIVP_UNZIP, APP_METADATA_MANAGER, APP_TEXT_VIEWER, APP_VIDEO_VIEWER } from "@kfs/mui/hox/windows";
import ImageViewer from "./ImageViewer";
import LivpUnZip from "./LivpUnZip";
import MetadataManager from "./MetadataManager";
import TextViewer from "./TextViewer";
import VideoViewer from "./VideoViewer";

export default function () {
    const [windows, setWindows] = useWindows();
    console.log(windows)
    return (
        Object.values(windows).map(w => (
            <div key={w.id}>
                {w.app === APP_TEXT_VIEWER && <TextViewer id={w.id} props={w.props} />}
                {w.app === APP_IMAGE_VIEWER && <ImageViewer id={w.id} props={w.props} />}
                {w.app === APP_VIDEO_VIEWER && <VideoViewer id={w.id} props={w.props} />}
                {w.app === APP_METADATA_MANAGER && <MetadataManager id={w.id} props={w.props} />}
                {w.app === APP_LIVP_UNZIP && <LivpUnZip id={w.id} props={w.props} />}
            </div>
        ))
    )
}
