import useResourceManager from '@kfs/common/hox/resourceManager';
import { Stack } from '@mui/material';
import AbsolutePath from './AbsolutePath';
import DriverFiles from "./DriverFiles/DriverFiles";

export default function () {
    const [resourceManager, setResourceManager] = useResourceManager();
    return (
        <Stack style={{ flex: "1", display: 'flex', flexDirection: 'column', minHeight: '0' }}>
            <AbsolutePath />
            <DriverFiles />
        </Stack>
    );
}
