import { Grid } from "@mui/material";
import Image from "./Image";
import { parseShotEquipment } from "api/utils/api";

export default function ({ exifMap, chosenShotEquipment }) {
    return (
        <Grid container spacing={1} style={{ overflowY: "scroll" }}>
            {Object.keys(exifMap).sort((a, b) => exifMap[a].shotTime - exifMap[b].shotTime)
                .filter(hash => chosenShotEquipment.includes(exifMap[hash].shotEquipment)).map(hash => {
                    return <Grid item style={{ width: "256px", height: "256px" }} key={hash}>
                        <Image hash={hash} exif={exifMap[hash]} />
                    </Grid>
                })}
        </Grid>
    );
}
