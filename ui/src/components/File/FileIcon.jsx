import { Box } from "@mui/material";
import { openDir, openFile } from "api/fs";
import { isDCIM, isViewable, modeIsDir } from "api/utils/api";
import { getMetadata } from "api/web/exif";
import ImageViewer from "components/FileViewer/ImageViewer";
import VideoViewer from "components/FileViewer/VideoViewer";
import useResourceManager from "hox/resourceManager";
import { getSysConfig } from "hox/sysConfig";
import { useState } from "react";
import SvgIcon from "../Icon/SvgIcon";

export default function ({ dirItem }) {
    const sysConfig = getSysConfig().sysConfig;
    const [open, setOpen] = useState(false);
    const [metadata, setMetadata] = useState();
    const [resourceManager, setResourceManager] = useResourceManager();
    const { driver, filePath } = resourceManager;
    const { name, mode } = dirItem;
    const curFilePath = filePath.concat(name);
    return (
        <Box className="file-icon-box">
            {modeIsDir(mode) ?
                <SvgIcon icon="folder1" className='file-icon file-icon-folder' fontSize="inherit" onClick={() => {
                    openDir(setResourceManager, driver, curFilePath);
                }} /> :
                isDCIM(name) ?
                    <img src={`${sysConfig.webServer}/thumbnail?size=64&hash=${dirItem.hash}`} loading="lazy" onClick={() => getMetadata(dirItem.hash).then(m => {
                        setMetadata(m);
                        setOpen(true);
                    })} />
                    : name.toLowerCase().endsWith(".txt") ?
                        <SvgIcon icon="txt3" className='file-icon file-icon-file file-icon-file-viewable' fontSize="inherit" onClick={() => {
                            openFile(setResourceManager, driverId, curFilePath, dirItem);
                        }} />
                        : isViewable(name) ?
                            <SvgIcon icon="file12" className='file-icon file-icon-file file-icon-file-viewable' fontSize="inherit" onClick={() => {
                                openFile(setResourceManager, driverId, curFilePath, dirItem);
                            }} />
                            : <SvgIcon icon="file12" className='file-icon file-icon-file' fontSize="inherit" />
            }
            {open && metadata?.fileType?.type === "video" && <VideoViewer open={open} setOpen={setOpen} metadata={metadata} hash={metadata.hash} attribute={dirItem} />}
            {open && metadata?.fileType?.type === "image" && <ImageViewer open={open} setOpen={setOpen} metadata={metadata} hash={metadata.hash} attribute={dirItem} />}
        </Box>
    );
}
