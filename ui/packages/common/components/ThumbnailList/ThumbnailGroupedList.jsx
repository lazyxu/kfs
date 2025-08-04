import { Box, Grid } from "@mui/material";
import { useEffect, useRef, useState } from "react";
import { calImageWidth, getScrollbarWidth } from "./ThumbnailList";
import Thumbnail from "@kfs/common/components/Thumbnail/Thumbnail";

export default function ({ metadataList, formatFn, spacing, minCol }) {
    let groupedMap = {};
    metadataList.forEach(metadata => {
        let group = formatFn(metadata.parsedTime);
        if (groupedMap.hasOwnProperty(group)) {
            groupedMap[group].push(metadata);
        } else {
            groupedMap[group] = [metadata];
        }
    })

    const ref = useRef(null);
    const [width, setWidth] = useState(0);

    useEffect(() => {
        const sbWidth = getScrollbarWidth();
        const updateWidth = () => {
            if (ref.current) {
                const gridWidth = ref.current.offsetWidth;
                const imgWidth = calImageWidth(minCol, gridWidth, sbWidth, spacing);
                console.log("Scrollbar width:", sbWidth);
                console.log("Grid width:", gridWidth);
                console.log("Calculated Thumbnail width:", imgWidth);
                setWidth(imgWidth);
            }
        };

        updateWidth();
        window.addEventListener("resize", updateWidth);
        return () => window.removeEventListener("resize", updateWidth);
    }, []);

    return (
        <Grid container ref={ref} spacing={spacing} sx={{ width: "100%", overflowX: "hidden" }}>
            {Object.keys(groupedMap).map(date => <Grid xs={12} style={{ width: "100%" }} key={date}>
                <Box sx={{ padding: "8px 0"}}>{date} {groupedMap[date].length}项</Box>
                <Grid container spacing={spacing} sx={{ width: "100%"}}>
                    {groupedMap[date].map(metadata => {
                        return <Grid item style={{ width: width, height: width }} key={metadata.id}>
                            <Thumbnail metadata={metadata} width={width} />
                        </Grid>
                    })}
                </Grid>
            </Grid>
            )}
        </Grid>
    );
}
