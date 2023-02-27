import {useRef} from 'react';
import File from "components/File";
import AbsolutePath from "components/AbsolutePath";
import useResourceManager from 'hox/resourceManager';
import DefaultContextMenu from "components/ContextMenu/DefaultContextMenu";
import useContextMenu from "hox/contextMenu";
import FileContextMenu from "components/ContextMenu/FileContextMenu";
import FileViewer from "./FileViewer/FileViewer";
import Dialog from "components/Dialog";
import {Box, Grid, Stack} from "@mui/material";

export default function () {
    const [resourceManager, setResourceManager] = useResourceManager();
    const [contextMenu, setContextMenu] = useContextMenu();
    const filesElm = useRef(null);
    return (
        <>
            <AbsolutePath/>
            {resourceManager.file ?
                <FileViewer file={resourceManager.file}/> :
                <Box style={{flex: "auto", overflowY: "scroll"}}>
                    <Grid container padding={1} spacing={1}
                          ref={filesElm} onContextMenu={(e) => {
                        e.preventDefault();
                        // console.log(e.target, e.currentTarget, e.target === e.currentTarget);
                        // if (e.target === e.currentTarget) {
                        const {clientX, clientY} = e;
                        let {x, y, width, height} = e.currentTarget.getBoundingClientRect();
                        setContextMenu({
                            type: 'default',
                            clientX, clientY,
                            x, y, width, height,
                        })
                        // }
                    }}>
                        {resourceManager.dirItems.map((dirItem, i) => (
                            <Grid item key={dirItem.name}>
                                <File filesElm={filesElm} dirItem={dirItem} key={dirItem.name}/>
                            </Grid>
                        ))}
                    </Grid>
                </Box>
            }
            {resourceManager.dirItems &&
                <Stack className='filePath'
                       direction="row"
                       justifyContent="flex-start"
                       alignItems="center"
                       spacing={1}
                >
                    共{resourceManager.dirItems.length}个项目
                </Stack>}
            <DefaultContextMenu/>
            <FileContextMenu/>
            <Dialog/>
        </>
    );
}
