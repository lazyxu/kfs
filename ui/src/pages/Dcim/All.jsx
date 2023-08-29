import { Grid } from "@mui/material";
import Image from "./Image";

export default function ({ metadataList }) {
    return (
        <Grid container spacing={1} style={{ overflowY: "scroll" }}>
            {metadataList.map(metadata => {
                return <Grid item style={{ width: "256px", height: "256px" }} key={metadata.hash}>
                    <Image metadata={metadata} />
                </Grid>
            })}
        </Grid>
    );
}
