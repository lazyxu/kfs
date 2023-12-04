import { Bookmark, ContentCopy, ContentCut, Delete, Download, DriveFileRenameOutline, History, IosShare, OpenInNew, Settings } from "@mui/icons-material";
import { ListItemText, MenuItem } from "@mui/material";
import { download, openDir } from "api/fs";
import { modeIsDir } from "api/utils/api";
import Menu from "components/Menu";
import useResourceManager from "hox/resourceManager";

export default function ({ contextMenu, setContextMenu, setFileAttribute }) {
    const [resourceManager, setResourceManager] = useResourceManager();
    const { driver, filePath, mode } = contextMenu;
    return (
        <Menu
            contextMenu={contextMenu}
            open={contextMenu !== null}
            onClose={() => setContextMenu(null)}
        >
            <MenuItem onClick={() => {
                if (modeIsDir(mode)) {
                    setContextMenu(null);
                    openDir(setResourceManager, driver, filePath);
                } else {
                    // openFile(setResourceManager, driverId, filePath, dirItem);
                }
            }}>
                <OpenInNew />
                <ListItemText>打开</ListItemText>
            </MenuItem>
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
