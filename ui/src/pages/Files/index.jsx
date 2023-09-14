import { Box } from "@mui/material";
import { listDriver } from "api/driver";
import { analyzeMetadata } from "api/web/exif";
import AbsolutePath from 'components/AbsolutePath';
import useResourceManager from 'hox/resourceManager';
import { useEffect } from "react";
import DirItems from "./DirItems";
import Drivers from "./Drivers";
import FileViewer from './FileViewer/FileViewer';

export default function () {
    const [resourceManager, setResourceManager] = useResourceManager();
    useEffect(() => {
        listDriver(setResourceManager);
        analyzeMetadata(true);
    }, []);
    return (
        <Box sx={{ display: 'flex', flex: "1", flexDirection: 'column', minHeight: '0' }}>
            <AbsolutePath />
            {resourceManager.drivers && <Drivers />}
            {resourceManager.file && <FileViewer file={resourceManager.file} />}
            {resourceManager.dirItems && <DirItems dirItems={resourceManager.dirItems} />}
        </Box>
    );
}
