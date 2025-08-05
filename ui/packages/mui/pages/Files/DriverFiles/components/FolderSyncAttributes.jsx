import { HourglassDisabled, HourglassTop, PlayArrow, Stop } from '@mui/icons-material';
import { Box, Grid, TextField, FormControlLabel, Checkbox, Button } from "@mui/material";
import IconButton from "@mui/material/IconButton";
import { useEffect, useState } from 'react';
import { updateBaiduPhoto } from '../../../../api/driver';
import { startBaiduPhotoTask } from '../../../../api/exif';
import { getSysConfig } from '@kfs/common/hox/sysConfig';
import { noteError } from '@kfs/mui/components/Notification';
import { EventStreamContentType, fetchEventSource } from '@microsoft/fetch-event-source';

function Attr({ k, children }) {
    return <>
        <Grid size={{ xs: 4 }} sx={{ overflowWrap: "anywhere" }}>
            <Box fontWeight="bold">{k}：</Box>
        </Grid>
        <Grid size={{ xs: 8 }} sx={{ overflowWrap: "anywhere" }}>
            {children}
        </Grid>
    </>
}

const StatusIdle = 0
const StatusFinished = 1
const StatusCanceled = 2
const StatusError = 3
const StatusWaitRunning = 4
const StatusWaitCanceled = 5
const StatusRunning = 6

export default ({ driverFile, onSave }) => {
    const [cron, setCron] = useState(driverFile.cron || '');
    const [enabled, setEnabled] = useState(driverFile.enabled || false);
    const [syncAtStart, setSyncAtStart] = useState(driverFile.sync_at_start || false);
    const [bdussBfess, setBdussBfess] = useState(driverFile.bduss_bfess || '');
    const [stoken, setStoken] = useState(driverFile.stoken || '');
    const [taskInfo, setTaskInfo] = useState();
    const controller = new AbortController();

    useEffect(() => {
        if (driverFile.sync_from_id > 0) {
            fetchEventSource(`${getSysConfig().webServer}/api/v1/event/baiduPhotoTask/${driverFile.id}`, {
                signal: controller.signal,
                async onopen(response) {
                    if (response.ok && response.headers.get('content-type').includes(EventStreamContentType)) {
                        return;
                    }
                    console.error(response);
                    noteError("event/baiduPhotoTask.onopen: " + response.status);
                },
                onmessage(msg) {
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
                    noteError("event/baiduPhotoTask.onclose");
                },
                onerror(err) {
                    console.error(err);
                }
            });
        }
        return () => {
            controller?.abort();
        }
    }, [driverFile.id, driverFile.sync_from_id]);

    const handleSave = () => {
        updateBaiduPhoto(driverFile.id, bdussBfess, stoken, cron, enabled, syncAtStart)
            .then(onSave);
    };

    return (
        <Box sx={{ mt: 2 }}>
            <Grid container spacing={1.5} sx={{ alignItems: "center" }}>
                <Attr k="cron表达式">
                    <TextField
                        size="small"
                        value={cron}
                        onChange={(e) => setCron(e.target.value)}
                        placeholder="0 * * * *"
                    />
                </Attr>

                <Attr k="开启同步">
                    <FormControlLabel
                        control={
                            <Checkbox
                                checked={enabled}
                                onChange={(e) => setEnabled(e.target.checked)}
                            />
                        }
                        label={enabled ? "已启用" : "已禁用"}
                    />
                </Attr>

                <Attr k="启动时同步">
                    <FormControlLabel
                        control={
                            <Checkbox
                                checked={syncAtStart}
                                onChange={(e) => setSyncAtStart(e.target.checked)}
                            />
                        }
                        label={syncAtStart ? "是" : "否"}
                    />
                </Attr>

                <Attr k="手动同步">
                    {(taskInfo?.status === undefined ||
                        taskInfo?.status === StatusIdle ||
                        taskInfo?.status === StatusFinished ||
                        taskInfo?.status === StatusCanceled ||
                        taskInfo?.status === StatusError) &&
                        <IconButton onClick={e => startBaiduPhotoTask(true, driverFile.id).catch(e => noteError(e.message))}>
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
                        <IconButton onClick={e => startBaiduPhotoTask(false, driverFile.id)}>
                            <Stop />
                        </IconButton>
                    }
                    {taskInfo ? String(taskInfo.cnt) + "/" + taskInfo.total : "loading..."}
                </Attr>

                <Attr k="BDUSS_BFESS">
                    <TextField
                        size="small"
                        value={bdussBfess}
                        onChange={(e) => setBdussBfess(e.target.value)}
                        placeholder="输入BDUSS_BFESS"
                    />
                </Attr>

                <Attr k="STOKEN">
                    <TextField
                        size="small"
                        value={stoken}
                        onChange={(e) => setStoken(e.target.value)}
                        placeholder="输入STOKEN"
                    />
                </Attr>
            </Grid>
            
            <Box sx={{ mt: 2, display: 'flex', justifyContent: 'flex-end' }}>
                <Button
                    variant="contained"
                    onClick={handleSave}
                    color="primary"
                >
                    保存同步设置
                </Button>
            </Box>
        </Box>
    );
};
