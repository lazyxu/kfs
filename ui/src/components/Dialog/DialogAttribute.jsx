import React from 'react';
import { Dialog, DialogContent, DialogTitle, Grid } from "@mui/material";
import useDialog from "hox/dialog";
import useResourceManager from "hox/resourceManager";
import IconButton from "@mui/material/IconButton";
import CloseIcon from '@mui/icons-material/Close';
import moment from "moment/moment";
import humanize from "humanize";
import { getPerm, modeIsDir } from "../../api/utils/api";
import styles from "./DialogAttribute.module.scss"

function formatTime(t) {
    return moment(t / 1000 / 1000).format("YYYY年MM月DD日 HH:mm:ss");
}

function Attr({k, children}) {
    return <>
        <Grid xs={4} item={true} className={styles.attrubute} sx={{color: (theme) => theme.palette.grey[500]}}>{k}：</Grid>
        <Grid xs={8} item={true} className={styles.attrubute} sx={{color: (theme) => theme.palette.grey[500]}}>{children}</Grid>
    </>
}

export default () => {
    const [dialog, setDialog] = useDialog();
    const [resourceManager, setResourceManager] = useResourceManager();
    let { filePath, branchName } = resourceManager;
    filePath = filePath.slice()
    filePath.push(dialog.dirItem.Name);
    const isDir = modeIsDir(dialog.dirItem.Mode)
    return (
        <Dialog sx={{ m: 0, p: 2 }} open={true} fullWidth={true} onClose={() => {
            setDialog(null)
        }}>
            <DialogTitle>
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
                    <CloseIcon />
                </IconButton>
            </DialogTitle>
            <DialogContent>
                <Grid container spacing={1.5}>
                    <Attr k="名称">{dialog.dirItem.Name}</Attr>
                    <Attr k="分支">{branchName}</Attr>
                    <Attr k="路径">{"/" + filePath.join("/")}</Attr>
                    <Attr k="哈希值">{dialog.dirItem.Hash}</Attr>
                    <Attr k="类型">{isDir ? "文件夹" : "文件"}</Attr>
                    <Attr k="文件大小">{humanize.filesize(dialog.dirItem.Size)}</Attr>
                    <Attr k="文件权限">{getPerm(dialog.dirItem.Mode).toString(8)}</Attr>
                    {isDir && <>
                        <Attr k="文件数量">{dialog.dirItem.Count}</Attr>
                        <Attr k="文件总数量">{dialog.dirItem.TotalCount}</Attr>
                    </>
                    }
                    <Attr k="创建时间">{formatTime(dialog.dirItem.CreateTime)}</Attr>
                    <Attr k="属性修改时间">{formatTime(dialog.dirItem.ChangeTime)}</Attr>
                    <Attr k="内容修改时间">{formatTime(dialog.dirItem.ModifyTime)}</Attr>
                    <Attr k="访问时间">{formatTime(dialog.dirItem.AccessTime)}</Attr>
                </Grid>
            </DialogContent>
        </Dialog>
    );
};
