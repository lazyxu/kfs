import './index.scss';
import { useClick } from "use";
import { list, openFile } from "api/fs";
import { modeIsDir } from "api/utils/api";
import useResourceManager from 'hox/resourceManager';
import useContextMenu from "../../hox/contextMenu";
import { Box, Stack } from "@mui/material";
import FileIcon from './FileIcon';

export default function ({ dirItem, filesElm }) {
    const [resourceManager, setResourceManager] = useResourceManager();
    const [contextMenu, setContextMenu] = useContextMenu();
    let { filePath, driverName } = resourceManager;
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
                e.stopPropagation();
                const { clientX, clientY } = e;
                let { x, y, width, height } = filesElm.current.getBoundingClientRect();
                setContextMenu({
                    type: 'file',
                    dirItem,
                    clientX, clientY,
                    x, y, width, height,
                })
            }}>
            <Box onMouseDown={useClick(null, () => {
                if (modeIsDir(mode)) {
                    list(setResourceManager, driverName, filePath);
                } else {
                    openFile(setResourceManager, driverName, filePath, dirItem);
                }
            })}>
                <FileIcon dirItem={dirItem}/>
            </Box>
            <Box kfs-attr="file" style={{ width: "100%", overflowWrap: "break-word", textAlign: "center" }}>{name}</Box>
        </Stack>
    )
};
