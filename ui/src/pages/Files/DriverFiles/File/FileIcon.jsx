import { Box } from "@mui/material";
import { getMetadata } from "api/exif";
import { isDCIM, isViewable, modeIsDir } from "api/utils/api";
import ImageViewer from "components/FileViewer/ImageViewer";
import VideoViewer from "components/FileViewer/VideoViewer";
import useResourceManager, { openDir } from "hox/resourceManager";
import { getSysConfig } from "hox/sysConfig";
import useWindows, { APP_TEXT_VIEWER, newWindow } from "hox/windows";
import { memo, useState } from "react";
import SvgIcon from "../../../../components/Icon/SvgIcon";
import ImgCancelable from "./ImgCancelable";

export default memo(({ driver, filePath, dirItem, hasBeenInView, inView }) => {
    const sysConfig = getSysConfig().sysConfig;
    const [open, setOpen] = useState(false);
    const [metadata, setMetadata] = useState();
    const [resourceManager, setResourceManager] = useResourceManager();
    const [windows, setWindows] = useWindows();
    const { name, mode } = dirItem;
    // console.log("===render", filePath, hasBeenInView)
    return (
        <Box className="file-icon-box">
            {modeIsDir(mode) ?
                <SvgIcon icon="folder1" className='file-icon file-icon-folder' fontSize="inherit" onClick={() => {
                    openDir(setResourceManager, driver, filePath);
                }} /> :
                isDCIM(name) ?
                    <ImgCancelable src={`${sysConfig.webServer}/thumbnail?size=64&hash=${dirItem.hash}`} inView={inView} onClick={() => getMetadata(dirItem.hash).then(m => {
                        setMetadata(m);
                        setOpen(true);
                    })} />
                    : name.toLowerCase().endsWith(".txt") ?
                        <SvgIcon icon="txt3" className='file-icon file-icon-file file-icon-file-viewable' fontSize="inherit" onClick={() => {
                            newWindow(setWindows, APP_TEXT_VIEWER, { driver, filePath, dirItem });
                        }} />
                        : isViewable(name) ?
                            <SvgIcon icon="file12" className='file-icon file-icon-file file-icon-file-viewable' fontSize="inherit" onClick={() => {
                                newWindow(setWindows, APP_TEXT_VIEWER, { driver, filePath, dirItem });
                            }} />
                            : <SvgIcon icon="file12" className='file-icon file-icon-file' fontSize="inherit" />
            }
            {open && metadata?.fileType?.type === "video" && <VideoViewer open={open} setOpen={setOpen} metadata={metadata} hash={metadata.hash} attribute={dirItem} />}
            {open && metadata?.fileType?.type === "image" && <ImageViewer open={open} setOpen={setOpen} metadata={metadata} hash={metadata.hash} attribute={dirItem} />}
        </Box>
    );
});
