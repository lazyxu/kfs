import { Box } from "@mui/material";
import { getMetadata } from "api/exif";
import { isImage, isVideo, isViewable, modeIsDir } from "api/utils/api";
import VideoViewer from "components/FileViewer/VideoViewer";
import useResourceManager, { openDir } from "hox/resourceManager";
import { getSysConfig } from "hox/sysConfig";
import useWindows, { APP_IMAGE_VIEWER, APP_TEXT_VIEWER, newWindow } from "hox/windows";
import { memo, useState } from "react";
import SvgIcon from "../../../../components/Icon/SvgIcon";
import ImgCancelable from "./ImgCancelable";

export default memo(({ driver, filePath, driverFile, hasBeenInView, inView }) => {
    const sysConfig = getSysConfig().sysConfig;
    const [open, setOpen] = useState(false);
    const [metadata, setMetadata] = useState();
    const [resourceManager, setResourceManager] = useResourceManager();
    const [windows, setWindows] = useWindows();
    const { name, mode } = driverFile;
    // console.log("===render", filePath, hasBeenInView)
    return (
        <Box className="file-icon-box">
            {modeIsDir(mode) ?
                <SvgIcon icon="folder1" className='file-icon file-icon-folder' fontSize="inherit" onClick={() => {
                    openDir(setResourceManager, driver, filePath);
                }} /> :
                isImage(name) ?
                    <ImgCancelable src={`${sysConfig.webServer}/thumbnail?size=64&hash=${driverFile.hash}`} inView={inView} onClick={() => {
                        newWindow(setWindows, APP_IMAGE_VIEWER, { driver, filePath });
                    }} /> :
                    isVideo(name) ?
                        <ImgCancelable src={`${sysConfig.webServer}/thumbnail?size=64&hash=${driverFile.hash}`} inView={inView} onClick={() => getMetadata(driverFile.hash).then(m => {
                            setMetadata(m);
                            setOpen(true);
                        })} />
                        : name.toLowerCase().endsWith(".txt") ?
                            <SvgIcon icon="txt3" className='file-icon file-icon-file file-icon-file-viewable' fontSize="inherit" onClick={() => {
                                newWindow(setWindows, APP_TEXT_VIEWER, { driver, filePath });
                            }} />
                            : isViewable(name) ?
                                <SvgIcon icon="file12" className='file-icon file-icon-file file-icon-file-viewable' fontSize="inherit" onClick={() => {
                                    newWindow(setWindows, APP_TEXT_VIEWER, { driver, filePath });
                                }} />
                                : <SvgIcon icon="file12" className='file-icon file-icon-file' fontSize="inherit" />
            }
            {open && metadata?.fileType?.type === "video" && <VideoViewer open={open} setOpen={setOpen} metadata={metadata} hash={metadata.hash} attribute={driverFile} />}
        </Box>
    );
});
