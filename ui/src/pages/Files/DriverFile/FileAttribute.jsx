import { Close } from '@mui/icons-material';
import { Box, Dialog, DialogContent, DialogTitle, Grid } from "@mui/material";
import IconButton from "@mui/material/IconButton";
import moment from "moment/moment";
import { useState } from 'react';

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
    const [attributeType, setAttributeType] = useState(0);
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
                    <Attr k="云盘ID">{driver.id}</Attr>
                    <Attr k="云盘名称">{driver.name}</Attr>
                    <Attr k="云盘描述">{driver.description}</Attr>
                    <Attr k="云盘类型">{getDriverType(driver)}</Attr>

                    <Attr k="名称">{dialog.dirItem.name}</Attr>
                    <Attr k="路径">{"/" + filePath.join("/")}</Attr>
                    <Attr k="哈希值">{dialog.dirItem.hash}</Attr>
                    <Attr k="类型">{isDir ? "文件夹" : "文件"}</Attr>
                    <Attr k="文件大小">{humanize.filesize(dialog.dirItem.size)}</Attr>
                    <Attr k="文件权限">{getPerm(dialog.dirItem.mode).toString(8)}</Attr>
                    {isDir && <>
                        <Attr k="文件数量">{dialog.dirItem.count}</Attr>
                        <Attr k="文件总数量">{dialog.dirItem.totalCount}</Attr>
                    </>
                    }
                    <Attr k="创建时间">{formatTime(dialog.dirItem.createTime)}</Attr>
                    <Attr k="属性修改时间">{formatTime(dialog.dirItem.changeTime)}</Attr>
                    <Attr k="内容修改时间">{formatTime(dialog.dirItem.modifyTime)}</Attr>
                    <Attr k="访问时间">{formatTime(dialog.dirItem.accessTime)}</Attr>
                </Grid>
            </DialogContent>
        </Dialog>
    );
};
