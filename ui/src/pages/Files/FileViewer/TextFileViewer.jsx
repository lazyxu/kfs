import {download} from "../../../api/fs";
import humanize from 'humanize';
import useResourceManager from "../../../hox/resourceManager";
import {Link} from "@mui/material";
import {getSysConfig} from "../../../hox/sysConfig";
import Editor from "@monaco-editor/react";

export default ({file}) => {
    const [resourceManager, setResourceManager] = useResourceManager();
    console.log("fileViewer", file)
    return (
        file.tooLarge ?
            <>
                文件大于{humanize.filesize(getSysConfig().sysConfig.maxContentSize)}，不支持在线查看，你可以选择
                <Link underline="hover" onClick={() => {
                    download(resourceManager.driverName, resourceManager.filePath)
                }}>下载该文件</Link>。
            </> :
            <Editor
                height="100%"
                defaultValue={file.content}
                defaultPath={file.name}
            />
    )
};
