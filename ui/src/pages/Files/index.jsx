import { Box } from "@mui/material";
import useResourceManager from 'hox/resourceManager';
import AbsolutePath from './AbsolutePath';
import DriverFiles from "./DriverFiles/DriverFiles";
import Drivers from "./Drivers";

export default function () {
    const [resourceManager, setResourceManager] = useResourceManager();
    return (
        <Box sx={{ display: 'flex', flex: "1", flexDirection: 'column', minHeight: '0' }}>
            <AbsolutePath />
            {resourceManager.hasOwnProperty("driver") ? <DriverFiles /> : <Drivers />}
        </Box>
    );
}
