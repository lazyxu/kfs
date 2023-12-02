import { Box, Stack } from "@mui/material";
import useResourceManager from 'hox/resourceManager';
import FileIcon from './FileIcon';
import './index.scss';

export default function ({ dirItem, setContextMenu }) {
    const [resourceManager, setResourceManager] = useResourceManager();
    const { driver, filePath } = resourceManager;
    const { name } = dirItem
    const curFilePath = filePath.concat(name);
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
                    driver, filePath, dirItem
                });
            }}
        >
            <Box>
                <FileIcon dirItem={dirItem} filePath={curFilePath} />
            </Box>
            <Box kfs-attr="file" style={{ width: "100%", overflowWrap: "break-word", textAlign: "center" }}>{name}</Box>
        </Stack>
    )
};
