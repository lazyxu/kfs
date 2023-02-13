import React, {useState} from 'react';
import {Dialog, DialogContent, DialogContentText, DialogTitle, Grid} from "@mui/material";
import useDialog from "hox/dialog";
import useResourceManager from "hox/resourceManager";
import useSysConfig from "hox/sysConfig";
import IconButton from "@mui/material/IconButton";
import CloseIcon from '@mui/icons-material/Close';
import moment from "moment/moment";
import humanize from "humanize";
import {getPerm, modeIsDir} from "../../api/utils/api";

export default () => {
    const [dialog, setDialog] = useDialog();
    let [name, setName] = useState("");
    const [resourceManager, setResourceManager] = useResourceManager();
    let {filePath, branchName} = resourceManager;
    const {sysConfig} = useSysConfig();
    const isDir = modeIsDir(dialog.dirItem.Mode)
    return (
        <Dialog sx={{m: 0, p: 2}} open={true} fullWidth={true} onClose={() => {
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
                    <CloseIcon/>
                </IconButton>
            </DialogTitle>
            <DialogContent>
                <DialogContentText>
                    <Grid container>
                        <Grid xs={4}>名称：</Grid>
                        <Grid xs={8}>{dialog.dirItem.Name}</Grid>
                        <Grid xs={4}>哈希值：</Grid>
                        <Grid xs={8}>{dialog.dirItem.Hash}</Grid>
                        <Grid xs={4}>类型：</Grid>
                        <Grid xs={8}>{isDir ? "文件夹" : "文件"}</Grid>
                        <Grid xs={4}>文件大小：</Grid>
                        <Grid xs={8}>{humanize.filesize(dialog.dirItem.Size)}</Grid>
                        <Grid xs={4}>文件权限：</Grid>
                        <Grid xs={8}>{getPerm(dialog.dirItem.Mode).toString(8)}</Grid>
                        {isDir && <>
                            <Grid xs={4}>文件数量：</Grid>
                            <Grid xs={8}>{dialog.dirItem.Count}</Grid>
                            <Grid xs={4}>文件总数量：</Grid>
                            <Grid xs={8}>{dialog.dirItem.TotalCount}</Grid>
                        </>
                        }
                        <Grid xs={4}>创建时间：</Grid>
                        <Grid
                            xs={8}>{moment(dialog.dirItem.CreateTime / 1000 / 1000).format("YYYY年MM月DD日 HH:mm:ss")}</Grid>
                        <Grid xs={4}>属性修改时间：</Grid>
                        <Grid
                            xs={8}>{moment(dialog.dirItem.ChangeTime / 1000 / 1000).format("YYYY年MM月DD日 HH:mm:ss")}</Grid>
                        <Grid xs={4}>内容修改时间：</Grid>
                        <Grid
                            xs={8}>{moment(dialog.dirItem.ModifyTime / 1000 / 1000).format("YYYY年MM月DD日 HH:mm:ss")}</Grid>
                        <Grid xs={4}>访问时间：</Grid>
                        <Grid
                            xs={8}>{moment(dialog.dirItem.AccessTime / 1000 / 1000).format("YYYY年MM月DD日 HH:mm:ss")}</Grid>
                    </Grid>
                </DialogContentText>
            </DialogContent>
        </Dialog>
    );
};
