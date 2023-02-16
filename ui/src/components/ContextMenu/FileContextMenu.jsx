import React from 'react';
import ContextMenu from './index';
import useContextMenu from "hox/contextMenu";
import useResourceManager from "hox/resourceManager";
import {download, open} from "api/api";
import {modeIsDir} from "api/utils/api";
import useDialog from "hox/dialog";
import useDialog2 from "../../hox/dialog";

export default () => {
    const [contextMenu, setContextMenu] = useContextMenu();
    const [resourceManager, setResourceManager] = useResourceManager();
    const [dialog, setDialog] = useDialog2();
    if (contextMenu === null || contextMenu.type !== "file") {
        return <div/>
    }
    let {filePath, branchName} = resourceManager;
    let {Name, Mode } = contextMenu.dirItem;
    filePath = filePath.concat(Name);
    return <ContextMenu
        left={contextMenu.clientX}
        top={contextMenu.clientY}
        right={contextMenu.x + contextMenu.width}
        bottom={contextMenu.y + contextMenu.height}
        maxWidth={200}
        maxHeight={10*50}
        options={{
            打开: () => {
                open(setResourceManager, branchName, filePath);
            },
            下载: {
                enabled: !modeIsDir(Mode) , fn: () => {
                    download(branchName, filePath);
                }
            },
            分享: null,
            剪切: null,
            复制: null,
            删除: null,
            重命名: null,
            属性: () => {
                setDialog({
                    title: "属性",
                    dirItem: contextMenu.dirItem,
                })
            },
            历史版本: {
                enabled: false, fn: () => {
                }
            },
            添加书签: null,
        }}
        onFinish={() => setContextMenu(null)}
    />
}
