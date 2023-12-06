import { Close } from '@mui/icons-material';
import { Box, Dialog, DialogContent, DialogTitle, Grid } from "@mui/material";
import IconButton from "@mui/material/IconButton";
import { getDriversDirCalculatedInfo } from 'api/driver';
import { getPerm, modeIsDir } from 'api/utils/api';
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

export default ({ fileAttribute, setFileAttribute }) => {
    const { driver, filePath, dirItem } = fileAttribute;
    const isDir = modeIsDir(dirItem.mode);
    const [attributes, setAttributes] = useState({});
    const { name, mode } = dirItem;
    const curFilePath = filePath.concat(name);
    useEffect(() => {
        getDriversDirCalculatedInfo(driver.id, curFilePath).then(setAttributes);
    }, []);
    return (
        <Dialog open={true} fullWidth={true} onClose={() => setFileAttribute(null)}>
            <DialogTitle sx={{
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.secondary
            }}>
                云盘属性
                <IconButton
                    aria-label="close"
                    onClick={() => setFileAttribute(null)}
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
                    <Attr k="云盘ID">{driver.id}</Attr>
                    <Attr k="云盘名称">{driver.name}</Attr>
                    <Attr k="云盘描述">{driver.description}</Attr>
                    <Attr k="云盘类型">{getDriverType(driver)}</Attr>
                    <Attr k="文件路径">{"/" + curFilePath.join("/")}</Attr>
                    <Attr k="哈希值">{dirItem.hash}</Attr>
                    <Attr k="类型">{isDir ? "文件夹" : "文件"}</Attr>
                    {!isDir && <Attr k="文件大小">{humanize.filesize(dirItem.size)}</Attr>}
                    <Attr k="文件权限">{getPerm(dirItem.mode).toString(8)}</Attr>
                    {isDir && <>
                        <Attr k="目录下总大小">{humanize.filesize(attributes.fileSize)}</Attr>
                        <Attr k="目录下文件总数量">{attributes.fileCount}</Attr>
                        <Attr k="目录下目录总数量">{attributes.dirCount}</Attr>
                    </>
                    }
                    <Attr k="创建时间">{formatTime(dirItem.createTime)}</Attr>
                    <Attr k="属性修改时间">{formatTime(dirItem.changeTime)}</Attr>
                    <Attr k="内容修改时间">{formatTime(dirItem.modifyTime)}</Attr>
                    <Attr k="访问时间">{formatTime(dirItem.accessTime)}</Attr>
                </Grid>
            </DialogContent>
        </Dialog>
    );
};
