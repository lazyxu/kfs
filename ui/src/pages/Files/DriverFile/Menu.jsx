import { Add, Bookmark, ContentPaste, CreateNewFolder, DriveFolderUpload, History, IosShare, Settings, UploadFile } from "@mui/icons-material";
import { ListItemText, MenuItem } from "@mui/material";
import Menu from "components/Menu";
import useResourceManager from "hox/resourceManager";

export default function ({ contextMenu, setContextMenu }) {
    const [resourceManager, setResourceManager] = useResourceManager();
    const { driver, filePath } = contextMenu;
    return (
        <Menu
            contextMenu={contextMenu}
            open={contextMenu !== null}
            onClose={() => setContextMenu(null)}
        >
            <MenuItem disabled>
                <UploadFile />
                <ListItemText>上传文件</ListItemText>
            </MenuItem>
            <MenuItem disabled>
                <DriveFolderUpload />
                <ListItemText>上传文件夹</ListItemText>
            </MenuItem>
            <MenuItem disabled>
                <Add />
                <ListItemText>新建文件</ListItemText>
            </MenuItem>
            <MenuItem disabled>
                <CreateNewFolder />
                <ListItemText>新建文件夹</ListItemText>
            </MenuItem>
            <MenuItem disabled>
                <ContentPaste />
                <ListItemText>粘贴</ListItemText>
            </MenuItem>
            <MenuItem disabled>
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
