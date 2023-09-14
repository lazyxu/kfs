import { EventStreamContentType, fetchEventSource } from "@microsoft/fetch-event-source";
import { Close, HourglassDisabled, HourglassTop, PlayArrow, RestartAlt, SettingsApplications, Stop } from "@mui/icons-material";
import { Box, Button, IconButton, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Typography } from "@mui/material";
import { deleteBackupTask, startBackupTask } from "api/web/backup";
import { noteError, noteSuccess } from "components/Notification/Notification";
import { getSysConfig } from "hox/sysConfig";
import moment from "moment";
import { useEffect, useState } from "react";

const StatusIdle = 0;
const StatusWaitRunning = 1;
const StatusRunning = 2;
const StatusFinished = 3;
const StatusCanceled = 4;
const StatusError = 5;
const StatusWaitCanceled = 6;
const StatusMsgs = {
    undefined: "空闲",
    0: "空闲",
    1: "等待运行",
    2: "正在运行",
    3: "已完成",
    4: "已取消",
    5: "错误",
    6: "正在取消",
};

export default function ({ setTaskDetail }) {
    const sysConfig = getSysConfig().sysConfig;
    const [taskInfos, setTaskInfos] = useState({ list: [], runningTasks: {} });
    const controller = new AbortController();
    useEffect(() => {
        fetchEventSource(`http://127.0.0.1:${getSysConfig().sysConfig.port}/api/v1/event/backupTask`, {
            signal: controller.signal,
            async onopen(response) {
                if (response.ok && response.headers.get('content-type').includes(EventStreamContentType)) {
                    return; // everything's good
                }
                console.error(response);
                noteError("event/backupTask.onopen: " + response.status);
            },
            onmessage(msg) {
                // if the server emits an error message, throw an exception
                // so it gets handled by the onerror callback below:
                if (msg.event === 'FatalError') {
                    console.error(msg);
                    noteError("event/backupTask.onmessage: " + msg);
                    return;
                }
                let info = JSON.parse(msg.data);
                console.log(info);
                setTaskInfos(info);
            },
            onclose() {
                // if the server closes the connection unexpectedly, retry:
                noteError("event/backupTask.onclose");
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
        <Box sx={{
            width: "100%", flex: "1",
            display: 'flex', flexDirection: 'column', minHeight: '0'
        }}>
            <Typography variant="h6" noWrap component="div" sx={{ marginLeft: 2 }}>
                任务列表
            </Typography>
            <TableContainer sx={{
                width: "100%", flex: "1", overflowY: 'auto', alignContent: "flex-start",
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.primary
            }}
            >
                <Table aria-label="simple table">
                    <TableHead>
                        <TableRow>
                            <TableCell align="left">任务名称</TableCell>
                            <TableCell align="left">任务状态</TableCell>
                            <TableCell align="left">源目录</TableCell>
                            <TableCell align="left">备份方式</TableCell>
                            <TableCell align="left">云盘名称</TableCell>
                            <TableCell align="left">目标目录</TableCell>
                            <TableCell align="left">上次执行结束时间</TableCell>
                            <TableCell align="left">操作</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {taskInfos.list.map((task) => (
                            <TableRow
                                key={task.name}
                                sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                            >
                                <TableCell component="th" scope="row">
                                    <Button onClick={() => setTaskDetail(task.name)}>{task.name}</Button>
                                </TableCell>
                                <TableCell align="left">{StatusMsgs[taskInfos.runningTasks[task.name]?.status]}</TableCell>
                                <TableCell align="left">{task.srcPath}</TableCell>
                                <TableCell align="left"><span title="单向上传：源目录的所有内容将会被上传更新到云盘，删除操作将不会造成云盘内容的对应删除。">单向上传</span></TableCell>
                                <TableCell align="left">{task.driverName}</TableCell>
                                <TableCell align="left">{task.dstPath}</TableCell>
                                <TableCell align="left">{taskInfos.runningTasks[task.name]?.lastDoneTime ? moment(taskInfos.runningTasks[task.name]?.lastDoneTime / 1000 / 1000).format("YYYY年MM月DD日 HH:mm:ss") : "无记录"}</TableCell>
                                <TableCell align="left">
                                    <IconButton disabled>
                                        <RestartAlt />
                                    </IconButton>
                                    {(taskInfos.runningTasks[task.name]?.status === undefined ||
                                        taskInfos.runningTasks[task.name]?.status === StatusIdle ||
                                        taskInfos.runningTasks[task.name]?.status === StatusFinished ||
                                        taskInfos.runningTasks[task.name]?.status === StatusCanceled ||
                                        taskInfos.runningTasks[task.name]?.status === StatusError) &&
                                        <IconButton onClick={e => startBackupTask(task.name, sysConfig.socketServer, true)
                                            .then(() => noteSuccess("运行备份任务：" + task.name))
                                            .catch(e => noteError(e.message))
                                        }>
                                            <PlayArrow />
                                        </IconButton>
                                    }
                                    {taskInfos.runningTasks[task.name]?.status === StatusWaitRunning &&
                                        <IconButton>
                                            <HourglassTop />
                                        </IconButton>
                                    }
                                    {taskInfos.runningTasks[task.name]?.status === StatusWaitCanceled &&
                                        <IconButton>
                                            <HourglassDisabled />
                                        </IconButton>
                                    }
                                    {taskInfos.runningTasks[task.name]?.status === StatusRunning &&
                                        <IconButton onClick={e => startBackupTask(task.name, sysConfig.socketServer, false)
                                            .then(() => noteSuccess("停止备份任务：" + task.name))
                                            .catch(e => noteError(e.message))
                                        }>
                                            <Stop />
                                        </IconButton>
                                    }

                                    <IconButton disabled>
                                        <SettingsApplications />
                                    </IconButton>
                                    <IconButton disabled={taskInfos.runningTasks[task.name]?.status === StatusRunning ||
                                        taskInfos.runningTasks[task.name]?.status === StatusWaitRunning ||
                                        taskInfos.runningTasks[task.name]?.status === StatusWaitCanceled
                                    }
                                        onClick={e => deleteBackupTask(task.name)
                                            .then(() => noteSuccess("删除备份任务：" + task.name))
                                            .catch(e => noteError(e.message))
                                        }
                                    >
                                        <Close />
                                    </IconButton>
                                </TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </TableContainer>
        </Box>
    );
}
