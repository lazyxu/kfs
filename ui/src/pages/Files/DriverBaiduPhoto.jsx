import { EventStreamContentType, fetchEventSource } from '@microsoft/fetch-event-source';
import DeleteIcon from '@mui/icons-material/Delete';
import DriveFileRenameOutlineIcon from '@mui/icons-material/DriveFileRenameOutline';
import { Box, Button, Card, CardActions, CardContent, Link, Stack } from "@mui/material";
import { list } from "api/fs";
import SvgIcon from "components/Icon/SvgIcon";
import { noteError } from 'components/Notification/Notification';
import useContextMenu from "hox/contextMenu";
import useResourceManager from 'hox/resourceManager';
import { getSysConfig } from 'hox/sysConfig';
import { useEffect, useState } from 'react';
import { deleteDriver } from "../../api/driver";

export default ({ driver, setDriverAttribute }) => {
    const [resourceManager, setResourceManager] = useResourceManager();
    const [contextMenu, setContextMenu] = useContextMenu();

    const [taskInfo, setTaskInfo] = useState();
    const controller = new AbortController();
    useEffect(() => {
        fetchEventSource(`${getSysConfig().sysConfig.webServer}/api/v1/event/baiduPhotoTask/${driver.name}`, {
            signal: controller.signal,
            async onopen(response) {
                if (response.ok && response.headers.get('content-type').includes(EventStreamContentType)) {
                    return; // everything's good
                }
                console.error(response);
                noteError("event/baiduPhotoTask.onopen: " + response.status);
            },
            onmessage(msg) {
                // if the server emits an error message, throw an exception
                // so it gets handled by the onerror callback below:
                if (msg.event === 'FatalError') {
                    console.error(msg);
                    noteError("event/baiduPhotoTask.onmessage: " + msg);
                    return;
                }
                let info = JSON.parse(msg.data);
                console.log(info);
                setTaskInfo(info);
            },
            onclose() {
                // if the server closes the connection unexpectedly, retry:
                noteError("event/baiduPhotoTask.onclose");
            },
            onerror(err) {
                console.error(err);
                // noteError("event/baiduPhotoTask.onerror: " + err.message);
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
        <Card sx={{ minWidth: 275 }} variant="outlined">
            <CardContent>
                <Link underline="hover" onClick={() => list(setResourceManager, driver.name, [])}>
                    <Stack
                        direction="row"
                        justifyContent="space-between"
                        alignItems="center"
                        spacing={2}
                    >
                        <SvgIcon icon="wangpan" fontSize="inherit" />
                        <Box sx={{ flex: 1 }}>{driver.name}</Box>
                    </Stack>
                </Link>
                <Box variant="body">
                    {/* <Box>文件总数：{driver.count}</Box> */}
                    {/* <Box>总大小：{humanize.filesize(driver.size)}</Box> */}
                    {/* <Typography>可修改该云盘的设备：any</Typography> */}
                    {/* <Typography>可读取该云盘的设备：any</Typography> */}
                </Box>
                <Box color="text.secondary">
                    {driver.description}
                </Box>
                <Box color="text.secondary">
                    <IconButton onClick={e => analyzeMetadata(true).catch(e => noteError(e.message))}>
                        <PlayArrow />
                    </IconButton>
                    [一刻相册] {taskInfo ? String(taskInfo.cnt) + "/" + taskInfo.total : "loading..."}
                </Box>
            </CardContent>
            <CardActions>
                <Button size="small" color="error" startIcon={<DeleteIcon />} variant="outlined"
                    onClick={() => deleteDriver(setResourceManager, driver.name)}>删除</Button>
                <Button size="small" startIcon={<DriveFileRenameOutlineIcon />} variant="outlined" disabled>重命名</Button>
                <Button size="small" variant="outlined" disabled>重置</Button>
                <Button size="small" variant="outlined"
                    onClick={() => setDriverAttribute(driver)} >属性</Button>
            </CardActions>
        </Card>
    )
};
