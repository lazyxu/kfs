import React from 'react';
import ContextMenu from './index';
import useContextMenu from "../../hox/contextMenu";
import {newDir, newFile} from "../../api/api";
import useSysConfig from "../../hox/sysConfig";
import useResourceManager from "../../hox/resourceManager";

export default () => {
    const [contextMenu, setContextMenu] = useContextMenu();
    const [resourceManager, setResourceManager] = useResourceManager();
    const {sysConfig} = useSysConfig();
    if (contextMenu === null) {
        return <div/>
    }
    let left = contextMenu.clientX;
    let top = contextMenu.clientY;
    let maxWidth = 200;
    let maxHeight = 200;
    if (left + maxWidth > contextMenu.x + contextMenu.width) {
        left = contextMenu.x + contextMenu.width - maxWidth;
    }
    if (top + maxHeight > contextMenu.y + contextMenu.height) {
        top = contextMenu.y + contextMenu.height - maxHeight;
    }
    let {filePath, branchName} = resourceManager;
    return <ContextMenu
        left={left}
        top={top}
        options={{
            // 上传文件: <UploadFile/>,
            新建文件: () => {
                newFile(sysConfig, setResourceManager, branchName, filePath);
            },
            新建文件夹: () => {
                newDir(sysConfig, setResourceManager, branchName, filePath);
            },
            // 刷新: () => this.context.cd(branch, pwd),
            // 粘贴: {
            //     enabled: clipboard && clipboard.file,
            //     fn: () => {
            //         const {branch: srcBranch, pathList} = clipboard.file;
            //         if (clipboard.cut) {
            //             this.context.mv(srcBranch, pathList, branch, pwd);
            //             this.context.setState({clipboard: undefined});
            //             return;
            //         }
            //         this.context.cp(srcBranch, pathList, branch, pwd);
            //         this.context.setState({clipboard: undefined});
            //     },
            // },
            历史版本: {
                enabled: false, fn: () => {
                }
            },
        }}
        onFinish={() => setContextMenu(null)}
    />
}
