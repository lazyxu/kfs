import {download} from "../../../api/fs";
import useResourceManager from "../../../hox/resourceManager";
import {Link} from "@mui/material";

const size10M = 10 * 1024 * 1024;

export default ({file}) => {
    const [resourceManager, setResourceManager] = useResourceManager();
    return (
        file.size < size10M ?
            <div>
                {file.content}
            </div> :
            <>
                文件大于100MB，不支持在线查看，你可以选择<Link underline="hover" onClick={() => {
                download(resourceManager.branchName, resourceManager.filePath)
            }}>下载该文件</Link>。
            </>
    )
};
