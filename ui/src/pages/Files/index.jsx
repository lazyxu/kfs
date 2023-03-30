import useResourceManager from 'hox/resourceManager';
import {Stack} from "@mui/material";
import DirItems from "./DirItems";
import Branches from "./Branches";
import {useEffect} from "react";
import {listBranch} from "api/branch";
import File from "./File";

export default function ({show}) {
    const [resourceManager, setResourceManager] = useResourceManager();
    useEffect(() => {
        listBranch(setResourceManager);
    }, []);
    if (!show) {
        return;
    }
    return (
        <Stack style={{width: "100%", height: "100%"}}>
            {resourceManager.branches && <Branches/>}
            {resourceManager.file && <File file={resourceManager.file}/>}
            {resourceManager.dirItems && <DirItems dirItems={resourceManager.dirItems}/>}
        </Stack>
    );
}
