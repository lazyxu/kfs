import React from 'react';
import ContextMenu from './index';
import useContextMenu from "hox/contextMenu";
import useResourceManager from "hox/resourceManager";
import useDialog2 from "hox/dialog";

export default () => {
    const [contextMenu, setContextMenu] = useContextMenu();
    const [resourceManager, setResourceManager] = useResourceManager();
    let {filePath, branchName} = resourceManager;
    const [dialog, setDialog] = useDialog2();
    console.log('useDialog2()', useDialog2());
    if (contextMenu === null || contextMenu.type !== "default") {
        return <div/>
    }
    return <ContextMenu
        left={contextMenu.clientX}
        top={contextMenu.clientY}
        right={contextMenu.x + contextMenu.width}
        bottom={contextMenu.y + contextMenu.height}
        maxWidth={200}
        maxHeight={7*50}
        options={{
            上传文件: null,
            上传文件夹: null,
            新建文件: () => {
                console.log("新建文件1")
                console.log(dialog, setDialog, useDialog2, setContextMenu, useContextMenu, setResourceManager)
                setDialog({
                    title: "新建文件",
                })
                console.log("新建文件2")
            },
            新建文件夹: () => {
                setDialog({
                    title: "新建文件夹",
                })
            },
            刷新: null,
            粘贴: null,
            历史版本: {
                enabled: false, fn: () => {
                }
            },
            添加书签: null,
        }}
        onFinish={() => {
            // console.log("onFinish")
            setContextMenu(null)
        }
        }
    />
}
