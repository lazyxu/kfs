import { Box, Stack } from "@mui/material";
import useResourceManager from 'hox/resourceManager';
import { useEffect, useRef } from "react";
import { useInView } from "react-intersection-observer";
import FileIcon from './FileIcon';
import './index.scss';

export default ({ dirItem, setContextMenu }) => {
    const [resourceManager, setResourceManager] = useResourceManager();
    const { driver, filePath } = resourceManager;
    const { name } = dirItem
    const curFilePath = useRef([]);
    const { ref, inView } = useInView({ threshold: 0 });
    const hasBeenInView = useRef(false);
    useEffect(() => {
        curFilePath.current = filePath.concat(name);
    }, []);
    useEffect(() => {
        if (!inView || hasBeenInView.current) {
            return;
        }
        hasBeenInView.current = true;
    }, [inView]);
    return (
        <Stack ref={ref} sx={{ ":hover": { backgroundColor: (theme) => theme.palette.action.hover } }}
            className='file-normal'
            justifyContent="flex-start"
            alignItems="center"
            spacing={1}
            onContextMenu={(e) => {
                e.preventDefault(); e.stopPropagation();
                setContextMenu({
                    mouseX: e.clientX, mouseY: e.clientY,
                    driver, filePath, dirItem
                });
            }}
        >
            <Box>
                <FileIcon dirItem={dirItem} filePath={curFilePath.current} hasBeenInView={hasBeenInView.current} driver={driver} />
            </Box>
            <Box kfs-attr="file" style={{ width: "100%", overflowWrap: "break-word", textAlign: "center" }}>{name}</Box>
        </Stack>
    )
};
