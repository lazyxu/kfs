import { Box } from "@mui/material";
import { modeIsDir } from "api/utils/api";
import useResourceManager, { openDir } from "hox/resourceManager";
import { getSysConfig } from "hox/sysConfig";
import useWindows, { APP_IMAGE_VIEWER, APP_TEXT_VIEWER, APP_VIDEO_VIEWER, getOpenApp, newWindow } from "hox/windows";
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
    if (modeIsDir(mode)) {
        return (
            <Box className="file-icon-box">
                <SvgIcon icon="folder1" className='file-icon file-icon-folder' fontSize="inherit" onClick={() => {
                    openDir(setResourceManager, driver, filePath);
                }} />
            </Box>
        );
    }
    const app = getOpenApp(name);
    const onClick = () => newWindow(setWindows, app, { driver, filePath, driverFile });
    return (
        <Box className="file-icon-box">
            {
                app === APP_IMAGE_VIEWER || app === APP_VIDEO_VIEWER ?
                    <ImgCancelable src={`${sysConfig.webServer}/thumbnail?size=64&hash=${driverFile.hash}`} inView={inView} onClick={onClick} /> :
                    app === APP_TEXT_VIEWER ?
                        <SvgIcon icon="txt3" className='file-icon file-icon-file file-icon-file-viewable' fontSize="inherit" onClick={onClick} />
                        : <SvgIcon icon="file12" className='file-icon file-icon-file' fontSize="inherit" />
            }
        </Box>
    );
});
