import { Box, Divider, Grid, Input } from "@mui/material";
import { getDriverLocalFile, getDriversDirCalculatedInfo } from 'api/web/driver';
import humanize from 'humanize';
import moment from "moment/moment";
import { useEffect, useState } from 'react';

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

export default ({ driver }) => {
    const [attributes, setAttributes] = useState({});
    const [localFileAttributes, setLocalFileAttributes] = useState();
    useEffect(() => {
        if (driver.type === "localFile") {
            getDriverLocalFile(driver.id).then(attrs => setLocalFileAttributes(attrs));
        }
        getDriversDirCalculatedInfo(driver.id).then(setAttributes);
    }, []);
    return (
        <Grid container spacing={1.5} sx={{ alignItems: "center" }}>
            <Attr k="云盘ID">{driver.id}</Attr>
            <Attr k="云盘名称">{driver.name}</Attr>
            <Attr k="云盘描述">{driver.description}</Attr>
            <Attr k="云盘类型">{getDriverType(driver)}</Attr>
            <Attr k="总大小">{humanize.filesize(attributes.fileSize)}</Attr>
            <Attr k="文件数量">{attributes.fileCount}</Attr>
            <Attr k="目录数量">{attributes.dirCount}</Attr>
            {/* <Box variant="body"> */}
            {/* <Box>文件总数：{driver.count}</Box> */}
            {/* <Box>总大小：{humanize.filesize(driver.size)}</Box> */}
            {/* <Typography>可修改该云盘的设备：any</Typography> */}
            {/* <Typography>可读取该云盘的设备：any</Typography> */}
            {/* </Box> */}
            {driver.type === "localFile" && <>
                <Grid xs={12} item sx={{ overflowWrap: "anywhere" }}><Divider /></Grid>
                <Attr k="设备ID">{localFileAttributes ? localFileAttributes.deviceId : "加载中..."}</Attr>
                <Attr k="本地文件夹路径">{localFileAttributes ?
                    <a title={localFileAttributes.srcPath} onClick={() => {
                        const { shell } = window.require('@electron/remote');
                        shell.openPath(localFileAttributes.srcPath);
                    }} >{localFileAttributes.srcPath}</a> : "加载中..."}</Attr>
                <Attr k="上传时压缩">{localFileAttributes ? localFileAttributes.encoder : "加载中..."}</Attr>
                <Attr k="过滤规则">
                    <Input disabled multiline sx={{ width: "100%" }} value={localFileAttributes?.ignores} />
                </Attr>
            </>}
        </Grid>
    );
};
