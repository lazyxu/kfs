import {useRef} from 'react';
import AbsolutePath from "components/AbsolutePath";
import useResourceManager from 'hox/resourceManager';
import Branch from "../../components/File/Branch";
import {Grid} from "@mui/material";

export default function () {
    const [resourceManager, setResourceManager] = useResourceManager();
    const filesElm = useRef(null);

    return (
        <>
            <AbsolutePath/>
            <Grid container margin={1} spacing={1}
                  style={{overflowY: "scroll"}}
                  bottom="0"
                  position="relative"
                  ref={filesElm}>
                {resourceManager.branches.map((branch, i) => (
                    <Grid item={true} key={branch.name}>
                        <Branch branch={branch}>{branch.name}</Branch>
                    </Grid>
                ))}
            </Grid>
        </>
    );
}
