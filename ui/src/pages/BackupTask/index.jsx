import useResourceManager from 'hox/resourceManager';
import {Stack} from "@mui/material";
import {useEffect} from "react";
import {listBranch} from "api/branch";

export default function () {
    const [resourceManager, setResourceManager] = useResourceManager();
    useEffect(() => {
        listBranch(setResourceManager);
    }, []);
    return (
        <Stack style={{width: "100%", height: "100%"}}>
            备份任务
        </Stack>
    );
}
