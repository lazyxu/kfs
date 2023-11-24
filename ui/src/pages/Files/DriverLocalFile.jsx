import { EventStreamContentType, fetchEventSource } from '@microsoft/fetch-event-source';
import { HourglassDisabled, HourglassTop, PlayArrow, Stop } from '@mui/icons-material';
import { Box, Grid, IconButton } from "@mui/material";
import { startDriverLocalFile } from 'api/web/exif';
import { noteError } from 'components/Notification/Notification';
import { getSysConfig } from 'hox/sysConfig';
import humanize from 'humanize';
import humanizeDuration from "humanize-duration";
import moment from 'moment';
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

export default ({ driver, attributes }) => {
    const [info, setInfo] = useState();
    const controller = new AbortController();
    useEffect(() => {
        fetchEventSource(`http://127.0.0.1:${getSysConfig().sysConfig.port}/api/v1/event/driverLocalFile/${driver.id}`, {
            signal: controller.signal,
            async onopen(response) {
                if (response.ok && response.headers.get('content-type').includes(EventStreamContentType)) {
                    return; // everything's good
                }
                console.error(response);
                noteError("event/driverLocalFile.onopen: " + response.status);
            },
            onmessage(msg) {
                // if the server emits an error message, throw an exception
                // so it gets handled by the onerror callback below:
                if (msg.event === 'FatalError') {
                    console.error(msg);
                    noteError("event/driverLocalFile.onmessage: " + msg);
                    return;
                }
                let info = JSON.parse(msg.data);
                console.log(info);
                if (info.errMsg) {
                    noteError(info.errMsg);
                }
                setInfo(info);
            },
            onclose() {
                // if the server closes the connection unexpectedly, retry:
                noteError("event/driverLocalFile.onclose");
            },
            onerror(err) {
                console.error(err);
                // noteError("event/driverLocalFile.onerror: " + err.message);
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
    let curFile = info?.curFile ? info.curFile : info?.curDir ? info.curDir : "";
    return (
        <>
            <Attr k="上次同步结束时间">{info?.lastDoneTime ? `${moment(info.lastDoneTime / 1000 / 1000).format("YYYY年MM月DD日 HH:mm:ss")}` : "?"}</Attr>
            <Attr k="同步">{(info?.status === undefined ||
                info?.status === StatusIdle ||
                info?.status === StatusFinished ||
                info?.status === StatusCanceled ||
                info?.status === StatusError) &&
                <IconButton onClick={e => startDriverLocalFile(true, driver.id, attributes.srcPath, attributes.encoder)}>
                    <PlayArrow />
                </IconButton>
            }
                {info?.status === StatusWaitRunning &&
                    <IconButton>
                        <HourglassTop />
                    </IconButton>
                }
                {info?.status === StatusWaitCanceled &&
                    <IconButton>
                        <HourglassDisabled />
                    </IconButton>
                }
                {info?.status === StatusRunning &&
                    <IconButton onClick={e => startDriverLocalFile(false, driver.id, attributes.srcPath, attributes.encoder)}>
                        <Stop />
                    </IconButton>
                }
            </Attr>
            <Attr k="耗时">{info ? `${humanizeDuration(Math.floor(info.cost / 100) * 100)}` : "?"}</Attr>
            <Attr k="当前文件或目录">
                <a variant="text" onClick={() => {
                    const { shell } = window.require('@electron/remote');
                    shell.openPath(curFile);
                }} >
                    <div style={{ whiteSpace: "nowrap", textOverflow: "ellipsis", overflow: "hidden", display: "block" }}>
                        {curFile}
                    </div>
                </a>
            </Attr>
            <Attr k="当前文件大小">{info?.curFile ? humanize.filesize(info.curSize) : ""}</Attr>
            <Attr k="当前目录下文件数量">{info?.curDir ? info.curDirItemCnt : ""}</Attr>
            <Attr k="同步大小">{info ? `${humanize.filesize(info.size)}/${humanize.filesize(info.totalSize)}` : "?"}</Attr>
            <Attr k="同步文件数量">{info ? `${info.fileCount}/${info.totalFileCount}` : "?"}</Attr>
            <Attr k="同步目录数量">{info ? `${info.dirCount}/${info.totalDirCount}` : "?"}</Attr>
            <Attr k="同步失败的文件或目录">{info ? info.warnings.map((w, i)=> (<ul key={i}><li>{w}</li></ul>)) : "?"}</Attr>
        </>
    )
};
