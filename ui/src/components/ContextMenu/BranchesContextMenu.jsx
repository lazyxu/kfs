import React from 'react';
import ContextMenu from './index';
import useContextMenu from "hox/contextMenu";
import useResourceManager from "hox/resourceManager";
import useDialog2 from "hox/dialog";
import {listBranch} from "../../api/branch";

export default () => {
    const [contextMenu, setContextMenu] = useContextMenu();
    const [resourceManager, setResourceManager] = useResourceManager();
    let {filePath, branchName} = resourceManager;
    const [dialog, setDialog] = useDialog2();
    console.log('useDialog2()', useDialog2());
    if (contextMenu === null || contextMenu.type !== "branches") {
        return <div/>
    }
    return <ContextMenu
        left={contextMenu.clientX}
        top={contextMenu.clientY}
        right={contextMenu.x + contextMenu.width}
        bottom={contextMenu.y + contextMenu.height}
        maxWidth={200}
        maxHeight={2 * 50}
        options={{
            新建同步文件夹: () => {
                setDialog({
                    title: "新建同步文件夹",
                })
            },
            刷新: async () => {
                await listBranch(setResourceManager)
            },
        }}
        onFinish={() => {
            // console.log("onFinish")
            setContextMenu(null)
        }
        }
    />
}
