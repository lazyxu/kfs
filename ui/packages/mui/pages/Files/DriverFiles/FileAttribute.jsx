import { formatLinuxPermission, modeIsDir } from '@kfs/common/api/utils';
import { Close } from '@mui/icons-material';
import { Box, Dialog, DialogContent, DialogTitle, Grid, TextField, FormControlLabel, Checkbox, Button, DialogActions } from "@mui/material";
import IconButton from "@mui/material/IconButton";
import humanize from 'humanize';
import moment from "moment/moment";
import { useEffect, useState } from 'react';
import { createBaiduPhoto, updateBaiduPhoto } from '../../../api/driver';

function formatTime(t) {
    return moment(t / 1000 / 1000).format("YYYY年MM月DD日 HH:mm:ss");
}

function Attr({ k, children }) {
    return <>
        <Grid size={{ xs: 4 }} sx={{ overflowWrap: "anywhere" }}>
            <Box fontWeight="bold">{k}：</Box>
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

    const [cron, setCron] = useState(driverFile.cron || '');
    const [enabled, setEnabled] = useState(driverFile.enabled || false);
    const [syncAtStart, setSyncAtStart] = useState(driverFile.sync_at_start || false);
    const [bdussBfess, setBdussBfess] = useState(driverFile.bduss_bfess || '');
    const [stoken, setStoken] = useState(driverFile.stoken || '');
    const [attributes, setAttributes] = useState({});

    useEffect(() => {
        // getDriversDirCalculatedInfo(driver.id, filePath).then(setAttributes);
    }, []);

    const handleConfirm = () => {
        updateBaiduPhoto(driverFile.id, bdussBfess, stoken, cron, enabled, syncAtStart)
            .then(onClose);
    };

    return (
        <Dialog open={true} fullWidth={true} onClose={onClose} maxWidth="sm">
            <DialogTitle>
                属性
                <IconButton
                    aria-label="close"
                    onClick={() => onClose()}
                    sx={{
                        position: 'absolute',
                        right: 8,
                        top: 8,
                    }}
                >
                    <Close />
                </IconButton>
            </DialogTitle>
            <DialogContent>
                <Box sx={{ flexGrow: 1, mt: 2 }}>
                    <Grid container spacing={1.5} sx={{ alignItems: "center" }}>
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
                        {driverFile.sync_from_id > 0 && <>
                            <Grid size={{ xs: 12 }} sx={{ overflowWrap: "anywhere" }}>
                                <Box fontWeight="bold" fontSize="1.1rem">同步设置</Box>
                            </Grid>
                            <Attr k="cron表达式">
                                <TextField 
                                    size="small"
                                    value={cron}
                                    onChange={(e) => setCron(e.target.value)}
                                    placeholder="0 * * * *"
                                />
                            </Attr>
                            
                            <Attr k="开启同步">
                                <FormControlLabel
                                    control={
                                        <Checkbox 
                                            checked={enabled} 
                                            onChange={(e) => setEnabled(e.target.checked)} 
                                        />
                                    }
                                    label={enabled ? "已启用" : "已禁用"}
                                />
                            </Attr>
                            
                            <Attr k="启动时同步">
                                <FormControlLabel
                                    control={
                                        <Checkbox 
                                            checked={syncAtStart} 
                                            onChange={(e) => setSyncAtStart(e.target.checked)} 
                                        />
                                    }
                                    label={syncAtStart ? "是" : "否"}
                                />
                            </Attr>
                            
                            <Attr k="BDUSS_BFESS">
                                <TextField 
                                    size="small"
                                    value={bdussBfess}
                                    onChange={(e) => setBdussBfess(e.target.value)}
                                    placeholder="输入BDUSS_BFESS"
                                />
                            </Attr>
                            
                            <Attr k="STOKEN">
                                <TextField 
                                    size="small"
                                    value={stoken}
                                    onChange={(e) => setStoken(e.target.value)}
                                    placeholder="输入STOKEN"
                                />
                            </Attr>
                        </>}
                    </Grid>
                </Box>
            </DialogContent>
            {driverFile.sync_from_id > 0 && <DialogActions>
                <Button onClick={onClose}>取消</Button>
                <Button 
                    variant="contained" 
                    onClick={handleConfirm}
                    color="primary"
                >
                    确定
                </Button>
            </DialogActions>}
        </Dialog>
    );
};