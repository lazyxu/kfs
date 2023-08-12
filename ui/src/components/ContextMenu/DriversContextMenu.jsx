import React from 'react';
import ContextMenu from './index';
import useContextMenu from "hox/contextMenu";
import useResourceManager from "hox/resourceManager";
import useDialog2 from "hox/dialog";
import {listDriver} from "../../api/driver";

export default () => {
    const [contextMenu, setContextMenu] = useContextMenu();
    const [resourceManager, setResourceManager] = useResourceManager();
    let {filePath, driverName} = resourceManager;
    const [dialog, setDialog] = useDialog2();
    console.log('useDialog2()', useDialog2());
    if (contextMenu === null || contextMenu.type !== "drivers") {
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
            新建云盘: () => {
                setDialog({
                    title: "新建云盘",
                })
            },
            刷新: async () => {
                await listDriver(setResourceManager)
            },
        }}
        onFinish={() => {
            // console.log("onFinish")
            setContextMenu(null)
        }
        }
    />
}
