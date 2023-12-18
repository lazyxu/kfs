import useResourceManager from "@kfs/common/hox/resourceManager";
import { getSysConfig } from "@kfs/common/hox/sysConfig";
import { noteError } from "@kfs/mui/components/Notification";
import File from "@kfs/mui/pages/Files/DriverFiles/File";
import { EventStreamContentType, fetchEventSource } from "@microsoft/fetch-event-source";
import { Grid, Stack } from "@mui/material";
import { useEffect, useState } from 'react';
import FileAttribute from "./FileAttribute";
import FileMenu from "./FileMenu";
import Menu from "./Menu";

export default function () {
    const [driverFiles, setDriverFiles] = useState([]);
    const [driverFilesTotal, setDriverFilesTotal] = useState(0);
    const [resourceManager, setResourceManager] = useResourceManager();
    const { driver, dirPath } = resourceManager;
    const controller = new AbortController();
    const [menu, setMenu] = useState(null);
    const [fileMenu, setFileMenu] = useState(null);
    const [fileAttribute, setFileAttribute] = useState(null);
    useEffect(() => {
        setDriverFiles([]);
        fetchEventSource(`${getSysConfig().webServer}/api/v1/event/list?driverId=${driver.id}&${dirPath.map(f => "filePath[]=" + f).join("&")}`, {
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
                    setDriverFiles(prev => [...prev, ...info.files]);
                } else {
                    setDriverFilesTotal(info.n);
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
    }, [resourceManager.dirPath]);
    return (
        <>
            <Grid container padding={1} spacing={1}
                sx={{ flex: "1", overflowY: 'auto', alignContent: "flex-start" }}
                onContextMenu={(e) => {
                    e.preventDefault(); e.stopPropagation();
                    setMenu({
                        mouseX: e.clientX, mouseY: e.clientY,
                        driver, dirPath,
                    });
                }}
            >
                {driverFiles.map((driverFile, i) => (
                    <Grid item key={driverFile.name}>
                        <File setContextMenu={setFileMenu} driverFile={driverFile} key={driverFile.name} />
                    </Grid>
                ))}
            </Grid>
            <Stack direction="row"
                justifyContent="flex-start"
                alignItems="center"
                spacing={1}
            >
                共{driverFiles.length}/{driverFilesTotal}个项目
            </Stack>
            {menu && <Menu contextMenu={menu} setContextMenu={setMenu} />}
            {fileMenu && <FileMenu contextMenu={fileMenu} setContextMenu={setFileMenu} setFileAttribute={setFileAttribute} />}
            {fileAttribute && <FileAttribute fileAttribute={fileAttribute} onClose={() => setFileAttribute(null)} />}
        </>
    );
}
