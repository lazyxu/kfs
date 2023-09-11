import { Box, Button, Stack } from "@mui/material";
import { useState } from "react";
import NewTask from "./NewTask";
import TaskDetail from "./TaskDetail";
import TaskList from "./TaskList";

export default function () {
    const [open, setOpen] = useState(false);
    const [taskDetail, setTaskDetail] = useState();
    return (
        <Box style={{ display: 'flex', flex: "1", flexDirection: 'column', minHeight: '0' }}>
            <TaskList setTaskDetail={setTaskDetail} />
            {taskDetail && <TaskDetail taskDetail={taskDetail} setTaskDetail={setTaskDetail} />}
            <Stack direction="row" justifyContent="flex-end" alignItems="center" spacing={1} >
                <Button variant="outlined" onClick={() => setOpen(true)}>创建新的备份任务</Button>
            </Stack>
            <NewTask open={open} setOpen={setOpen} />
        </Box>
    );
}
