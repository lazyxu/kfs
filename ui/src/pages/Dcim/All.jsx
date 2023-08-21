import { Grid } from "@mui/material";
import Image from "./Image";
import { timeSortFn } from "api/utils/api";

export default function ({ exifMap, chosenShotEquipment }) {
    return (
        <Grid container spacing={1} style={{ overflowY: "scroll" }}>
            {Object.keys(exifMap).filter(hash => chosenShotEquipment.includes(exifMap[hash].shotEquipment))
                .sort((a, b) => timeSortFn(exifMap, a, b)).map(hash => {
                    return <Grid item style={{ width: "256px", height: "256px" }} key={hash}>
                        <Image hash={hash} exif={exifMap[hash]} />
                    </Grid>
                })}
        </Grid>
    );
}
