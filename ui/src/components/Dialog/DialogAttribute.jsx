import React from 'react';
import {Dialog, DialogContent, DialogTitle, Grid} from "@mui/material";
import useDialog from "hox/dialog";
import useResourceManager from "hox/resourceManager";
import IconButton from "@mui/material/IconButton";
import CloseIcon from '@mui/icons-material/Close';
import moment from "moment/moment";
import humanize from "humanize";
import {getPerm, modeIsDir} from "../../api/utils/api";

function formatTime(t) {
    return moment(t / 1000 / 1000).format("YYYY年MM月DD日 HH:mm:ss");
}

function Attr({k, children}) {
    return <>
        <Grid xs={4} item sx={{overflowWrap: "anywhere"}}>{k}：</Grid>
        <Grid xs={8} item sx={{overflowWrap: "anywhere"}}>{children}</Grid>
    </>
}

export default () => {
    const [dialog, setDialog] = useDialog();
    const [resourceManager, setResourceManager] = useResourceManager();
    let {filePath, driverName} = resourceManager;
    filePath = filePath.slice()
    filePath.push(dialog.dirItem.name);
    const isDir = modeIsDir(dialog.dirItem.mode)
    return (
        <Dialog open={true} fullWidth={true} onClose={() => {
            setDialog(null)
        }}>
            <DialogTitle sx={{
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.secondary
            }}>
                {dialog.title}
                <IconButton
                    aria-label="close"
                    onClick={() => {
                        setDialog(null);
                    }}
                    sx={{
                        position: 'absolute',
                        right: 8,
                        top: 8,
                        color: (theme) => theme.palette.grey[500],
                    }}
                >
                    <CloseIcon/>
                </IconButton>
            </DialogTitle>
            <DialogContent sx={{
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.primary
            }}>
                <Grid container spacing={1.5}>
                    <Attr k="名称">{dialog.dirItem.name}</Attr>
                    <Attr k="云盘">{driverName}</Attr>
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
