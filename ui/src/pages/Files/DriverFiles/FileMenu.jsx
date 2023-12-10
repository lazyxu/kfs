import { Bookmark, ContentCopy, ContentCut, Delete, Download, DriveFileRenameOutline, History, IosShare, OpenInNew, Settings } from "@mui/icons-material";
import { ListItemText, MenuItem } from "@mui/material";
import { getDriverLocalFile } from "api/driver";
import { download } from "api/fs";
import { modeIsDir } from "api/utils/api";
import Menu from "components/Menu";
import useResourceManager, { openDir } from "hox/resourceManager";
import useWindows, { getOpenApp, newWindow } from "hox/windows";

export default function ({ contextMenu, setContextMenu, setFileAttribute }) {
    const [resourceManager, setResourceManager] = useResourceManager();
    const [windows, setWindows] = useWindows();
    const { driver, filePath, driverFile } = contextMenu;
    const { name, mode } = driverFile;
    return (
        <Menu
            contextMenu={contextMenu}
            open={contextMenu !== null}
            onClose={() => setContextMenu(null)}
        >
            <MenuItem onClick={() => {
                setContextMenu(null);
                if (modeIsDir(mode)) {
                    openDir(setResourceManager, driver, filePath);
                } else {
                    const app = getOpenApp(name);
                    newWindow(setWindows, app, { driver, filePath, driverFile });;
                }
            }}>
                <OpenInNew />
                <ListItemText>打开</ListItemText>
            </MenuItem>
            {process.env.REACT_APP_PLATFORM !== 'web' && <>
                <MenuItem onClick={() => {
                    setContextMenu(null);
                    getDriverLocalFile(driver.id).then(driverLocalFile => {
                        // console.log(driver, driverLocalFile, filePath);
                        const { shell } = window.require('@electron/remote');
                        shell.openPath(driverLocalFile.srcPath + "\\" + filePath.join("\\"));
                    });
                }}>
                    <OpenInNew />
                    <ListItemText>打开本地文件</ListItemText>
                </MenuItem>
                <MenuItem onClick={() => {
                    setContextMenu(null);
                    getDriverLocalFile(driver.id).then(driverLocalFile => {
                        // console.log(driver, driverLocalFile, filePath);
                        const { shell } = window.require('@electron/remote');
                        shell.showItemInFolder(driverLocalFile.srcPath + "\\" + filePath.join("\\"));
                    });
                }}>
                    <OpenInNew />
                    <ListItemText>打开本地文件位置</ListItemText>
                </MenuItem>
            </>}
            <MenuItem disabled={modeIsDir(mode)} onClick={() => {
                setContextMenu(null);
                download(driver.id, filePath);
            }}>
                <Download />
                <ListItemText>下载</ListItemText>
            </MenuItem>
            <MenuItem disabled>
                <ContentCut />
                <ListItemText>剪切</ListItemText>
            </MenuItem>
            <MenuItem disabled>
                <ContentCopy />
                <ListItemText>复制</ListItemText>
            </MenuItem>
            <MenuItem disabled>
                <Delete />
                <ListItemText>删除</ListItemText>
            </MenuItem>
            <MenuItem disabled>
                <DriveFileRenameOutline />
                <ListItemText>重命名</ListItemText>
            </MenuItem>
            <MenuItem onClick={() => {
                setContextMenu(null);
                setFileAttribute(contextMenu);
            }} >
                <Settings />
                <ListItemText>属性</ListItemText>
            </MenuItem>
            <MenuItem disabled>
                <IosShare />
                <ListItemText>分享</ListItemText>
            </MenuItem>
            <MenuItem disabled>
                <History />
                <ListItemText>历史版本</ListItemText>
            </MenuItem>
            <MenuItem disabled>
                <Bookmark />
                <ListItemText>添加书签</ListItemText>
            </MenuItem>
        </Menu>
    );
}
