import SvgIcon from "../Icon/SvgIcon";
import { modeIsDir, isDCIM, isViewable } from "api/utils/api";
import { Box, Stack } from "@mui/material";
import { getMetadata } from "api/web/exif";
import ImageViewer from "components/FileViewer/ImageViewer";
import { useState } from "react";
import VideoViewer from "components/FileViewer/VideoViewer";
import { getSysConfig } from "hox/sysConfig";
import { list, openFile } from "api/fs";
import useResourceManager from "hox/resourceManager";

export default function ({ dirItem }) {
    const sysConfig = getSysConfig().sysConfig;
    const [open, setOpen] = useState(false);
    const [metadata, setMetadata] = useState();
    const [resourceManager, setResourceManager] = useResourceManager();
    let { filePath, driverName } = resourceManager;
    filePath = filePath.concat(dirItem.name);
    return (
        <Box className="file-icon-box">
            {modeIsDir(dirItem.mode) ?
                <SvgIcon icon="folder1" className='file-icon file-icon-folder' fontSize="inherit" onClick={() => {
                    list(setResourceManager, driverName, filePath);
                }} /> :
                isDCIM(dirItem.name) ?
                    <img src={`${sysConfig.webServer}/thumbnail?size=64&hash=${dirItem.hash}`} loading="lazy" onClick={() => getMetadata(dirItem.hash).then(m => {
                        setMetadata(m);
                        setOpen(true);
                    })} />
                    : dirItem.name.toLowerCase().endsWith(".txt") ?
                        <SvgIcon icon="txt3" className='file-icon file-icon-file file-icon-file-viewable' fontSize="inherit" onClick={() => {
                            openFile(setResourceManager, driverName, filePath, dirItem);
                        }} />
                        : isViewable(dirItem.name) ?
                            <SvgIcon icon="file12" className='file-icon file-icon-file file-icon-file-viewable' fontSize="inherit" onClick={() => {
                                openFile(setResourceManager, driverName, filePath, dirItem);
                            }} />
                            : <SvgIcon icon="file12" className='file-icon file-icon-file' fontSize="inherit" />
            }
            {open && metadata?.fileType?.type === "video" && <VideoViewer open={open} setOpen={setOpen} metadata={metadata} hash={metadata.hash} attribute={dirItem} />}
            {open && metadata?.fileType?.type === "image" && <ImageViewer open={open} setOpen={setOpen} metadata={metadata} hash={metadata.hash} attribute={dirItem} />}
        </Box>
    );
}
