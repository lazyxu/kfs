import useResourceManager, { openDir } from '@kfs/common/hox/resourceManager';
import { deleteDriver, getDriverLocalFile, resetDriver } from '@kfs/mui/api/driver';
import SvgIcon from "@kfs/mui/components/Icon/SvgIcon";
import Menu from '@kfs/mui/components/Menu';
import { ClearAll, ContentCopy, Delete, DriveFileRenameOutline, OpenInNew, Settings } from '@mui/icons-material';
import { Box, ListItemText, MenuItem, Stack, Typography } from "@mui/material";
import { useState } from 'react';

export default ({ driver, setDriverAttribute, onDelete }) => {
    const [resourceManager, setResourceManager] = useResourceManager();
    const [contextMenu, setContextMenu] = useState(null);
    return (
        <>
            <Box
                title={`云盘名称：${driver.name}\n云盘描述：${driver.description}`}
                onContextMenu={(e) => {
                    e.preventDefault(); e.stopPropagation();
                    setContextMenu({ mouseX: e.clientX, mouseY: e.clientY });
                }}
                sx={{ padding: "1em", ":hover": { cursor: "pointer", backgroundColor: (theme) => theme.palette.action.hover } }}
            >
                <Stack
                    direction="row"
                    alignItems="center"
                    spacing={2}
                    onClick={() => openDir(setResourceManager, driver, [])}
                >
                    <Box sx={{ height: "64px", width: "64px" }}>
                        {driver.type === 'baiduPhoto' ?
                            <img src='baiduPhoto.png' style={{ maxWidth: "100%", maxHeight: "100%" }} />
                            : driver.type === 'localFile' ?
                                <SvgIcon icon="shangchuan" fontSize="inherit" style={{ height: "64px", width: "64px" }} />
                                :
                                <SvgIcon icon="wangpan" fontSize="inherit" style={{ height: "64px", width: "64px" }} />
                        }
                    </Box>
                    <Stack direction="column" sx={{ width: "10em" }}>
                        <Typography sx={{ whiteSpace: "nowrap", textOverflow: "ellipsis", overflow: "hidden", display: "block" }}>{driver.name}</Typography>
                        <Typography sx={{ flex: 1, whiteSpace: "nowrap", textOverflow: "ellipsis", overflow: "hidden", display: "block" }} color="text.secondary">
                            {driver.description}
                        </Typography>
                    </Stack>
                </Stack>
            </Box>
            <Menu
                contextMenu={contextMenu}
                open={contextMenu !== null}
                onClose={() => setContextMenu(null)}
            >
                <MenuItem onClick={() => openDir(setResourceManager, driver, [])}>
                    <OpenInNew />
                    <ListItemText>打开</ListItemText>
                </MenuItem>
                {/* TODO: device id */}
                {window.kfsEnv.VITE_APP_PLATFORM !== 'web' && driver.type === 'localFile' &&
                    <MenuItem onClick={() => {
                        setContextMenu(null);
                        getDriverLocalFile(driver.id).then(driverLocalFile => {
                            const { shell } = window.require('@electron/remote');
                            shell.showItemInFolder(driverLocalFile.srcPath);
                        });
                    }}>
                        <OpenInNew />
                        <ListItemText>打开本地文件位置</ListItemText>
                    </MenuItem>
                }
                <MenuItem onClick={() => resetDriver(driver.id).then(() => setContextMenu(null))} >
                    <ClearAll />
                    <ListItemText>重置</ListItemText>
                </MenuItem>
                <MenuItem disabled>
                    <ContentCopy />
                    <ListItemText>复制</ListItemText>
                </MenuItem>
                <MenuItem onClick={() => deleteDriver(driver.id).then(() => { setContextMenu(null); onDelete(); })} disableRipple>
                    <Delete />
                    <ListItemText>删除</ListItemText>
                </MenuItem>
                <MenuItem disabled>
                    <DriveFileRenameOutline />
                    <ListItemText>重命名</ListItemText>
                </MenuItem>
                <MenuItem onClick={() => { setContextMenu(null); setDriverAttribute(driver) }}>
                    <Settings />
                    <ListItemText>属性</ListItemText>
                </MenuItem>
            </Menu>
        </>
    )
};
