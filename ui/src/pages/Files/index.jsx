import { Box } from "@mui/material";
import AbsolutePath from 'components/AbsolutePath';
import useResourceManager from 'hox/resourceManager';
import DriverFiles from "./DriverFile/DriverFiles";
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
