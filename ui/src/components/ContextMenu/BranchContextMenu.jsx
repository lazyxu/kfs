import { list } from "api/fs";
import useContextMenu from "hox/contextMenu";
import useResourceManager from "hox/resourceManager";
import { deleteDriver } from "../../api/driver";
import useDialog2 from "../../hox/dialog";
import ContextMenu from './index';

export default function () {
    const [contextMenu, setContextMenu] = useContextMenu();
    const [resourceManager, setResourceManager] = useResourceManager();
    const [dialog, setDialog] = useDialog2();
    if (contextMenu === null || contextMenu.type !== "driver") {
        return <div />
    }
    let { driver } = contextMenu;
    return <ContextMenu
        left={contextMenu.clientX}
        top={contextMenu.clientY}
        right={contextMenu.x + contextMenu.width}
        bottom={contextMenu.y + contextMenu.height}
        maxWidth={200}
        maxHeight={10 * 50}
        options={{
            打开: async () => {
                await list(setResourceManager, driver.id, driver.name, []);
            },
            删除: async () => {
                await deleteDriver(setResourceManager, driver.id);
            },
            重命名: null,
            属性: () => {
                setDialog({
                    title: "属性",
                    driver,
                })
            },
        }}
        onFinish={() => setContextMenu(null)}
    />
}
