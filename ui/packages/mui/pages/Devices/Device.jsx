import useSysConfig from '@kfs/common/hox/sysConfig';
import { deleteDevice } from '@kfs/mui/api/device';
import SvgIcon from "@kfs/mui/components/Icon/SvgIcon";
import { Delete, DriveFileRenameOutline, MoreVert } from '@mui/icons-material';
import { Box, Card, CardContent, IconButton, ListItemText, Menu, MenuItem, Stack, alpha, styled } from "@mui/material";
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

export default ({ device, setDevices }) => {
    const { sysConfig, setSysConfig } = useSysConfig();
    const [anchorEl, setAnchorEl] = useState(null);
    let deviceIcon = "weizhishebei";
    let os = device.os.toLowerCase();
    if (os.includes("windows")) {
        deviceIcon = "windows";
    }
    if (os.includes("linux")) {
        deviceIcon = "linux";
    }
    if (os.includes("mac")) {
        deviceIcon = "MAC";
    }
    if (os.includes("android")) {
        deviceIcon = "Android";
    }
    if (os.includes("iphone")) {
        deviceIcon = "mobileios";
    }
    if (device.hostname === "") {
        deviceIcon = "HTML";
    }
    return (
        <Card sx={{ minWidth: 275 }} variant="outlined">
            <CardContent>
                <Stack
                    direction="row"
                    alignItems="center"
                    spacing={2}
                >
                    <Box sx={{ height: "64px", width: "64px" }} >
                        <SvgIcon icon={deviceIcon} fontSize="inherit" style={{ height: "64px", width: "64px" }} />
                    </Box>
                    <Stack
                        direction="column"
                        sx={{ flex: 1 }}
                    >
                        <Box>名称：{device.name}</Box>
                        <Box>id：{device.id} {sysConfig.deviceId === device.id && "（当前设备）"}</Box>
                        <Box color="text.secondary" sx={{ flex: 1 }}>
                            系统：{device.os}
                        </Box>
                        <Box color="text.secondary" sx={{ flex: 1 }}>
                            userAgent：{device.userAgent}
                        </Box>
                        <Box color="text.secondary" sx={{ flex: 1 }}>
                            hostname：{device.hostname}
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
                        <MenuItem disabled>
                            <DriveFileRenameOutline />
                            <ListItemText>重命名</ListItemText>
                        </MenuItem>
                        <MenuItem onClick={() => { setAnchorEl(null); deleteDevice(setDevices, device.id); }} disableRipple>
                            <Delete />
                            <ListItemText>删除</ListItemText>
                        </MenuItem>
                    </StyledMenu>
                </Stack>
            </CardContent>
        </Card>
    )
};
