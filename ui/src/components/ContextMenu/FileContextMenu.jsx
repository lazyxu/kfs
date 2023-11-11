import { download, list, openFile } from "api/fs";
import { modeIsDir } from "api/utils/api";
import useContextMenu from "hox/contextMenu";
import useResourceManager from "hox/resourceManager";
import useDialog2 from "../../hox/dialog";
import ContextMenu from './index';

export default () => {
    const [contextMenu, setContextMenu] = useContextMenu();
    const [resourceManager, setResourceManager] = useResourceManager();
    const [dialog, setDialog] = useDialog2();
    if (contextMenu === null || contextMenu.type !== "file") {
        return <div />
    }
    let { driverId, driverName, filePath } = resourceManager;
    let { dirItem } = contextMenu;
    let { name, mode } = dirItem;
    filePath = filePath.concat(name);
    return <ContextMenu
        left={contextMenu.clientX}
        top={contextMenu.clientY}
        right={contextMenu.x + contextMenu.width}
        bottom={contextMenu.y + contextMenu.height}
        maxWidth={200}
        maxHeight={10 * 50}
        options={{
            打开: async () => {
                if (modeIsDir(mode)) {
                    await list(setResourceManager, driverId, driverName, filePath);
                } else {
                    await openFile(setResourceManager, driverId, filePath, dirItem);
                }
            },
            下载: {
                enabled: !modeIsDir(mode), fn: async () => {
                    await download(driverId, filePath);
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
                    dirItem,
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
