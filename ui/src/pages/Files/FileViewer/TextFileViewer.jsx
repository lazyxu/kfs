import { Link, Typography } from "@mui/material";
import humanize from 'humanize';
import { download } from "../../../api/fs";
import useResourceManager from "../../../hox/resourceManager";
import { getSysConfig } from "../../../hox/sysConfig";

export default ({ file }) => {
    const [resourceManager, setResourceManager] = useResourceManager();
    console.log("fileViewer", file)
    return (
        file.tooLarge ?
            <>
                文件大于{humanize.filesize(getSysConfig().sysConfig.maxContentSize)}，不支持在线查看，你可以选择
                <Link underline="hover" onClick={() => {
                    download(resourceManager.driverId, resourceManager.filePath)
                }}>下载该文件</Link>。
            </> :
            <Typography sx={{ wordBreak: "break-all" }}>
                {file.content}
            </Typography>
    )
};
