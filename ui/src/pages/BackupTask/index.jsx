import { Box, Button, IconButton, Paper, Stack, Tab, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Tabs } from "@mui/material";
import { useEffect, useState } from "react";
import FastBackup from "./FastBackup";
import { deleteBackupTask, listBackupTask } from "api/web/backup";
import NewTask from "./NewTask";
import { noteError, noteInfo, noteSuccess } from "components/Notification/Notification";
import { EventStreamContentType, fetchEventSource } from "@microsoft/fetch-event-source";
import { Close, Info, PlayArrow, RestartAlt, SettingsApplications, Start, Stop } from "@mui/icons-material";

function createData(name, calories, fat, carbs, protein) {
    return { name, calories, fat, carbs, protein };
}

const rows = [
    createData('Frozen yoghurt', 159, 6.0, 24, 4.0),
    createData('Ice cream sandwich', 237, 9.0, 37, 4.3),
    createData('Eclair', 262, 16.0, 24, 6.0),
    createData('Cupcake', 305, 3.7, 67, 4.3),
    createData('Gingerbread', 356, 16.0, 49, 3.9),
];

export default function ({ show }) {
    const [open, setOpen] = useState(false);
    const [backupTasks, setBackupTasks] = useState([]);
    useEffect(() => {
        fetchEventSource('http://127.0.0.1:11235/api/v1/event/backupTask', {
            async onopen(response) {
                if (response.ok && response.headers.get('content-type') === EventStreamContentType) {
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
                setBackupTasks(JSON.parse(msg.data));
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
    }, []);
    return (
        <Box style={{ display: show ? 'flex' : "none", flex: "1", flexDirection: 'column', minHeight: '0' }}>
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
                            <TableCell align="left">最后执行时间</TableCell>
                            <TableCell align="left">操作</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {backupTasks.map((task) => (
                            <TableRow
                                key={task.name}
                                sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                            >
                                <TableCell component="th" scope="row">
                                    {task.name}
                                </TableCell>
                                <TableCell align="left">未知状态</TableCell>
                                <TableCell align="left">{task.srcPath}</TableCell>
                                <TableCell align="left">未知方式</TableCell>
                                <TableCell align="left">{task.driverName}</TableCell>
                                <TableCell align="left">{task.dstPath}</TableCell>
                                <TableCell align="left">未知时间</TableCell>
                                <TableCell align="left">
                                    <IconButton disabled>
                                        <RestartAlt />
                                    </IconButton>
                                    <IconButton disabled>
                                        <PlayArrow />
                                    </IconButton>
                                    <IconButton disabled>
                                        <Stop />
                                    </IconButton>
                                    <IconButton disabled>
                                        <SettingsApplications />
                                    </IconButton>
                                    <IconButton
                                        onClick={e => {
                                            deleteBackupTask(task.name)
                                                .then(() => noteSuccess("删除备份任务：" + task.name))
                                                .catch(e => noteError(e.message))
                                        }} >
                                        <Close />
                                    </IconButton>
                                </TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </TableContainer>
            <Stack className='filePath'
                direction="row"
                justifyContent="flex-end"
                alignItems="center"
                spacing={1}
            >
                <Button variant="outlined" onClick={() => setOpen(true)}>创建新的备份任务</Button>
            </Stack>
            <NewTask open={open} setOpen={setOpen} />
        </Box>
    );
}
