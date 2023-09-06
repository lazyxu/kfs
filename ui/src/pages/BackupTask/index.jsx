import { Box, Button, Paper, Stack, Tab, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Tabs } from "@mui/material";
import { useEffect, useState } from "react";
import FastBackup from "./FastBackup";
import { listBackupTask } from "api/web/backup";
import NewTask from "./NewTask";

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
        listBackupTask().then(setBackupTasks);
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
                                <TableCell align="right">未知状态</TableCell>
                                <TableCell align="right">{task.srcPath}</TableCell>
                                <TableCell align="left">未知方式</TableCell>
                                <TableCell align="right">{task.driverName}</TableCell>
                                <TableCell align="right">{task.dstPath}</TableCell>
                                <TableCell align="left">未知时间</TableCell>
                                <TableCell align="left">
                                    <Button>重新运行</Button>
                                    <Button>停止</Button>
                                    <Button>设置</Button>
                                    <Button>删除</Button>
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
