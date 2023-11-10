import { ClearAll, ContentCopy, Delete, DriveFileRenameOutline, MoreVert, OpenInNew, Settings } from '@mui/icons-material';
import { Box, Card, CardContent, IconButton, Link, ListItemText, Menu, MenuItem, Stack, alpha, styled } from "@mui/material";
import { deleteDriver } from 'api/driver';
import { list } from "api/fs";
import SvgIcon from "components/Icon/SvgIcon";
import useContextMenu from "hox/contextMenu";
import useResourceManager from 'hox/resourceManager';
import { useState } from 'react';

const StyledMenu = styled((props) => (
    <Menu
        elevation={0}
        anchorOrigin={{
            vertical: 'bottom',
            horizontal: 'right',
        }}
        transformOrigin={{
            vertical: 'top',
            horizontal: 'right',
        }}
        {...props}
    />
))(({ theme }) => ({
    '& .MuiPaper-root': {
        borderRadius: 6,
        marginTop: theme.spacing(1),
        minWidth: 180,
        color:
            theme.palette.mode === 'light' ? 'rgb(55, 65, 81)' : theme.palette.grey[300],
        boxShadow:
            'rgb(255, 255, 255) 0px 0px 0px 0px, rgba(0, 0, 0, 0.05) 0px 0px 0px 1px, rgba(0, 0, 0, 0.1) 0px 10px 15px -3px, rgba(0, 0, 0, 0.05) 0px 4px 6px -2px',
        '& .MuiMenu-list': {
            padding: '4px 0',
        },
        '& .MuiMenuItem-root': {
            '& .MuiSvgIcon-root': {
                fontSize: 18,
                color: theme.palette.text.secondary,
                marginRight: theme.spacing(1.5),
            },
            '&:active': {
                backgroundColor: alpha(
                    theme.palette.primary.main,
                    theme.palette.action.selectedOpacity,
                ),
            },
        },
    },
}));

export default ({ driver, setDriverAttribute }) => {
    const [anchorEl, setAnchorEl] = useState(null);
    const [resourceManager, setResourceManager] = useResourceManager();
    const [contextMenu, setContextMenu] = useContextMenu();
    return (
        <Card sx={{ minWidth: 275 }} variant="outlined">
            <CardContent>
                <Stack
                    direction="row"
                    alignItems="center"
                    spacing={2}
                >
                    <Link onClick={() => list(setResourceManager, driver.id, driver.name, [])}>
                        <Box sx={{ height: "64px", width: "64px" }} >
                            {driver.type === 'baiduPhoto' ?
                                <img src='baiduPhoto.png' style={{ maxWidth: "100%", maxHeight: "100%" }} />
                                :
                                <SvgIcon icon="wangpan" fontSize="inherit" style={{ height: "64px", width: "64px" }} />
                            }
                        </Box>
                    </Link>
                    <Stack
                        direction="column"
                        sx={{ flex: 1 }}
                    >
                        <Link underline="hover" onClick={() => list(setResourceManager, driver.id, driver.name, [])} sx={{ flex: 1 }}>
                            <Box>{driver.name}</Box>
                        </Link>
                        <Box color="text.secondary" sx={{ flex: 1 }}>
                            {driver.description}
                        </Box>
                    </Stack>
                    <IconButton aria-label="settings" onClick={e => {
                        setAnchorEl(e.currentTarget);
                    }}>
                        <MoreVert />
                    </IconButton>
                    <StyledMenu
                        anchorEl={anchorEl}
                        open={Boolean(anchorEl)}
                        onClose={() => setAnchorEl(null)}
                    >
                        <MenuItem onClick={() => list(setResourceManager, driver.id, driver.name, [])}>
                            <OpenInNew />
                            <ListItemText>打开</ListItemText>
                        </MenuItem>
                        <MenuItem onClick={() => { setAnchorEl(null); setDriverAttribute(driver) }}>
                            <Settings />
                            <ListItemText>设置</ListItemText>
                        </MenuItem>
                        <MenuItem disabled>
                            <DriveFileRenameOutline />
                            <ListItemText>重命名</ListItemText>
                        </MenuItem>
                        <MenuItem disabled>
                            <ContentCopy />
                            <ListItemText>复制</ListItemText>
                        </MenuItem>
                        <MenuItem disabled>
                            <ClearAll />
                            <ListItemText>重置</ListItemText>
                        </MenuItem>
                        <MenuItem onClick={() => { setAnchorEl(null); deleteDriver(setResourceManager, driver.id); }} disableRipple>
                            <Delete />
                            <ListItemText>删除</ListItemText>
                        </MenuItem>
                    </StyledMenu>
                </Stack>
            </CardContent>
        </Card>
    )
};
