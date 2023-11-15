import { Close } from '@mui/icons-material';
import { Box, Dialog, DialogContent, DialogTitle, Divider, Grid, MenuItem, Select, Switch } from "@mui/material";
import IconButton from "@mui/material/IconButton";
import { getDriverLocalFile, getDriverSync, getDriversDirCount, getDriversFileCount, getDriversFileSize, updateDriverSync } from 'api/web/driver';
import { noteError } from 'components/Notification/Notification';
import humanize from 'humanize';
import moment from "moment/moment";
import { useEffect, useState } from 'react';
import DriverBaiduPhoto from './DriverBaiduPhoto';
import DriverLocalFile from './DriverLocalFile';

function formatTime(t) {
    return moment(t / 1000 / 1000).format("YYYY年MM月DD日 HH:mm:ss");
}

function Attr({ k, children }) {
    return <>
        <Grid xs={4} item sx={{ overflowWrap: "anywhere" }}><Box>{k}：</Box></Grid>
        <Grid xs={8} item sx={{ overflowWrap: "anywhere" }}>{children}</Grid>
    </>
}

function getDriverType(driver) {
    switch (driver.type) {
        case "baiduPhoto":
            return "一刻相册备份盘";
        case "localFile":
            return "本地文件备份盘";
        case "":
            return "普通云盘";
        default:
            break;
    }
}

