import { Close } from '@mui/icons-material';
import ContentCopyIcon from '@mui/icons-material/ContentCopy';
import DeleteIcon from '@mui/icons-material/Delete';
import FileDownloadIcon from '@mui/icons-material/FileDownload';
import ModeEditIcon from '@mui/icons-material/ModeEdit';
import { Box, Dialog, DialogContent, DialogTitle, Link, Stack, Tooltip, Typography } from "@mui/material";
import IconButton from "@mui/material/IconButton";
import { download, openFile } from "api/fs";
import { getSysConfig } from "hox/sysConfig";
import useWindows, { closeWindow } from 'hox/windows';
import humanize from 'humanize';
import moment from 'moment';
import { useEffect, useState } from 'react';
import styles from './TextViewer.module.scss';

export default ({ id, props }) => {
    let { driver, filePath, dirItem } = props;
    dirItem = dirItem || {};
    const [loaded, setLoaded] = useState();
    console.log("fileViewer", id, props)
    const [windows, setWindows] = useWindows();
    let time = moment(dirItem.modifyTime / 1000 / 1000).format("YYYY年MM月DD日 HH:mm:ss");
    useEffect(() => {
        openFile(driver.id, filePath).then(setLoaded);
    }, []);
    return (
        <Dialog open={true} fullScreen={true} onClose={() => closeWindow(setWindows, id)}>
            <DialogTitle sx={{
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.secondary
            }}>
                文本编辑器
                <IconButton
                    aria-label="close"
                    onClick={() => closeWindow(setWindows, id)}
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
                <Stack
                    className={styles.fileHeaderViewer}
                    direction="row"
                    justifyContent="space-between"
                    alignItems="center"
                    spacing={0.5}
                >
                    {humanize.filesize(dirItem.size)} | {time}
                    <Stack
                        direction="row"
                        justifyContent="flex-end"
                        alignItems="flex-end"
                        spacing={0.5}
                    >
                        <Tooltip title="编辑">
                            <span><IconButton color="inherit" disabled={true}>
                                <ModeEditIcon fontSize="small" />
                            </IconButton></span>
                        </Tooltip>
                        <Tooltip title="复制文本内容">
                            <span><IconButton onClick={() => {
                                navigator.clipboard.writeText(loaded.content);
                            }} disabled={!loaded || loaded.tooLarge}>
                                <ContentCopyIcon fontSize="small" />
                            </IconButton></span>
                        </Tooltip>
                        <Tooltip title="下载">
                            <span><IconButton onClick={() => {
                                download(driver.id, filePath)
                            }}>
                                <FileDownloadIcon fontSize="small" />
                            </IconButton></span>
                        </Tooltip>
                        <Tooltip title="删除">
                            <span><IconButton disabled={true}>
                                <DeleteIcon fontSize="small" />
                            </IconButton></span>
                        </Tooltip>
                    </Stack>
                </Stack>
                <Box style={{ flex: "auto" }} className={styles.fileViewer}>
                    {!loaded ?
                        <>加载中...</> :
                        loaded.tooLarge ?
                            <>
                                文件大于{humanize.filesize(getSysConfig().sysConfig.maxContentSize)}，不支持在线查看，你可以选择
                                <Link underline="hover" onClick={() => {
                                    download(driver.id, filePath)
                                }}>下载该文件</Link>。
                            </> :
                            <Typography sx={{ wordBreak: "break-all" }}>
                                {loaded.content}
                            </Typography>
                    }
                </Box>
            </DialogContent>
        </Dialog>
    )
};
