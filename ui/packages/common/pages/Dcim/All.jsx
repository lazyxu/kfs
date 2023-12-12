import { Grid } from "@mui/material";
import { useEffect, useRef, useState } from "react";
import Image from "./Image";

function calImageWith(gridWith) {
    const n = gridWith / 100;
    return gridWith / Math.ceil(n);
}

export default function ({ metadataList }) {
    const ref = useRef(null);

    const [width, setWidth] = useState(0);

    useEffect(() => {
        setWidth(calImageWith(ref.current.offsetWidth));
        const onResize = () => {
            console.log("resize")
            setWidth(calImageWith(ref.current.offsetWidth));
        };
        window.addEventListener("resize", onResize);
        return () => {
            window.removeEventListener("resize", onResize);
        }
    }, []);
    return (
        <Grid container ref={ref} sx={{ width: "100%" }}>
            {width && metadataList.map(metadata => {
                return <Grid item style={{ width: width, height: width }} key={metadata.hash}>
                    <Image metadata={metadata} />
                </Grid>
            })}
        </Grid>
    );
}
