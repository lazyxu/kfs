import { Box, Stack } from "@mui/material";
import useResourceManager from 'hox/resourceManager';
import FileIcon from './FileIcon';
import './index.scss';

export default function ({ dirItem, setContextMenu }) {
    const [resourceManager, setResourceManager] = useResourceManager();
    let { driverId, driverName, filePath } = resourceManager;
    const { name, mode } = dirItem;
    filePath = filePath.concat(name);
    return (
        <Stack component="span" sx={{ ":hover": { backgroundColor: (theme) => theme.palette.action.hover } }}
            className='file-normal'
            justifyContent="flex-start"
            alignItems="center"
            spacing={1}
            onContextMenu={(e) => {
                e.preventDefault();
                setContextMenu({
                    mouseX: e.clientX, mouseY: e.clientY,
                    driverId, driverName, filePath, mode
                });
            }}
        >
            <Box>
                <FileIcon dirItem={dirItem} filePath={filePath} />
            </Box>
            <Box kfs-attr="file" style={{ width: "100%", overflowWrap: "break-word", textAlign: "center" }}>{name}</Box>
        </Stack>
    )
};
