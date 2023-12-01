import { EventStreamContentType, fetchEventSource } from "@microsoft/fetch-event-source";
import { Box, Grid, Stack } from "@mui/material";
import DefaultContextMenu from "components/ContextMenu/DefaultContextMenu";
import FileContextMenu from "components/ContextMenu/FileContextMenu";
import Dialog from "components/Dialog";
import File from "components/File";
import { noteError } from "components/Notification/Notification";
import useContextMenu from "hox/contextMenu";
import useResourceManager from "hox/resourceManager";
import { getSysConfig } from "hox/sysConfig";
import { useEffect, useRef, useState } from 'react';

export default function () {
    const [dirItems, setDirItems] = useState([]);
    const [dirItemsTotal, setDirItemsTotal] = useState(0);
    const [resourceManager, setResourceManager] = useResourceManager();
    let { driverId, filePath } = resourceManager;
    const [contextMenu, setContextMenu] = useContextMenu();
    const filesElm = useRef(null);
    const controller = new AbortController();
    useEffect(() => {
        setDirItems([]);
        fetchEventSource(`${getSysConfig().sysConfig.webServer}/api/v1/event/list?driverId=${driverId}&${filePath.map(f=>"filePath[]="+f).join("&")}`, {
            signal: controller.signal,
            async onopen(response) {
                if (response.ok && response.headers.get('content-type').includes(EventStreamContentType)) {
                    return; // everything's good
                }
                console.error(response);
                noteError("event/list.onopen: " + response.status);
            },
            onmessage(msg) {
                // if the server emits an error message, throw an exception
                // so it gets handled by the onerror callback below:
                if (msg.event === 'FatalError') {
                    console.error(msg);
                    noteError("event/list.onmessage: " + msg);
                    return;
                }
                let info = JSON.parse(msg.data);
                console.log("===info===", info);
                if (info.errMsg) {
                    noteError(info.errMsg);
                }
                if (info.file) {
                    // TODO: create too many arrays?
                    setDirItems(prev => [...prev, info.file]);
                } else {
                    setDirItemsTotal(info.n);
                }
            },
            onclose() {
                // if the server closes the connection unexpectedly, retry:
                // noteError("event/list.onclose");
            },
            onerror(err) {
                console.error(err);
                // noteError("event/list.onerror: " + err.message);
                // if (err instanceof FatalError) {
                //     throw err; // rethrow to stop the operation
                // } else {
                //     // do nothing to automatically retry. You can also
                //     // return a specific retry interval here.
                // }
            }
        });
        return () => {
            console.log("===abort===");
            controller.abort();
        }
    }, [resourceManager.filePath]);
    return (
        <>
            <Box style={{ flex: "1", overflowY: 'auto', alignContent: "flex-start" }}
                ref={filesElm} onContextMenu={(e) => {
                    e.preventDefault();
                    // console.log(e.target, e.currentTarget, e.target === e.currentTarget);
                    // if (e.target === e.currentTarget) {
                    const { clientX, clientY } = e;
                    let { x, y, width, height } = e.currentTarget.getBoundingClientRect();
                    setContextMenu({
                        type: 'default',
                        clientX, clientY,
                        x, y, width, height,
                    })
                    // }
                }}>
                <Grid container padding={1} spacing={1}>
                    {dirItems.map((dirItem, i) => (
                        <Grid item key={dirItem.name}>
                            <File filesElm={filesElm} dirItem={dirItem} key={dirItem.name} />
                        </Grid>
                    ))}
                </Grid>
            </Box>
            <Stack className='filePath'
                direction="row"
                justifyContent="flex-start"
                alignItems="center"
                spacing={1}
            >
                共{dirItems.length}/{dirItemsTotal}个项目
            </Stack>
            <DefaultContextMenu />
            <FileContextMenu />
            <Dialog />
        </>
    );
}
