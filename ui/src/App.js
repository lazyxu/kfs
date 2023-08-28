import SystemConfig from 'pages/Setting/SystemConfig';
import Version from 'components/Version';
import Files from "./pages/Files";
import { AppBar, Box, Divider, Drawer, IconButton, List, ListItem, ListItemButton, ListItemIcon, ListItemText, Stack, Toolbar, Typography, styled, useColorScheme } from "@mui/material";
import React, { useEffect } from "react";
import useMenu from "./hox/menu";
import useSysConfig from "./hox/sysConfig";
import BackupTask from "./pages/BackupTask";
import Dcim from 'pages/Dcim';
import SvgIcon from 'components/Icon/SvgIcon';
import MenuIcon from '@mui/icons-material/Menu';
import { Inbox, Mail } from '@mui/icons-material';
import DedicatedSpace from 'pages/DedicatedSpace/DedicatedSpace';

function App() {
    const { sysConfig } = useSysConfig();
    const { menu, setMenu } = useMenu();
    const { mode, setMode } = useColorScheme();
    const [open, setOpen] = React.useState(false);
    const toggleDrawer = (open) => (event) => {
        if (event.type === 'keydown' && (event.key === 'Tab' || event.key === 'Shift')) {
            return;
        }

        setOpen(open);
    };
    useEffect(() => {
        // document.body.setAttribute('data-theme', sysConfig.theme);
        console.log("mode:", mode, "=>", sysConfig.theme);
        setMode(sysConfig.theme);
    }, [sysConfig.theme]);
    const DrawerHeader = styled('div')(({ theme }) => ({
        display: 'flex',
        alignItems: 'center',
        padding: theme.spacing(0, 1),
        // necessary for content to be below app bar
        ...theme.mixins.toolbar,
        justifyContent: 'flex-end',
    }));
    return (
        <Stack sx={{ width: "100%", height: "100%", position: "fixed" }} direction="row">
            <AppBar position="fixed" open={open}>
                <Toolbar>
                    <IconButton
                        color="inherit"
                        aria-label="open drawer"
                        onClick={() => setOpen(true)}
                        edge="start"
                        sx={{
                            marginRight: 5,
                            ...(open && { display: 'none' }),
                        }}
                    >
                        <MenuIcon />
                    </IconButton>
                    <Typography variant="h6" noWrap component="div">
                        {menu}
                    </Typography>
                </Toolbar>
            </AppBar>
            <Drawer
                anchor="left"
                open={open}
                onClose={toggleDrawer(false)}
            >
                <Box
                    sx={{ width: 250 }}
                    role="presentation"
                    onClick={toggleDrawer(false)}
                    onKeyDown={toggleDrawer(false)}
                >
                    <List>
                        {(process.env.REACT_APP_PLATFORM === 'web' ? [
                            { icon: 'wangpan', name: '我的云盘' },
                            { icon: 'DCIM', name: '我的相册' },
                            { icon: 'peizhi', name: '设置' },
                            { icon: 'equipment_data-02_fn', name: '存储空间' },
                        ] : [
                            { icon: 'wangpan', name: '我的云盘' },
                            { icon: 'DCIM', name: '我的相册' },
                            { icon: 'yuntongbu', name: '备份任务' },
                            { icon: 'peizhi', name: '设置' },
                            { icon: 'equipment_data-02_fn', name: '存储空间' },
                        ]).map((item, index) => (
                            <ListItem key={item.name} disablePadding onClick={() => setMenu(item.name)}>
                                <ListItemButton>
                                    <ListItemIcon>
                                        <SvgIcon icon={item.icon} style={{ height: "24px", width: "24px" }} />
                                    </ListItemIcon>
                                    <ListItemText primary={item.name} />
                                </ListItemButton>
                            </ListItem>
                        ))}
                    </List>
                    <Divider />
                    <List>
                        {[
                            { icon: 'devices', name: '设备列表' },
                            { icon: '', name: '文件类型' },
                            { icon: '', name: '文件大小' },
                            { icon: 'swapVertical', name: '传输列表' },
                            { icon: 'system', name: '资源监控' },
                            { icon: '', name: '我的书签' },
                            { icon: '', name: '分享历史' },
                            { icon: '', name: '操作历史' },
                            { icon: '', name: '搜索' },
                        ].map((item, index) => (
                            <ListItem key={item.name} disablePadding onClick={() => setMenu(item.name)}>
                                <ListItemButton>
                                    <ListItemIcon>
                                        <SvgIcon icon={item.icon} />
                                    </ListItemIcon>
                                    <ListItemText primary={item.name} />
                                </ListItemButton>
                            </ListItem>
                        ))}
                    </List>
                    <Divider />
                    <List>
                        {['All mail', 'Trash', 'Spam'].map((text, index) => (
                            <ListItem key={text} disablePadding>
                                <ListItemButton>
                                    <ListItemIcon>
                                        {index % 2 === 0 ? <Inbox /> : <Mail />}
                                    </ListItemIcon>
                                    <ListItemText primary={text} />
                                </ListItemButton>
                            </ListItem>
                        ))}
                    </List>
                </Box>
                <Version />
            </Drawer>
            <Box sx={{
                flex: 1,
                position: "relative",
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.primary
            }}>
                <DrawerHeader />
                <Files show={menu === '我的云盘'} />
                <Dcim show={menu === '我的相册'} />
                {process.env.REACT_APP_PLATFORM !== 'web' && <BackupTask show={menu === '备份任务'} />}
                <SystemConfig show={menu === '设置'} />
                <DedicatedSpace show={menu === '存储空间'} />
            </Box>
        </Stack>
    );
}

export default App;
