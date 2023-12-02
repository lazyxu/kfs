import { ClearAll, ContentCopy, Delete, DriveFileRenameOutline, OpenInNew, Settings } from '@mui/icons-material';
import { Box, ListItemText, MenuItem, Stack, Typography } from "@mui/material";
import { deleteDriver, resetDriver } from 'api/driver';
import { openDir } from "api/fs";
import SvgIcon from "components/Icon/SvgIcon";
import Menu from 'components/Menu';
import useResourceManager from 'hox/resourceManager';
import { useState } from 'react';

export default ({ driver, setDriverAttribute }) => {
    const [resourceManager, setResourceManager] = useResourceManager();
    const [contextMenu, setContextMenu] = useState(null);
    return (
        <>
            <Box
                title={`云盘名称：${driver.name}\n云盘描述：${driver.description}`}
                onContextMenu={(e) => {
                    e.preventDefault();
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
                <MenuItem onClick={() => { setContextMenu(null); setDriverAttribute(driver) }}>
                    <Settings />
                    <ListItemText>属性</ListItemText>
                </MenuItem>
                <MenuItem disabled>
                    <DriveFileRenameOutline />
                    <ListItemText>重命名</ListItemText>
                </MenuItem>
                <MenuItem disabled>
                    <ContentCopy />
                    <ListItemText>复制</ListItemText>
                </MenuItem>
                <MenuItem onClick={() => resetDriver(driver.id).then(() => setContextMenu(null))} >
                    <ClearAll />
                    <ListItemText>重置</ListItemText>
                </MenuItem>
                <MenuItem onClick={() => deleteDriver(setResourceManager, driver.id).then(() => setContextMenu(null))} disableRipple>
                    <Delete />
                    <ListItemText>删除</ListItemText>
                </MenuItem>
            </Menu>
        </>
    )
};
