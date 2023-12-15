import useResourceManager from '@kfs/common/hox/resourceManager';
import { Box } from "@mui/material";
import AbsolutePath from './AbsolutePath';
import DriverFiles from "./DriverFiles/DriverFiles";
import Drivers from "./Drivers";

export default function () {
    const [resourceManager, setResourceManager] = useResourceManager();
    return (
        <Box sx={{ flex: "1", display: 'flex', flexDirection: 'column', minHeight: '0' }}>
            <AbsolutePath />
            {resourceManager.hasOwnProperty("driver") ? <DriverFiles /> : <Drivers />}
        </Box>
    );
}
