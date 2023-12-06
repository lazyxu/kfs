import { Box, Drawer } from "@mui/material";
import { listDriverFileByHash } from "api/fs";
import { useEffect } from "react";

export default function ({ hash, open, setOpen, sameFiles, setSameFiles }) {
    useEffect(() => {
        listDriverFileByHash(hash).then(setSameFiles);
    }, []);
    return (
        <Drawer
            anchor="right"
            open={open}
            onClose={() => setOpen(false)}
            sx={{ zIndex: 1350 }}
            SlideProps={{ sx: { maxWidth: "90%" } }}
        >
            <Box
                sx={{ whiteSpace: "pre" }}
            >
                {sameFiles.map((f, i) => <Box key={i}>{f.driverName}{f.dirPath.length ? ("/" + f.dirPath.join("/") + "/" + f.name) : ("/" + f.name)}</Box>)}
            </Box>
        </Drawer>
    );
}
