import { getSysConfig } from "@kfs/common/hox/sysConfig";
import { analyzeMetadata } from "@kfs/mui/api/exif";
import { noteError, noteWarning } from "@kfs/mui/components/Notification";
import { TitleBar, Window, WorkingArea } from "@kfs/mui/components/Window/Window";
import { EventStreamContentType, fetchEventSource } from "@microsoft/fetch-event-source";
import { HourglassDisabled, HourglassTop, PlayArrow, Replay, Stop } from "@mui/icons-material";
import { Box, IconButton } from "@mui/material";
import moment from "moment";
import { useEffect, useState } from "react";

const StatusIdle = 0
const StatusFinished = 1
const StatusCanceled = 2
const StatusError = 3
const StatusWaitRunning = 4
const StatusWaitCanceled = 5
const StatusRunningCollect = 6
const StatusRunningAnalyze = 7
const StatusMsgs = {
    [undefined]: "加载中",
    [StatusIdle]: "空闲",
    [StatusFinished]: "已完成",
    [StatusCanceled]: "已取消",
    [StatusError]: "错误",
    [StatusWaitRunning]: "等待运行",
    [StatusWaitCanceled]: "正在取消",
    [StatusRunningCollect]: "正在收集",
    [StatusRunningAnalyze]: "正在解析",
};

export default function ({ id }) {
    const [taskInfo, setTaskInfo] = useState();
    const controller = new AbortController();
    useEffect(() => {
        setTaskInfo();
        fetchEventSource(`${getSysConfig().webServer}/api/v1/event/metadataAnalysisTask`, {
            signal: controller.signal,
            async onopen(response) {
                if (response.ok && response.headers.get('content-type').includes(EventStreamContentType)) {
                    return; // everything's good
                }
                console.error(response);
                noteError("event/metadataAnalysisTask.onopen: " + response.status);
            },
            onmessage(msg) {
                // if the server emits an error message, throw an exception
                // so it gets handled by the onerror callback below:
                if (msg.event === 'FatalError') {
                    console.error(msg);
                    noteError("event/metadataAnalysisTask.onmessage: " + msg);
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
                noteError("event/metadataAnalysisTask.onclose");
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
    }, []);
    return (
        <Window id={id}>
            <TitleBar id={id} title="元数据管理" />
            <WorkingArea>
                <Box sx={{ padding: "0.5em 1em" }}>
                    {(taskInfo?.status === undefined ||
                        taskInfo?.status === StatusIdle ||
                        taskInfo?.status === StatusFinished ||
                        taskInfo?.status === StatusCanceled ||
                        taskInfo?.status === StatusError) &&
                        <IconButton onClick={e => analyzeMetadata(true, true).catch(e => noteError(e.message))}>
                            <Replay />
                        </IconButton>
                    }
                    {(taskInfo?.status === undefined ||
                        taskInfo?.status === StatusIdle ||
                        taskInfo?.status === StatusFinished ||
                        taskInfo?.status === StatusCanceled ||
                        taskInfo?.status === StatusError) &&
                        <IconButton onClick={e => analyzeMetadata(true, false).catch(e => noteError(e.message))}>
                            <PlayArrow />
                        </IconButton>
                    }
                    {taskInfo?.status === StatusWaitRunning &&
                        <IconButton>
                            <HourglassTop />
                        </IconButton>
                    }
                    {taskInfo?.status === StatusWaitCanceled &&
                        <IconButton>
                            <HourglassDisabled />
                        </IconButton>
                    }
                    {(taskInfo?.status === StatusRunningCollect ||
                        taskInfo?.status === StatusRunningAnalyze) &&
                        <IconButton onClick={e => analyzeMetadata(false).catch(e => noteError(e.message))}>
                            <Stop />
                        </IconButton>
                    }
                    任务：解析图片与视频的元数据
                </Box>
                <Box sx={{ padding: "0.5em 1em" }}>
                    状态：{StatusMsgs[taskInfo?.status]}
                </Box>
                <Box sx={{ padding: "0.5em 1em" }}>
                    上次执行结束时间：{taskInfo?.lastDoneTime ? moment(taskInfo?.lastDoneTime / 1000 / 1000).format("YYYY年MM月DD日 HH:mm:ss") : "无记录"}
                </Box>
                {(taskInfo?.status === StatusRunningAnalyze) &&
                    <Box sx={{ padding: "0.5em 1em" }}>
                        当前进度：{taskInfo?.cnt}/{taskInfo?.total}
                    </Box>
                }
                {(taskInfo?.status === StatusFinished ||
                    taskInfo?.status === StatusCanceled ||
                    taskInfo?.status === StatusError) &&
                    <Box sx={{ padding: "0.5em 1em" }}>
                        上次执行结果：{taskInfo?.cnt}/{taskInfo?.total}
                    </Box>
                }
            </WorkingArea>
        </Window>
    )
}
