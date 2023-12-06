import { Box, Stack } from "@mui/material";
import useResourceManager from 'hox/resourceManager';
import { useCallback, useEffect, useRef, useState } from "react";
import { useInView } from "react-intersection-observer";
import FileIcon from './FileIcon';
import './index.scss';

const useGetState = (initiateState) => {
    const [state, setState] = useState(initiateState);
    const stateRef = useRef(state);
    stateRef.current = state;
    const getState = useCallback(() => stateRef.current, []);
    return [state, setState, getState];
};

export default ({ dirItem, setContextMenu }) => {
    const [resourceManager, setResourceManager] = useResourceManager();
    const { driver, dirPath } = resourceManager;
    const { name } = dirItem
    const filePath = useRef([]);
    const { ref, inView } = useInView({ threshold: 0 });
    const [hasBeenInView, setHasBeenInView, getHasBeenInView] = useGetState(false);
    useEffect(() => {
        filePath.current = dirPath.concat(name);
    }, []);
    useEffect(() => {
        const hasBeenIn = getHasBeenInView();
        if (!inView || hasBeenIn) {
            return;
        }
        setHasBeenInView(true);
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
                    driver, filePath: filePath.current, dirItem
                });
            }}
        >
            {filePath.current.length && <FileIcon dirItem={dirItem} filePath={filePath.current} hasBeenInView={hasBeenInView} driver={driver} inView={inView}/>}
            <Box kfs-attr="file" style={{ width: "100%", overflowWrap: "break-word", textAlign: "center" }}>{name}</Box>
        </Stack>
    )
};
