import {download} from "../../../api/fs";
import humanize from 'humanize';
import useResourceManager from "../../../hox/resourceManager";
import {Link, Typography} from "@mui/material";
import {getSysConfig} from "../../../hox/sysConfig";

export default ({file}) => {
    const [resourceManager, setResourceManager] = useResourceManager();
    return (
        file.tooLarge ?
            <>
                文件大于{humanize.filesize(getSysConfig().sysConfig.maxContentSize)}，不支持在线查看，你可以选择
                <Link underline="hover" onClick={() => {
                    download(resourceManager.branchName, resourceManager.filePath)
                }}>下载该文件</Link>。
            </> :
            <Typography style={{whiteSpace: "pre-wrap", overflowWrap: "anywhere"}}>
                {file.content}
            </Typography>
    )
};
