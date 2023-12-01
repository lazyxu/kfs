import { Box } from "@mui/material";
import { listDriver } from "api/driver";
import AbsolutePath from 'components/AbsolutePath';
import useResourceManager from 'hox/resourceManager';
import { useEffect } from "react";
import DriverFiles from "./DriverFiles";
import Drivers from "./Drivers";

export default function () {
    const [resourceManager, setResourceManager] = useResourceManager();
    useEffect(() => {
        listDriver(setResourceManager);
    }, []);
    console.log("resourceManager", resourceManager, resourceManager.hasOwnProperty("driverId"));
    return (
        <Box sx={{ display: 'flex', flex: "1", flexDirection: 'column', minHeight: '0' }}>
            <AbsolutePath />
            {resourceManager.drivers && <Drivers />}
            {resourceManager.hasOwnProperty("driverId") && <DriverFiles />}
        </Box>
    );
}