export default ({ setOpen, driver }) => {
    // TODO: get more calculated attributes from server.
    const [attributes, setAttributes] = useState({});
    const [syncAttributes, setSyncAttributes] = useState();
    const [localFileAttributes, setLocalFileAttributes] = useState();
    useEffect(() => {
        if (driver.type === "baiduPhoto" || driver.type === "localFile") {
            getDriverSync(driver.id).then(n => setSyncAttributes(n)).catch(e => noteError(e.message));
        }
        if (driver.type === "localFile") {
            getDriverLocalFile(driver.id).then(n => setLocalFileAttributes(n));
        }
        getDriversFileSize(driver.id).then(n => setAttributes(prev => { return { ...prev, fileSize: n }; })).catch(e => noteError(e.message));
        getDriversFileCount(driver.id).then(n => setAttributes(prev => { return { ...prev, fileCount: n }; })).catch(e => noteError(e.message));
        getDriversDirCount(driver.id).then(n => setAttributes(prev => { return { ...prev, dirCount: n }; })).catch(e => noteError(e.message));
    }, []);
    return (
        <Dialog open={true} fullWidth={true} onClose={() => setOpen(false)}>
            <DialogTitle sx={{
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.secondary
            }}>
                云盘属性
                <IconButton
                    aria-label="close"
                    onClick={() => setOpen(false)}
                    sx={{
                        position: 'absolute',
                        right: 8,
                        top: 8,
                        color: (theme) => theme.palette.grey[500],
                    }}
                >
                    <Close />
                </IconButton>
            </DialogTitle>
            <DialogContent sx={{
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.primary
            }}>
                <Grid container spacing={1.5} sx={{ alignItems: "center" }}>
                    <Attr k="id">{driver.id}</Attr>
                    <Attr k="名称">{driver.name}</Attr>
                    <Attr k="描述">{driver.description}</Attr>
                    <Attr k="类型">{getDriverType(driver)}</Attr>
                    <Attr k="总大小">{humanize.filesize(attributes.fileSize)}</Attr>
                    <Attr k="文件数量">{attributes.fileCount}</Attr>
                    <Attr k="目录数量">{attributes.dirCount}</Attr>
                    {/* <Box variant="body"> */}
                    {/* <Box>文件总数：{driver.count}</Box> */}
                    {/* <Box>总大小：{humanize.filesize(driver.size)}</Box> */}
                    {/* <Typography>可修改该云盘的设备：any</Typography> */}
                    {/* <Typography>可读取该云盘的设备：any</Typography> */}
                    {/* </Box> */}
                    {driver.type === "baiduPhoto" && <>
                        <Grid xs={12} item sx={{ overflowWrap: "anywhere" }}><Divider /></Grid>
                        <Attr k="同步"><DriverBaiduPhoto driver={driver} /></Attr>
                        <Attr k="定时同步">
                            {syncAttributes ? <>
                                <Switch checked={syncAttributes.sync} onChange={e => updateDriverSync(driver.id, e.target.checked, syncAttributes.h, syncAttributes.m).then(setSyncAttributes(prev => { return { ...prev, sync: e.target.checked }; })).catch(e => noteError(e.message))} />
                                <Select variant="standard" size="small" sx={{ marginLeft: "1em" }} value={syncAttributes.h} onChange={e => updateDriverSync(driver.id, syncAttributes.sync, e.target.value, syncAttributes.m).then(setSyncAttributes(prev => { return { ...prev, h: e.target.value }; })).catch(e => noteError(e.message))}>
                                    {[...Array(24).keys()].map(value =>
                                        <MenuItem key={value} value={value}>{value.toString().padStart(2, 0)}</MenuItem>
                                    )}
                                </Select>时
                                <Select variant="standard" size="small" sx={{ marginLeft: "1em" }} value={syncAttributes.m} onChange={e => updateDriverSync(driver.id, syncAttributes.sync, syncAttributes.h, e.target.value).then(setSyncAttributes(prev => { return { ...prev, m: e.target.value }; })).catch(e => noteError(e.message))}>
                                    {[...Array(60).keys()].map(value =>
                                        <MenuItem key={value} value={value}>{value.toString().padStart(2, 0)}</MenuItem>
                                    )}
                                </Select>分
                            </> : <>配置加载中...</>}
                        </Attr>
                    </>}
                    {driver.type === "localFile" && <>
                        <Grid xs={12} item sx={{ overflowWrap: "anywhere" }}><Divider /></Grid>
                        <Attr k="设备ID">{localFileAttributes ? localFileAttributes.deviceId : "加载中..."}</Attr>
                        <Attr k="本地文件夹路径">{localFileAttributes ? localFileAttributes.srcPath : "加载中..."}</Attr>
                        <Attr k="上传时压缩">{localFileAttributes ? localFileAttributes.encoder : "加载中..."}</Attr>
                        <Grid xs={12} item sx={{ overflowWrap: "anywhere" }}><Divider /></Grid>
                        <DriverLocalFile driver={driver} attributes={localFileAttributes} />
                        <Attr k="定时同步">
                            {syncAttributes ? <>
                                <Switch checked={syncAttributes.sync} onChange={e => updateDriverSync(driver.id, e.target.checked, syncAttributes.h, syncAttributes.m).then(setSyncAttributes(prev => { return { ...prev, sync: e.target.checked }; })).catch(e => noteError(e.message))} />
                                <Select variant="standard" size="small" sx={{ marginLeft: "1em" }} value={syncAttributes.h} onChange={e => updateDriverSync(driver.id, syncAttributes.sync, e.target.value, syncAttributes.m).then(setSyncAttributes(prev => { return { ...prev, h: e.target.value }; })).catch(e => noteError(e.message))}>
                                    {[...Array(24).keys()].map(value =>
                                        <MenuItem key={value} value={value}>{value.toString().padStart(2, 0)}</MenuItem>
                                    )}
                                </Select>时
                                <Select variant="standard" size="small" sx={{ marginLeft: "1em" }} value={syncAttributes.m} onChange={e => updateDriverSync(driver.id, syncAttributes.sync, syncAttributes.h, e.target.value).then(setSyncAttributes(prev => { return { ...prev, m: e.target.value }; })).catch(e => noteError(e.message))}>
                                    {[...Array(60).keys()].map(value =>
                                        <MenuItem key={value} value={value}>{value.toString().padStart(2, 0)}</MenuItem>
                                    )}
                                </Select>分
                            </> : <>配置加载中...</>}
                        </Attr>
                    </>}
                </Grid>
            </DialogContent>
        </Dialog>
    );
};
