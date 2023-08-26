import { Grid } from "@mui/material";
import Image from "./Image";
import { timeSortFn } from "api/utils/api";

export default function ({ metadataList, chosenShotEquipment, chosenFileType }) {
    return (
        <Grid container spacing={1} style={{ overflowY: "scroll" }}>
            {metadataList.filter(metadata => chosenShotEquipment.includes(metadata.shotEquipment) && chosenFileType.includes(metadata.fileType.extension))
                .sort(timeSortFn).map(metadata => {
                    return <Grid item style={{ width: "256px", height: "256px" }} key={metadata.hash}>
                        <Image metadata={metadata} />
                    </Grid>
                })}
        </Grid>
    );
}
