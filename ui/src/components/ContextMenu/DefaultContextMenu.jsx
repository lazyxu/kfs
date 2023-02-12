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
    if (contextMenu === null || contextMenu.type !== "default") {
        return <div/>
    }
    let {filePath, branchName} = resourceManager;
    return <ContextMenu
        left={contextMenu.clientX}
        top={contextMenu.clientY}
        right={contextMenu.x + contextMenu.width}
        bottom={contextMenu.y + contextMenu.height}
        maxWidth={200}
        maxHeight={150}
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
