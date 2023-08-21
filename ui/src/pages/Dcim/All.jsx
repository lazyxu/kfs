import { Grid } from "@mui/material";
import Image from "./Image";

export default function ({ exifMap, chosenModel }) {
    return (
        <Grid container spacing={1} style={{ overflowY: "scroll" }}>
            {Object.keys(exifMap).sort((a, b) => exifMap[a].DateTime - exifMap[b].DateTime)
                .filter(hash => chosenModel.includes(exifMap[hash].Model)).map(hash => {
                    return <Grid item style={{ width: "256px", height: "256px" }} key={hash}>
                        <Image hash={hash} exif={exifMap[hash]} />
                    </Grid>
                })}
        </Grid>
    );
}
