import { EventStreamContentType, fetchEventSource } from '@microsoft/fetch-event-source';
import { HourglassDisabled, HourglassTop, PlayArrow, Stop } from '@mui/icons-material';
import { IconButton } from "@mui/material";
import { startBaiduPhotoTask } from 'api/web/exif';
import { noteError } from 'components/Notification/Notification';
import { getSysConfig } from 'hox/sysConfig';
import { useEffect, useState } from 'react';

const StatusIdle = 0
const StatusFinished = 1
const StatusCanceled = 2
const StatusError = 3
const StatusWaitRunning = 4
const StatusWaitCanceled = 5
const StatusRunning = 6

export default ({ driver }) => {
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
        <>
            {(taskInfo?.status === undefined ||
                taskInfo?.status === StatusIdle ||
                taskInfo?.status === StatusFinished ||
                taskInfo?.status === StatusCanceled ||
                taskInfo?.status === StatusError) &&
                <IconButton onClick={e => startBaiduPhotoTask(true, driver.name).catch(e => noteError(e.message))}>
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
                <IconButton onClick={e => startBaiduPhotoTask(false, driver.name)}>
                    <Stop />
                </IconButton>
            }
            {taskInfo ? String(taskInfo.cnt) + "/" + taskInfo.total : "loading..."}
        </>
    )
};
