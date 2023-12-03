import { Box, Stack } from "@mui/material";
import useResourceManager from 'hox/resourceManager';
import { useInView } from "react-intersection-observer";
import FileIcon from './FileIcon';
import './index.scss';

export default function ({ dirItem, setContextMenu }) {
    const [resourceManager, setResourceManager] = useResourceManager();
    const { driver, filePath } = resourceManager;
    const { name } = dirItem
    const curFilePath = filePath.concat(name);
    const { ref, inView } = useInView({
      threshold: 0
    });
    return (
        <Stack ref={ref} sx={{ ":hover": { backgroundColor: (theme) => theme.palette.action.hover } }}
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
                <FileIcon dirItem={dirItem} filePath={curFilePath} inView={inView} />
            </Box>
            <Box kfs-attr="file" style={{ width: "100%", overflowWrap: "break-word", textAlign: "center" }}>{name}</Box>
        </Stack>
    )
};
