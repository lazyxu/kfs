import { EventStreamContentType, fetchEventSource } from "@microsoft/fetch-event-source";
import { Box, Grid, Stack } from "@mui/material";
import File from "components/File";
import { noteError } from "components/Notification/Notification";
import useResourceManager from "hox/resourceManager";
import { getSysConfig } from "hox/sysConfig";
import { useEffect, useState } from 'react';
import FileAttribute from "./FileAttribute";
import FileMenu from "./FileMenu";
import Menu from "./Menu";

export default function () {
    const [dirItems, setDirItems] = useState([]);
    const [dirItemsTotal, setDirItemsTotal] = useState(0);
    const [resourceManager, setResourceManager] = useResourceManager();
    let { driver, filePath } = resourceManager;
    const controller = new AbortController();
    const [menu, setMenu] = useState(null);
    const [fileMenu, setFileMenu] = useState(null);
    let [fileAttribute, setFileAttribute] = useState(null);
    useEffect(() => {
        setDirItems([]);
        fetchEventSource(`${getSysConfig().sysConfig.webServer}/api/v1/event/list?driverId=${driver.id}&${filePath.map(f => "filePath[]=" + f).join("&")}`, {
            signal: controller.signal,
            openWhenHidden: true,
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
                console.log(info);
                if (info.errMsg) {
                    noteError(info.errMsg);
                }
                if (info.files) {
                    setDirItems(prev => [...prev, ...info.files]);
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
            <Box style={{ flex: "1", overflowY: 'auto', alignContent: "flex-start" }} >
                <Grid container padding={1} spacing={1}
                    onContextMenu={(e) => {
                        e.preventDefault(); e.stopPropagation();
                        setMenu({
                            mouseX: e.clientX, mouseY: e.clientY,
                            driver, filePath,
                        });
                    }}
                >
                    {dirItems.map((dirItem, i) => (
                        <Grid item key={dirItem.name}>
                            <File setContextMenu={setFileMenu} dirItem={dirItem} key={dirItem.name} />
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
            {menu && <Menu contextMenu={menu} setContextMenu={setMenu} />}
            {fileMenu && <FileMenu contextMenu={fileMenu} setContextMenu={setFileMenu} setFileAttribute={setFileAttribute} />}
            {fileAttribute && <FileAttribute fileAttribute={fileAttribute} setFileAttribute={setFileAttribute} />}
        </>
    );
}
