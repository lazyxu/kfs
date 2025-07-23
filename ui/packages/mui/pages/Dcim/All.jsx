import { Grid } from "@mui/material";
import { useEffect, useRef, useState } from "react";
import Image from "./Image";

export function getScrollbarWidth() {
    const scrollDiv = document.createElement('div');
    scrollDiv.style.cssText = `
        position: absolute;
        top: -9999px;
        width: 100px;
        height: 100px;
        overflow: scroll;
    `;
    document.body.appendChild(scrollDiv);
    const scrollbarWidth = scrollDiv.offsetWidth - scrollDiv.clientWidth;
    document.body.removeChild(scrollDiv);
    return scrollbarWidth;
}

export function calImageWidth(minCol, gridWidth, scrollbarWidth, spacing) {
    const colCount = Math.ceil(Math.max(minCol, Math.floor((gridWidth - scrollbarWidth) / 256)));
    const totalGap = (colCount - 1) * 8 * spacing;
    const availableWidth = gridWidth - scrollbarWidth - totalGap;
    return availableWidth / Math.ceil(colCount);
}

const spacing = 0; // 1 for 8px spacing
const minCol = 10;

export default function ({ metadataList }) {
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
                console.log("Calculated image width:", imgWidth);
                setWidth(imgWidth);
            }
        };

        updateWidth();
        window.addEventListener("resize", updateWidth);
        return () => window.removeEventListener("resize", updateWidth);
    }, []);

    return (
        <Grid container ref={ref} spacing={spacing} sx={{ width: "100%", overflowX: "hidden" }}>
            {width > 0 && metadataList.map(metadata => (
                <Grid item style={{ width: width, height: width }} key={metadata.id}>
                    <Image metadata={metadata} width={width} />
                </Grid>
            ))}
        </Grid>
    );
}
