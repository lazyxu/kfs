import { formatLinuxPermission, modeIsDir } from '@kfs/common/api/utils';
import { Close } from '@mui/icons-material';
import { Box, Dialog, DialogContent, DialogTitle, Grid } from "@mui/material";
import IconButton from "@mui/material/IconButton";
import humanize from 'humanize';
import moment from "moment/moment";
import { useEffect, useState } from 'react';

function formatTime(t) {
    return moment(t / 1000 / 1000).format("YYYY年MM月DD日 HH:mm:ss");
}

function Attr({ k, children }) {
    return <>
        <Grid size={{ xs: 4 }} sx={{ overflowWrap: "anywhere" }}>
            <Box>{k}：</Box>
        </Grid>
        <Grid size={{ xs: 8 }} sx={{ overflowWrap: "anywhere" }}>
            {children}
        </Grid>
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

export default ({ fileAttribute, onClose }) => {
    const { filePath, driverFile } = fileAttribute;
    const isDir = modeIsDir(driverFile.mode);
    const { name, mode } = driverFile;
    const [attributes, setAttributes] = useState({});
    useEffect(() => {
        // getDriversDirCalculatedInfo(driver.id, filePath).then(setAttributes);
    }, []);
    return (
        <Dialog open={true} fullWidth={true} onClose={onClose}>
            <DialogTitle>
                属性
                <IconButton
                    aria-label="close"
                    onClick={() => onClose()}
                    sx={{
                        position: 'absolute',
                        right: 8,
                        top: 8,
                        // color: (theme) => theme.palette.grey[500],
                    }}
                >
                    <Close />
                </IconButton>
            </DialogTitle>
            <DialogContent>
                <Box sx={{ flexGrow: 1 }}>
                <Grid container spacing={1.5} sx={{ alignItems: "center" }}>
                    {/* <Attr k="云盘ID">{driver.id}</Attr>
                    <Attr k="云盘名称">{driver.name}</Attr>
                    <Attr k="云盘描述">{driver.description}</Attr>
                    <Attr k="云盘类型">{getDriverType(driver)}</Attr> */}
                    <Attr k="类型">{isDir ? "文件夹" : "文件"}</Attr>
                    <Attr k="文件路径">{"/" + filePath.join("/")}</Attr>
                    {!isDir && <>
                        <Attr k="文件ID">{driverFile.file_id.Int64}</Attr>
                        <Attr k="文件大小">{humanize.filesize(driverFile.size)}</Attr>
                        <Attr k="md5">{driverFile.md5}</Attr>
                        <Attr k="sha256">{driverFile.sha256}</Attr>
                    </>}
                    <Attr k="文件权限">{formatLinuxPermission(driverFile.mode)}</Attr>
                    {/* {isDir && <>
                        <Attr k="目录下文件">共 {attributes.fileCount} 个，总大小 {humanize.filesize(attributes.fileSize)}</Attr>
                        <Attr k="目录下重复文件">共 {attributes.fileCount - attributes.distinctFileCount} 个，总大小  {humanize.filesize(attributes.fileSize - attributes.distinctFileSize)}</Attr>
                        <Attr k="子目录">共 {attributes.dirCount} 个</Attr>
                    </>
                    } */}
                    <Attr k="创建时间">{driverFile.create_time}</Attr>
                    <Attr k="属性修改时间">{driverFile.change_time}</Attr>
                    <Attr k="内容修改时间">{driverFile.modify_time}</Attr>
                    <Attr k="访问时间">{driverFile.access_time}</Attr>
                </Grid>
                </Box>
            </DialogContent>
        </Dialog>
    );
};
