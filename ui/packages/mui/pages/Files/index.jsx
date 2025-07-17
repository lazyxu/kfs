import useResourceManager from '@kfs/common/hox/resourceManager';
import { Box } from "@mui/material";
import AbsolutePath from './AbsolutePath';
import DriverFiles from "./DriverFiles/DriverFiles";

export default function () {
    const [resourceManager, setResourceManager] = useResourceManager();
    return (
        <Box sx={{ flex: "1", display: 'flex', flexDirection: 'column', minHeight: '0',
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.primary
         }}>
            <AbsolutePath />
            <DriverFiles />
        </Box>
    );
}
