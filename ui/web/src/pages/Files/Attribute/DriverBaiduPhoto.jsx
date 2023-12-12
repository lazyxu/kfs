import { getDriverSync, updateDriverSync } from '@/api/driver';
import { startBaiduPhotoTask } from '@/api/exif';
import { noteError } from '@/components/Notification/Notification';
import { getSysConfig } from '@/hox/sysConfig';
import { EventStreamContentType, fetchEventSource } from '@microsoft/fetch-event-source';
import { HourglassDisabled, HourglassTop, PlayArrow, Stop } from '@mui/icons-material';
import { Box, Grid, IconButton, MenuItem, Select, Switch } from "@mui/material";
import { useEffect, useState } from 'react';

const StatusIdle = 0
const StatusFinished = 1
const StatusCanceled = 2
const StatusError = 3
const StatusWaitRunning = 4
const StatusWaitCanceled = 5
const StatusRunning = 6

function Attr({ k, children }) {
    return <>
        <Grid xs={4} item sx={{ overflowWrap: "anywhere" }}><Box>{k}：</Box></Grid>
        <Grid xs={8} item sx={{ overflowWrap: "anywhere" }}>{children}</Grid>
    </>
}

export default ({ driver }) => {
    const [taskInfo, setTaskInfo] = useState();
    const [syncAttributes, setSyncAttributes] = useState();
    useEffect(() => {
        getDriverSync(driver.id).then(n => setSyncAttributes(n));
    }, []);
    const controller = new AbortController();
    useEffect(() => {
        fetchEventSource(`${getSysConfig().sysConfig.webServer}/api/v1/event/baiduPhotoTask/${driver.id}`, {
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
                if (info.errMsg) {
                    noteError(info.errMsg);
                }
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
        <Grid container spacing={1.5} sx={{ alignItems: "center" }}>
            <Attr k="定时同步">
        {syncAttributes ? <>
            <Switch checked={syncAttributes.sync} onChange={e => updateDriverSync(e.target.checked, syncAttributes.h, syncAttributes.m)} />
            <Select variant="standard" size="small" sx={{ marginLeft: "1em" }} value={syncAttributes.h} onChange={e => updateDriverSync(syncAttributes.sync, e.target.value, syncAttributes.m)}>
                {[...Array(24).keys()].map(value =>
                    <MenuItem key={value} value={value}>{value.toString().padStart(2, 0)}</MenuItem>
                )}
            </Select>时
            <Select variant="standard" size="small" sx={{ marginLeft: "1em" }} value={syncAttributes.m} onChange={e => updateDriverSync(syncAttributes.sync, syncAttributes.h, e.target.value)}>
                {[...Array(60).keys()].map(value =>
                    <MenuItem key={value} value={value}>{value.toString().padStart(2, 0)}</MenuItem>
                )}
            </Select>分
        </> : <>配置加载中...</>}
    </Attr>
    <Attr k="同步">
            {(taskInfo?.status === undefined ||
                taskInfo?.status === StatusIdle ||
                taskInfo?.status === StatusFinished ||
                taskInfo?.status === StatusCanceled ||
                taskInfo?.status === StatusError) &&
                <IconButton onClick={e => startBaiduPhotoTask(true, driver.id).catch(e => noteError(e.message))}>
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
            {taskInfo?.status === StatusRunning &&
                <IconButton onClick={e => startBaiduPhotoTask(false, driver.id)}>
                    <Stop />
                </IconButton>
            }
            {taskInfo ? String(taskInfo.cnt) + "/" + taskInfo.total : "loading..."}
            </Attr>
        </Grid>
    )
};
