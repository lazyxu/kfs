import useResourceManager from 'hox/resourceManager';
import {Stack} from "@mui/material";
import DirItems from "./DirItems";
import Drivers from "./Drivers";
import {useEffect} from "react";
import {listDriver} from "api/driver";
import File from "./File";

export default function ({show}) {
    const [resourceManager, setResourceManager] = useResourceManager();
    useEffect(() => {
        listDriver(setResourceManager);
    }, []);
    return (
        <Stack style={{width: "100%", height: "100%", display: show?undefined:"none"}}>
            {resourceManager.drivers && <Drivers/>}
            {resourceManager.file && <File file={resourceManager.file}/>}
            {resourceManager.dirItems && <DirItems dirItems={resourceManager.dirItems}/>}
        </Stack>
    );
}
