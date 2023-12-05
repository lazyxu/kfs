import { Close, ContentCopy, Save } from '@mui/icons-material';
import { default as Delete } from '@mui/icons-material/Delete';
import { default as FileDownload } from '@mui/icons-material/FileDownload';
import { Box, Dialog, DialogContent, Link, Stack } from "@mui/material";
import IconButton from "@mui/material/IconButton";
import { download, openFile } from "api/fs";
import { getSysConfig } from "hox/sysConfig";
import useWindows, { closeWindow } from 'hox/windows';
import humanize from 'humanize';
import moment from 'moment';
import { useEffect, useState } from 'react';

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
            <Stack direction="row" justifyContent="space-between"
                title={driver.name + ":/" + filePath.join("/")} sx={{
                    color: theme => theme.context.secondary,
                    backgroundColor: theme => theme.background.secondary,
                }}
            >
                <Box sx={{ height: "28px", lineHeight: "28px", paddingLeft: "1em" }}>
                    {dirItem.name} - 文本编辑器
                </Box>
                <IconButton aria-label="close" onClick={() => closeWindow(setWindows, id)}
                    sx={{
                        padding: "4px 12px", borderRadius: '0',
                        color: theme => theme.context.secondary,
                        ":hover": {
                            backgroundColor: "red",
                        }
                    }}
                >
                    <Close sx={{ width: "20px", height: "20px" }} />
                </IconButton>
            </Stack>
            <Box sx={{
                height: "28px", padding: "0",
                color: theme => theme.context.secondary,
                backgroundColor: theme => theme.background.secondary,
            }}>
                <Stack direction="row" alignItems="center" sx={{
                    height: "28px", padding: "0", paddingLeft: "1em",
                    backgroundColor: theme => theme.background.secondary,
                    color: theme => theme.context.secondary
                }}
                >
                    <IconButton title="保存" disabled
                        sx={{ height: "24px", width: "24px", color: theme => theme.context.secondary }}
                    >
                        <Save fontSize="small" sx={{ width: "20px", height: "20px" }} />
                    </IconButton>
                    <IconButton title="复制文本内容" onClick={() => {
                        navigator.clipboard.writeText(loaded.content);
                    }} disabled={!loaded || loaded.tooLarge}
                        sx={{ height: "24px", width: "24px", color: theme => theme.context.secondary }}
                    >
                        <ContentCopy fontSize="small" sx={{ width: "20px", height: "20px" }} />
                    </IconButton>
                    <IconButton title="下载" onClick={() => {
                        download(driver.id, filePath)
                    }}
                        sx={{ height: "24px", width: "24px", color: theme => theme.context.secondary }}
                    >
                        <FileDownload fontSize="small" sx={{ width: "20px", height: "20px" }} />
                    </IconButton>
                    <IconButton title="删除" disabled
                        sx={{ height: "24px", width: "24px", color: theme => theme.context.secondary }}
                    >
                        <Delete fontSize="small" sx={{ width: "20px", height: "20px" }} />
                    </IconButton>
                </Stack>
            </Box>
            <DialogContent sx={{
                padding: "0", paddingLeft: "5px",
                color: theme => theme.context.primary,
                backgroundColor: theme => theme.background.primary,
            }}>
                {!loaded ?
                    <>加载中...</> :
                    loaded.tooLarge ?
                        <>
                            文件大于{humanize.filesize(getSysConfig().sysConfig.maxContentSize)}，不支持在线查看，你可以选择
                            <Link underline="hover" onClick={() => {
                                download(driver.id, filePath)
                            }}>下载该文件</Link>。
                        </> :
                        <p style={{ wordBreak: "break-all", whiteSpace: "break-spaces", outline: "none" }} contentEditable>
                            {loaded.content}
                        </p>
                }
            </DialogContent>
            <Box sx={{
                flex: "0 0 auto", padding: "8px",
                color: theme => theme.context.secondary,
                backgroundColor: theme => theme.background.secondary,
            }}>
                <Stack direction="row" justifyContent="space-between"
                    title={driver.name + ":/" + filePath.join("/")} sx={{
                    }}
                >
                    <Box >
                        {humanize.filesize(dirItem.size)}
                    </Box>
                    <Box >
                        {time}
                    </Box>
                </Stack>
            </Box>
        </Dialog>
    )
};
