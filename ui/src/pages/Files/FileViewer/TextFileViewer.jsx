import {download} from "../../../api/fs";
import useResourceManager from "../../../hox/resourceManager";
import {Link} from "@mui/material";

const size10M = 10 * 1024 * 1024;

export default ({file}) => {
    const [resourceManager, setResourceManager] = useResourceManager();
    return (
        file.Size < size10M ?
            <>
                {(new TextDecoder("utf-8")).decode((file.Content))}
            </> :
            <>
                文件大于100MB，不支持在线查看，你可以选择<Link underline="hover" onClick={() => {
                download(resourceManager.branchName, resourceManager.filePath)
            }}>下载该文件</Link>。
            </>
    )
};

