import useResourceManager from 'hox/resourceManager';
import { Box, Stack } from "@mui/material";
import DirItems from "./DirItems";
import Drivers from "./Drivers";
import { useEffect } from "react";
import { listDriver } from "api/driver";
import AbsolutePath from 'components/AbsolutePath';
import FileViewer from './FileViewer/FileViewer';

export default function ({ show }) {
    const [resourceManager, setResourceManager] = useResourceManager();
    useEffect(() => {
        listDriver(setResourceManager);
    }, []);
    return (
        <Stack sx={{ display: show ? undefined : "none", overflowY: 'auto' }}>
            <AbsolutePath />
            {resourceManager.drivers && <Drivers />}
            {resourceManager.file && <FileViewer file={resourceManager.file} />}
            {resourceManager.dirItems && <DirItems dirItems={resourceManager.dirItems} />}
        </Stack>
    );
}
