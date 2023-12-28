import { modeIsDir } from "@kfs/common/api/utils";
import ImgCancelable from "@kfs/common/components/ImgCancelable";
import useResourceManager, { openDir } from "@kfs/common/hox/resourceManager";
import { getSysConfig } from "@kfs/common/hox/sysConfig";
import useWindows, { APP_IMAGE_VIEWER, APP_TEXT_VIEWER, APP_VIDEO_VIEWER, getOpenApp, newWindow } from "@kfs/mui/hox/windows";
import { Box, Skeleton } from "@mui/material";
import { memo } from "react";
import SvgIcon from "../../../../components/Icon/SvgIcon";

export default memo(({ driver, filePath, driverFile, hasBeenInView, inView }) => {
    const sysConfig = getSysConfig();
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
                    <ImgCancelable src={`${sysConfig.webServer}/thumbnail?size=64&hash=${driverFile.hash}`} inView={inView}
                        renderImg={(url) => <img title={name} src={url} loading="lazy" onClick={onClick} />}
                        renderSkeleton={() => <Skeleton title={name} variant="rectangular" animation={false} width="100%" height="100%" />}
                    /> :
                    app === APP_TEXT_VIEWER ?
                        <SvgIcon icon="txt3" className='file-icon file-icon-file file-icon-file-viewable' fontSize="inherit" onClick={onClick} />
                        : <SvgIcon icon="file12" className='file-icon file-icon-file' fontSize="inherit" />
            }
        </Box>
    );
});
