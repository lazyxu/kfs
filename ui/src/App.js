import { Inbox, Mail } from '@mui/icons-material';
import MenuIcon from '@mui/icons-material/Menu';
import { AppBar, Box, Divider, Drawer, IconButton, List, ListItem, ListItemButton, ListItemIcon, ListItemText, Toolbar, Typography, styled, useColorScheme } from "@mui/material";
import SvgIcon from 'components/Icon/SvgIcon';
import MetadataAnalysis from 'components/MetadataAnalysis';
import { SnackbarAction } from 'components/Notification/Notification';
import Version from 'components/Version';
import { SnackbarProvider } from 'notistack';
import Dcim from 'pages/Dcim';
import DedicatedSpace from 'pages/DedicatedSpace/DedicatedSpace';
import SystemConfig from 'pages/Setting/SystemConfig';
import React, { useEffect } from "react";
import useMenu from "./hox/menu";
import useSysConfig from "./hox/sysConfig";
import BackupTask from "./pages/BackupTask";
import Files from "./pages/Files";

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
        <SnackbarProvider action={SnackbarAction} >
            <Box sx={{
                position: 'fixed', width: "100%", height: "100%",
                display: 'flex', flexDirection: 'column',
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.primary
            }}>
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
                        <Box>
                            <MetadataAnalysis />
                        </Box>
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
                <DrawerHeader />
                {menu === '我的云盘' && <Files />}
                {menu === '我的相册' && <Dcim />}
                {menu === '备份任务' && <BackupTask />}
                {menu === '设置' && <SystemConfig />}
                {menu === '存储空间' && <DedicatedSpace />}
            </Box>
        </SnackbarProvider>
    );
}

export default App;
