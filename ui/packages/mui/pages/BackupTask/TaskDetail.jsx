import { getSysConfig } from "@kfs/common/hox/sysConfig";
import { noteError, noteWarning } from "@kfs/mui/components/Notification";
import { EventStreamContentType, fetchEventSource } from "@microsoft/fetch-event-source";
import { Close } from "@mui/icons-material";
import { Alert, Box, IconButton, Typography } from "@mui/material";
import humanize from "humanize";
import humanizeDuration from "humanize-duration";
import { useEffect, useState } from "react";
import LinearProgressWithLabel from "./LinearProgressWithLabel";

export default function ({ taskDetail, setTaskDetail }) {
    const sysConfig = getSysConfig();
    const [taskInfo, setTaskInfo] = useState();
    const controller = new AbortController();
    useEffect(() => {
        setTaskInfo();
        fetchEventSource(`http://127.0.0.1:${getSysConfig().port}/api/v1/event/backupTask/${taskDetail}`, {
            signal: controller.signal,
            async onopen(response) {
                if (response.ok && response.headers.get('content-type').includes(EventStreamContentType)) {
                    return; // everything's good
                }
                console.error(response);
                noteError("event/backupTask/detail.onopen: " + response.status);
            },
            onmessage(msg) {
                // if the server emits an error message, throw an exception
                // so it gets handled by the onerror callback below:
                if (msg.event === 'FatalError') {
                    console.error(msg);
                    noteError("event/backupTask/detail.onmessage: " + msg);
                    return;
                }
                let info = JSON.parse(msg.data);
                console.log(info);
                if (info?.errMsg) {
                    noteError(info?.errMsg);
                }
                if (info?.data?.errMsg) {
                    noteWarning(info?.data?.errMsg);
                }
                setTaskInfo(info);
            },
            onclose() {
                // if the server closes the connection unexpectedly, retry:
                noteError("event/backupTask/detail.onclose");
            },
            onerror(err) {
                console.error(err);
                // noteError("event/backupTask.onerror: " + err.message);
                // if (err instanceof FatalError) {
                //     throw err; // rethrow to stop the operation
                // } else {
                //     // do nothing to automatically retry. You can also
                //     // return a specific retry interval here.
                // }
            }
        });
        return () => {
            controller.abort();
        }
    }, [taskDetail]);
    return (
        <Box sx={{
            width: "100%", flex: "1",
            display: 'flex', flexDirection: 'column', minHeight: '0'
        }}>
            <Typography variant="h6" noWrap component="div" sx={{ marginLeft: 2 }}>
                任务日志： {taskDetail}
                <IconButton
                    aria-label="close"
                    onClick={() => setTaskDetail()}
                    sx={{
                        color: (theme) => theme.palette.grey[500],
                    }}
                >
                    <Close />
                </IconButton>
            </Typography>
            {taskInfo?.data?.cost ?
                <Alert variant="outlined" classes={{ message: "width100" }} severity={taskInfo.finished ? "success" : "info"}>
                    {taskInfo.data.pushedAllToStack &&
                        <LinearProgressWithLabel variant="determinate" value={taskInfo.data.size / taskInfo.data.totalSize * 100} />
                    }
                    <Typography>进度：{humanize.filesize(taskInfo.data.size)}/{humanize.filesize(taskInfo.data.totalSize)}</Typography>
                    <Typography>耗时：{humanizeDuration(Math.floor(taskInfo.data.cost / 1000) * 1000)}</Typography>
                    <Typography>平均上传速度：{humanize.filesize(taskInfo.data.size * 1000 / taskInfo.data.cost)}/s</Typography>
                    <Typography>预计剩余时间：{humanizeDuration(Math.floor(taskInfo.data.totalSize === taskInfo.data.size ? 0 : taskInfo.data.cost / taskInfo.data.size * (taskInfo.data.totalSize - taskInfo.data.size) / 1000) * 1000)}</Typography>
                    <Typography>文件：{taskInfo.data.fileCount}/{taskInfo.data.totalFileCount}</Typography>
                    <Typography>目录：{taskInfo.data.dirCount}/{taskInfo.data.totalDirCount}</Typography>
                    <Typography>上传列表：</Typography>
                    {taskInfo.data.processes.map((process, i) =>
                        process.filePath ? <Typography key={i}>{i + 1}： {StatusList[process.status]} {humanize.filesize(process.size)} {process.filePath}</Typography>
                            : <Typography key={i}>{i + 1}：空闲</Typography>
                    )}
                </Alert> : <Box />}
        </Box>
    );
}

const StatusList = ["正在上传", "已经存在", "上传成功"];
