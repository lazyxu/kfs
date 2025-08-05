import useResourceManager from "@kfs/common/hox/resourceManager";
import { getSysConfig } from "@kfs/common/hox/sysConfig";
import { noteError } from "@kfs/mui/components/Notification";
import File from "@kfs/mui/pages/Files/DriverFiles/File";
import { EventStreamContentType, fetchEventSource } from "@microsoft/fetch-event-source";
import { Grid, Stack } from "@mui/material";
import { useEffect, useRef, useState } from 'react';
import DirMenu from "./DirMenu";
import FileMenu from "./FileMenu";
import FileAttributeWithPath from "./FileAttributeWithPath";

export default function () {
    const [driverFiles, setDriverFiles] = useState([]);
    const [driverFilesTotal, setDriverFilesTotal] = useState(0);
    const [dirLogicFile, setDirLogicFile] = useState();
    const [resourceManager, setResourceManager] = useResourceManager();
    const { dirPath } = resourceManager;
    const abortControllerRef = useRef(null);
    const [dirMenu, setDirMenu] = useState(null);
    const [fileMenu, setFileMenu] = useState(null);
    const [fileAttribute, setFileAttribute] = useState(null);
    const refresh = () => {
        console.log("===refresh===");
        setDriverFiles([]);
        setDriverFilesTotal(0);
        if (abortControllerRef.current) {
            abortControllerRef.current.abort();
        }
        abortControllerRef.current = new AbortController();
        fetchEventSource(`${getSysConfig().webServer}/api/v1/event/list?${dirPath.map(f => "filePath[]=" + f).join("&")}`, {
            signal: abortControllerRef.current.signal,
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
                    setDirLogicFile(info.dir);
                }
            },
            onclose() {
                // if the server closes the connection unexpectedly, retry:
                // noteError("event/list.onclose");
            },
            onerror(err) {
                console.error(err);
            }
        });
    }
    useEffect(() => {
        refresh();
        return () => {
            console.log("===abort===");
            if (abortControllerRef.current) {
                abortControllerRef.current.abort();
            }
        }
    }, [resourceManager.dirPath]);
    console.log("===render DriverFiles", dirPath, driverFiles, dirLogicFile, driverFilesTotal);
    return (
        <>
            <DirMenu contextMenu={dirMenu} setContextMenu={setDirMenu}
                refresh={() => refresh()}
                setFileAttribute={setFileAttribute}
                dirPath={dirPath}
                driverFile={dirLogicFile}
            />
            <Grid container padding={1} spacing={1}
                sx={{ flex: "1", overflowY: 'auto', alignContent: "flex-start" }}
                onContextMenu={(e) => {
                    e.preventDefault(); e.stopPropagation();
                    setDirMenu({
                        mouseX: e.clientX, mouseY: e.clientY, open: true,
                        dirPath,
                        driverFile: dirLogicFile, filePath: dirPath, // for FileAttribute
                    });
                }}
            >
                {driverFiles.map((driverFile, i) => (
                    <Grid key={driverFile.name}>
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
            {fileMenu && <FileMenu contextMenu={fileMenu} setContextMenu={setFileMenu} setFileAttribute={setFileAttribute} />}
            {fileAttribute && <FileAttributeWithPath fileAttribute={fileAttribute} onClose={() => setFileAttribute(null)} />}
        </>
    );
}
