import {useRef} from 'react';
import File from "components/File";
import AbsolutePath from "components/AbsolutePath";
import DefaultContextMenu from "components/ContextMenu/DefaultContextMenu";
import useContextMenu from "hox/contextMenu";
import FileContextMenu from "components/ContextMenu/FileContextMenu";
import Dialog from "components/Dialog";
import {Box, Grid, Stack} from "@mui/material";

export default function ({dirItems}) {
    const [contextMenu, setContextMenu] = useContextMenu();
    const filesElm = useRef(null);
    return (
        <>
            <AbsolutePath/>
            <Box style={{flex: "auto"}}
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
                <Grid container padding={1} spacing={1}>
                    {dirItems.map((dirItem, i) => (
                        <Grid item key={dirItem.name}>
                            <File filesElm={filesElm} dirItem={dirItem} key={dirItem.name}/>
                        </Grid>
                    ))}
                </Grid>
            </Box>
            <Stack className='filePath'
                   direction="row"
                   justifyContent="flex-start"
                   alignItems="center"
                   spacing={1}
            >
                共{dirItems.length}个项目
            </Stack>
            <DefaultContextMenu/>
            <FileContextMenu/>
            <Dialog/>
        </>
    );
}
